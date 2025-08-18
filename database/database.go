package database

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
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
	err = DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Employee{}, &models.Customer{}, &models.TroubleTicket{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migrated successfully!")

	// Seed initial data if table is empty
	seedInitialData()
}

func seedInitialData() {
	// Seed roles
	var roleCount int64
	DB.Model(&models.Role{}).Count(&roleCount)

	if roleCount == 0 {
		roles := []models.Role{
			{Name: "owner", Description: "Owner"},
			{Name: "customer_service", Description: "Customer Service"},
			{Name: "noc", Description: "Network Operations Center"},
			{Name: "technician", Description: "Field Technician"},
		}

		for _, role := range roles {
			DB.Create(&role)
		}
		log.Println("Roles seeded successfully!")
	}

	// Seed users
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)

	if userCount == 0 {
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
		log.Println("Users seeded successfully!")
	}

	// Seed employees
	var employeeCount int64
	DB.Model(&models.Employee{}).Count(&employeeCount)

	if employeeCount == 0 {
		// Hash password: "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		employees := []models.Employee{
			{Name: "Owner User", Email: "admin@example.com", Password: string(hashedPassword), RoleID: 1},
			{Name: "Customer Service", Email: "cs@example.com", Password: string(hashedPassword), RoleID: 2},
			{Name: "NOC Operator", Email: "noc@example.com", Password: string(hashedPassword), RoleID: 3},
			{Name: "Technician User", Email: "tech@example.com", Password: string(hashedPassword), RoleID: 4},
		}

		for _, employee := range employees {
			DB.Create(&employee)
		}
		log.Println("Employees seeded successfully!")
	}

	// Seed customers
	var customerCount int64
	DB.Model(&models.Customer{}).Count(&customerCount)

	if customerCount == 0 {
		// Hash password: "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		customers := []models.Customer{
			{
				Name:     "Customer One",
				Phone:    "081234567890",
				Email:    "customer1@example.com",
				Password: string(hashedPassword),
				Status:   "active",
			},
			{
				Name:     "Customer Two",
				Phone:    "081234567891",
				Email:    "customer2@example.com",
				Password: string(hashedPassword),
				Status:   "active",
			},
		}

		for _, customer := range customers {
			DB.Create(&customer)
		}
		log.Println("Customers seeded successfully!")
	}
}
