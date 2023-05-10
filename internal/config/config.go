package config

import (
	"os"
	"time"

	"github.com/emPeeGee/raffinance/pkg/log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	defaultMaxHeaderBytes = 1 << 20 // 1 MB
	defaultReadTimeout    = 10 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	path                  = "configs"
	fileName              = "config"
)

type Config struct {
	Server
	DB
}

type Server struct {
	Addr           string
	MaxHeaderBytes int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

type DB struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}

func Get(logger log.Logger) (*Config, error) {
	if err := initializeConfig(); err != nil {
		logger.Fatalf("Error initializing config: %s", err.Error())
		return nil, err
	}

	if err := godotenv.Load(); err != nil {
		logger.Fatalf("Error loading env variables: %s", err.Error())
		return nil, err
	}

	db := DB{
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
	}

	port := os.Getenv("PORT") // Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "9000" // localhost
	}

	server := Server{
		Addr:           ":" + port,
		MaxHeaderBytes: defaultMaxHeaderBytes,
		ReadTimeout:    defaultReadTimeout,
		WriteTimeout:   defaultWriteTimeout,
	}

	return &Config{server, db}, nil

}

// TODO:
// Deprecated: Viper is no longer used
func initializeConfig() error {
	viper.AddConfigPath(path)
	viper.SetConfigName(fileName)

	return viper.ReadInConfig()
}
