package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/emPeeGee/raffinance/internal/account"
	"github.com/emPeeGee/raffinance/internal/analytics"
	"github.com/emPeeGee/raffinance/internal/auth"
	"github.com/emPeeGee/raffinance/internal/category"
	"github.com/emPeeGee/raffinance/internal/config"
	"github.com/emPeeGee/raffinance/internal/connection"
	"github.com/emPeeGee/raffinance/internal/contact"
	"github.com/emPeeGee/raffinance/internal/cors"
	"github.com/emPeeGee/raffinance/internal/entity"
	"github.com/emPeeGee/raffinance/internal/hub"
	"github.com/emPeeGee/raffinance/internal/seeder"
	"github.com/emPeeGee/raffinance/internal/tag"
	"github.com/emPeeGee/raffinance/internal/transaction"
	"github.com/emPeeGee/raffinance/pkg/accesslog"
	"github.com/emPeeGee/raffinance/pkg/errorutil"
	"github.com/emPeeGee/raffinance/pkg/log"
	"github.com/emPeeGee/raffinance/pkg/validatorutil"
	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

const Version = "1.0.0"

// TODO: jwt config to be extracted
// TODO: more linting linters to be added, like the linter for err check

// RUN: Before autoMigrate -> CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
func main() {
	logger := log.New().With(nil, "version", Version)

	if err := os.Setenv("TZ", "Universal"); err != nil {
		logger.Fatalf("Error setting environment variable")
	}

	cfg, err := config.Get(logger)
	if err != nil {
		logger.Fatalf("failed to initialize config: %s", err.Error())
	}

	db, err := connection.NewPostgresDB(cfg.DB)
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}

	err = db.AutoMigrate(&entity.User{}, &entity.Contact{}, &entity.Account{}, &entity.Transaction{}, &entity.TransactionType{}, &entity.Category{}, &entity.Tag{}, &entity.TransactionTag{})
	if err != nil {
		logger.Fatalf("failed to auto migrate gorm", err.Error())
	}

	seeder := seeder.NewSeeder(db, logger)
	if err := seeder.Run(); err != nil {
		logger.Fatalf("Error has occurred while seeding: %s", err.Error())
	}

	server := new(connection.Server)
	valid := validator.New()
	if err := valid.RegisterValidation("currency", validatorutil.CurrencyValidator); err != nil {
		logger.Fatalf("failed to register currency validator: %s", err.Error())
	}

	if err := valid.RegisterValidation("transactiontype", validatorutil.TransactionType); err != nil {
		logger.Fatalf("failed to register transaction type validator: %s", err.Error())
	}

	// TODO: Error handling here
	valid.RegisterStructValidation(transaction.ValidateCreateTransaction, transaction.CreateTransactionDTO{})
	valid.RegisterStructValidation(transaction.ValidateUpdateTransaction, transaction.UpdateTransactionDTO{})
	valid.RegisterStructValidation(analytics.ValidateDateRange, analytics.RangeDateParams{})

	hub := hub.NewHub()

	go func() {
		if err := server.Run(cfg.Server, buildHandler(db, valid, logger, hub)); err != nil {
			logger.Fatalf("Error occurred while running http server: %s", err.Error())
		}
	}()

	logger.Info("Raffinance Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("Raffinance Shutting Down")

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Fatalf("error occurred on server shutting down: %s", err.Error())
	}
}

// TODO: How the dependencies injection can be done better?
// TODO: Logger is not passed as ref
// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(db *gorm.DB, valid *validator.Validate, logger log.Logger, hub *hub.Hub) http.Handler {
	router := gin.New()
	router.Use(accesslog.Handler(logger), errorutil.Handler(logger), cors.Handler())

	authRg := router.Group("/auth")
	apiRg := router.Group("/api", auth.HandleUserIdentity(logger))

	apiRg.GET("/websocket", func(c *gin.Context) {
		handleWebSocketConnection(c, hub)
	})

	// transaction service is used in account as well
	transactionService := transaction.NewTransactionService(transaction.NewTransactionRepository(db, logger), logger, hub)

	auth.RegisterHandlers(
		authRg,
		apiRg,
		auth.NewAuthService(auth.NewAuthRepository(db, logger), logger),
		valid,
		logger,
	)

	contact.RegisterHandlers(
		apiRg,
		contact.NewContactService(contact.NewContactRepository(db, logger), logger),
		valid,
		logger,
	)

	account.RegisterHandlers(
		apiRg,
		account.NewAccountService(transactionService, account.NewAccountRepository(db, logger), logger),
		valid,
		logger,
	)

	transaction.RegisterHandlers(
		apiRg,
		transactionService,
		valid,
		logger,
	)

	category.RegisterHandlers(
		apiRg,
		category.NewCategoryService(category.NewCategoryRepository(db, logger), logger),
		valid,
		logger,
	)

	tag.RegisterHandlers(
		apiRg,
		tag.NewTagService(tag.NewTagRepository(db, logger), logger),
		valid,
		logger,
	)

	analytics.RegisterHandlers(
		apiRg,
		analytics.NewAnalyticsService(analytics.NewAnalyticsRepository(db, logger), logger),
		valid,
		logger,
	)

	return router
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow requests from a specific origin
		fmt.Println(r.Header.Get("Origin"))
		return r.Header.Get("Origin") == "http://localhost:3000"
	},
}

func handleWebSocketConnection(c *gin.Context, huub *hub.Hub) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	// Get the user ID from the URL params
	userID, err := auth.GetUserId(c)
	if err != nil || userID == nil {
		fmt.Println("ERRORR", err.Error())
		errorutil.Unauthorized(c, err.Error(), "you are not authorized")
		return
	}

	fmt.Println("USER ID", userID)

	// Create a new client for this connection
	client := hub.NewClient(*userID, conn)

	// Add the client to the repository
	huub.AddClient(client)

	// Start listening for incoming messages
	go client.Listen(huub)
}

// func (r *Hub) SendToClient(id string, messageType int, p []byte) error {
// 	r.Lock.RLock()
// 	defer r.Lock.RUnlock()

// 	if client, ok := r.Clients[id]; ok {
// 		return client.Conn.WriteMessage(messageType, p)
// 	}

// 	return fmt.Errorf("client not found: %s", id)
// }

// func handleWebSocketConnection(c *gin.Context) {
// 	fmt.Println("Connect")
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer conn.Close()

// 	for {
// 		// Get the current time and serialize it as a JSON message
// 		currentTime := time.Now().Format(time.RFC3339)
// 		message := []byte(`{"time":"` + currentTime + `"}`)

// 		// Send the JSON message to the client
// 		err := conn.WriteMessage(websocket.TextMessage, message)
// 		if err != nil {
// 			fmt.Println(err)
// 			break
// 		}

// 		// Wait for 10 seconds before sending the next message
// 		time.Sleep(10 * time.Second)
// 	}
// }
