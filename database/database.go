package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"sistem-internal/config"
	"sistem-internal/models"
)

var DB *gorm.DB

func ConnectDB() {
	// Database configuration
	dbHost := config.DB_HOST
	dbPort := config.DB_PORT
	dbUser := config.DB_USER
	dbPassword := config.DB_PASSWORD
	dbName := config.DB_NAME

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Info)

	// Connect to database
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully!")

	// Auto migrate the schema
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migrated successfully!")

	// Seed initial data if table is empty
	seedInitialData()
}

func seedInitialData() {
	var count int64
	DB.Model(&models.User{}).Count(&count)

	if count == 0 {
		users := []models.User{
			{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			{
				Name:  "Jane Smith",
				Email: "jane@example.com",
			},
		}

		for _, user := range users {
			DB.Create(&user)
		}

		log.Println("Initial data seeded successfully!")
	}
}
