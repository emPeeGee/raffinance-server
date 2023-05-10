package connection

import (
	"fmt"
	"os"

	"github.com/emPeeGee/raffinance/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg config.DB) (*gorm.DB, error) {

	// TODO: BAD
	dbPort := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", dbHost, dbPort, username, dbName, password)
	fmt.Println(dsn)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		return nil, err

	}

	return gormDB, nil
}
