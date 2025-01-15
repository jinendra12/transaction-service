package database

import (
	"fmt"
	"transaction-service/internal/config"
	"transaction-service/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initialize(cfg *config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&model.Transaction{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return db, nil
}
