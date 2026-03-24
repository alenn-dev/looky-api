package database

import (
	"fmt"
	"log"
	"looky/internal/config"
	"looky/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := config.Env.DatabaseURL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the DB:", err)
	}

	fmt.Println("Successful connection")

	// Auto migrations
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Order{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.OrderItem{})
	db.AutoMigrate(&models.Restaurant{})
	db.AutoMigrate(&models.OrderStatusHistory{})

	DB = db
}
