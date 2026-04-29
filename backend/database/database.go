package database

import (
	"log"
	"time"

	"rental-backend/config"
	"rental-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	maxRetries := 10
	retryInterval := 3 * time.Second

	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(config.GetDBConnectionString()), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			log.Printf("Retrying in %v...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	if err != nil {
		return err
	}

	log.Println("Database connected successfully")

	if err = DB.AutoMigrate(
		&models.User{},
		&models.House{},
		&models.Favorite{},
		&models.Message{},
	); err != nil {
		return err
	}

	log.Println("Database migrated successfully")
	return nil
}
