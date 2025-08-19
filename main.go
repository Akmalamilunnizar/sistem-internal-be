package main

import (
	"log"
	"net/http"

	"sistem-internal/database"
	"sistem-internal/handlers"
	"sistem-internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	database.ConnectDB()

	// Create a new Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:3000", // Nuxt.js default port
		"http://localhost:3001", // Alternative Nuxt.js port
		"http://localhost:3002", // Your working frontend port
		"http://localhost:3003", // Alternative Nuxt.js port
		"http://localhost:3005", // Another alternative port
		"http://127.0.0.1:3000",
		"http://127.0.0.1:3001",
		"http://127.0.0.1:3002",
		"http://127.0.0.1:3003",
		"http://127.0.0.1:3005",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Define routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Sistem Internal API",
			"status":  "running",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "healthy",
			})
		})

		// User routes
		api.GET("/users", handlers.GetUsers)
		api.GET("/users/:id", handlers.GetUser)
		api.POST("/users", handlers.CreateUser)
		api.PUT("/users/:id", handlers.UpdateUser)
		api.DELETE("/users/:id", handlers.DeleteUser)
		api.GET("/users/count", handlers.GetUserCount)

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/employee/login", handlers.EmployeeLogin)
			auth.POST("/customer/login", handlers.CustomerLogin)
		}

		// Tickets workflow routes (RBAC)
		tickets := api.Group("/tickets")
		{
			// Owner can view all tickets
			tickets.GET("", middleware.AuthRequired(), middleware.RoleRequired("owner"), handlers.ListTickets)

			// Customer Service creates ticket and forwards to NOC
			tickets.POST("", middleware.AuthRequired(), middleware.RoleRequired("customer_service"), handlers.CreateTicket)
			tickets.POST(":id/forward/noc", middleware.AuthRequired(), middleware.RoleRequired("customer_service"), handlers.ForwardToNOC)

			// NOC diagnoses
			tickets.POST(":id/noc/diagnose", middleware.AuthRequired(), middleware.RoleRequired("noc"), handlers.NOCDiagnose)

			// Technician resolves
			tickets.POST(":id/technician/resolve", middleware.AuthRequired(), middleware.RoleRequired("technician"), handlers.TechnicianResolve)
		}

		// Customer management routes
		customers := api.Group("/customers")
		{
			customers.GET("", middleware.AuthRequired(), handlers.GetCustomers)
			customers.GET("/:id", middleware.AuthRequired(), handlers.GetCustomer)
			customers.POST("", middleware.AuthRequired(), middleware.RoleRequired("customer_service"), handlers.CreateCustomer)
			customers.PUT("/:id", middleware.AuthRequired(), middleware.RoleRequired("customer_service"), handlers.UpdateCustomer)
			customers.DELETE("/:id", middleware.AuthRequired(), middleware.RoleRequired("owner"), handlers.DeleteCustomer)
			customers.GET("/:id/tickets", middleware.AuthRequired(), handlers.GetCustomerTickets)
		}

		// Reports and analytics routes
		reports := api.Group("/reports")
		{
			reports.GET("/tickets", middleware.AuthRequired(), handlers.GetAllTickets)
			reports.GET("/trouble-types", middleware.AuthRequired(), handlers.GetTroubleTypeStats)
			reports.GET("/geographic", middleware.AuthRequired(), handlers.GetGeographicTroubleData)
			reports.GET("/summary", middleware.AuthRequired(), handlers.GetTroubleSummary)
		}

		// Enhanced user management routes (Owner only)
		staff := api.Group("/staff")
		{
			staff.GET("", middleware.AuthRequired(), middleware.RoleRequired("owner"), handlers.GetAllStaff)
			staff.GET("/:id", middleware.AuthRequired(), middleware.RoleRequired("owner"), handlers.GetStaffById)
			staff.POST("", middleware.AuthRequired(), middleware.RoleRequired("owner"), handlers.CreateEmployee)
			staff.PUT("/:id", middleware.AuthRequired(), middleware.RoleRequired("owner"), handlers.UpdateStaff)
			staff.DELETE("/:id", middleware.AuthRequired(), middleware.RoleRequired("owner"), handlers.DeleteStaff)
		}

		// Roles management
		api.GET("/roles", middleware.AuthRequired(), middleware.RoleRequired("owner"), handlers.GetRoles)

		// Owner employee management (keeping for backward compatibility)
		employees := api.Group("/employees")
		{
			employees.GET("", middleware.AuthRequired(), middleware.RoleRequired("owner"), handlers.ListEmployees)
			employees.POST("", middleware.AuthRequired(), middleware.RoleRequired("owner"), handlers.CreateEmployee)
		}
	}

	// Start server
	log.Println("Server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
