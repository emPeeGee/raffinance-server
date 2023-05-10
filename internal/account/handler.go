package account

import (
	"net/http"
	"strconv"
	"time"

	"github.com/emPeeGee/raffinance/internal/auth"
	"github.com/emPeeGee/raffinance/pkg/errorutil"
	"github.com/emPeeGee/raffinance/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func RegisterHandlers(apiRg *gin.RouterGroup, service Service, validate *validator.Validate, logger log.Logger) {
	h := handler{service, logger, validate}

	api := apiRg.Group("/accounts")
	{
		api.POST("", h.createAccount)
		api.PUT("/:id", h.updateAccount)
		api.DELETE("/:id", h.deleteAccount)
		api.GET("/:id/bal", h.balance)

		api.GET("", h.getAccounts)
		api.GET("/:id", h.getAccount)
		// api.GET("/:id/transactions", h.getAccountTransactionsByMonth)
	}
}

type handler struct {
	service  Service
	logger   log.Logger
	validate *validator.Validate
}

func (h *handler) createAccount(c *gin.Context) {
	var input createAccountDTO

	userId, err := auth.GetUserId(c)
	if err != nil || userId == nil {
		errorutil.Unauthorized(c, err.Error(), "you are not authorized")
		return
	}

	if err := c.BindJSON(&input); err != nil {
		errorutil.BadRequest(c, "your request looks incorrect", err.Error())
		return
	}

	if err := h.validate.Struct(input); err != nil {
		errorutil.BadRequest(c, "your request did not pass validation", err.Error())
		return
	}

	createdAccount, err := h.service.createAccount(*userId, input)
	if err != nil {
		errorutil.InternalServer(c, "It looks like name is already used", err.Error())
		return
	}

	c.JSON(http.StatusOK, createdAccount)
}

func (h *handler) updateAccount(c *gin.Context) {
	var input updateAccountDTO

	userId, err := auth.GetUserId(c)
	if err != nil || userId == nil {
		errorutil.Unauthorized(c, err.Error(), "you are not authorized")
		return
	}

	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errorutil.BadRequest(c, "wrong account id", err.Error())
		return
	}

	if err := c.BindJSON(&input); err != nil {
		errorutil.BadRequest(c, "your request looks incorrect", err.Error())
		return
	}

	if err := h.validate.Struct(input); err != nil {
		errorutil.BadRequest(c, "your request did not pass validation", err.Error())
		return
	}

	updatedAccount, err := h.service.updateAccount(*userId, uint(accountId), input)
	if err != nil {
		errorutil.BadRequest(c, "error", err.Error())
		return
	}

	c.JSON(http.StatusOK, updatedAccount)
}

func (h *handler) deleteAccount(c *gin.Context) {
	userId, err := auth.GetUserId(c)
	if err != nil || userId == nil {
		errorutil.Unauthorized(c, err.Error(), "you are not authorized")
		return
	}

	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errorutil.BadRequest(c, err.Error(), "the id must be an integer")
		return
	}

	if err := h.service.deleteAccount(*userId, uint(accountId)); err != nil {
		h.logger.Info(err.Error())
		errorutil.NotFound(c, err.Error(), "Not found")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}

func (h *handler) getAccounts(c *gin.Context) {
	userId, err := auth.GetUserId(c)
	if err != nil || userId == nil {
		errorutil.Unauthorized(c, err.Error(), "you are not authorized")
		return
	}

	accounts, err := h.service.getAccounts(*userId)
	if err != nil {
		errorutil.InternalServer(c, "something went wrong, we are working", err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h *handler) getAccount(c *gin.Context) {
	userId, err := auth.GetUserId(c)
	if err != nil || userId == nil {
		errorutil.Unauthorized(c, err.Error(), "you are not authorized")
		return
	}

	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errorutil.BadRequest(c, err.Error(), "the id must be an integer")
		return
	}

	includeTxnsQuery := c.DefaultQuery("include_transactions", "true")
	includeTxns, err := strconv.ParseBool(includeTxnsQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for include_transactions parameter"})
		return
	}

	if includeTxns {
		accounts, err := h.service.getAccountWithTransactions(*userId, uint(accountId))
		if err != nil {
			errorutil.InternalServer(c, "something went wrong, we are working", err.Error())
			return
		}

		c.JSON(http.StatusOK, accounts)
	} else {
		accounts, err := h.service.getAccount(*userId, uint(accountId))
		if err != nil {
			errorutil.InternalServer(c, "something went wrong, we are working", err.Error())
			return
		}

		c.JSON(http.StatusOK, accounts)
	}
}

func (h *handler) getAccountTransactionsByMonth(c *gin.Context) {
	accountIdStr := c.Param("id")
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	accountId, err := strconv.ParseUint(accountIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
		return
	}

	transactions, err := h.service.getAccountTransactionsByMonth(uint(accountId), year, time.Month(month))
	if err != nil {
		errorutil.InternalServer(c, "Failed to get transactions", err.Error())
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// TODO: The func is useless now. Just for test purpose it exists
func (h *handler) balance(c *gin.Context) {
	userId, err := auth.GetUserId(c)
	if err != nil || userId == nil {
		errorutil.Unauthorized(c, err.Error(), "you are not authorized")
		return
	}

	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errorutil.BadRequest(c, err.Error(), "the id must be an integer")
		return
	}

	bal, err := h.service.getAccountBalance(*userId, uint(accountId))
	if err != nil {
		h.logger.Info(err.Error())
		errorutil.NotFound(c, err.Error(), "Not found")
		return
	}

	userBal, err := h.service.getUserBalance(*userId)
	if err != nil {
		h.logger.Info(err.Error())
		errorutil.NotFound(c, err.Error(), "Not found")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accountBal": bal,
		"userBal":    userBal,
	})
}
