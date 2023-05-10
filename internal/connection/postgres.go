package connection

import (
	"fmt"

	"github.com/emPeeGee/raffinance/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg config.DB) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Name, cfg.Password)
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
