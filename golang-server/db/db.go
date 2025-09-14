package db

import (
	"log"

	"github.com/example/golang-postgres-crud/config"
	"github.com/example/golang-postgres-crud/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(postgres.Open(config.DB_URL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	DB = database
	DB.AutoMigrate(&models.User{})
}
