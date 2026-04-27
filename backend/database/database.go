package database

import (
	"log"

	"rental-backend/config"
	"rental-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(postgres.Open(config.GetDBConnectionString()), &gorm.Config{})
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
