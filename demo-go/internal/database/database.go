package database

import (
	"demo-go/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) *gorm.DB {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    if err := db.AutoMigrate(&models.Order{}, &models.Product{}); err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    return db
}
