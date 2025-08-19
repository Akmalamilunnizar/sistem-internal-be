package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"sistem-internal/database"
	"sistem-internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCustomers handles GET /api/customers with pagination and search
func GetCustomers(c *gin.Context) {
	var customers []models.Customer
	var total int64

	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := database.DB.Model(&models.Customer{})

	// Add search filter
	if search != "" {
		searchLower := strings.ToLower(search)
		query = query.Where("LOWER(name) LIKE ? OR LOWER(phone_number) LIKE ?", "%"+searchLower+"%", "%"+searchLower+"%")
	}

	// Get total count
	query.Count(&total)

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&customers).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  customers,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetCustomer handles GET /api/customers/:id
func GetCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer

	err := database.DB.First(&customer, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// CreateCustomer handles POST /api/customers
func CreateCustomer(c *gin.Context) {
	var customer models.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if customer.Name == "" || customer.Phone == "" || customer.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, phone, and email are required"})
		return
	}

	// Check if email already exists
	var existingCustomer models.Customer
	if err := database.DB.Where("email = ?", customer.Email).First(&existingCustomer).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Check if phone already exists
	if err := database.DB.Where("phone_number = ?", customer.Phone).First(&existingCustomer).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Phone number already exists"})
		return
	}

	// Set default password if not provided
	if customer.Password == "" {
		customer.Password = "default123" // You should hash this password
	}

	err := database.DB.Create(&customer).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": customer})
}

// UpdateCustomer handles PUT /api/customers/:id
func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer

	// Check if customer exists
	if err := database.DB.First(&customer, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customer"})
		return
	}

	// Bind update data
	var updateData models.Customer
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if updateData.Name != "" {
		customer.Name = updateData.Name
	}
	if updateData.Phone != "" {
		customer.Phone = updateData.Phone
	}
	if updateData.Email != "" {
		customer.Email = updateData.Email
	}
	if updateData.Address != "" {
		customer.Address = updateData.Address
	}
	if updateData.GPSLat != 0 {
		customer.GPSLat = updateData.GPSLat
	}
	if updateData.GPSLong != 0 {
		customer.GPSLong = updateData.GPSLong
	}
	if updateData.Status != "" {
		customer.Status = updateData.Status
	}

	err := database.DB.Save(&customer).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// DeleteCustomer handles DELETE /api/customers/:id
func DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer

	// Check if customer exists
	if err := database.DB.First(&customer, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customer"})
		return
	}

	err := database.DB.Delete(&customer).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

// GetCustomerTickets handles GET /api/customers/:id/tickets
func GetCustomerTickets(c *gin.Context) {
	customerID := c.Param("id")
	var tickets []models.TroubleTicket

	err := database.DB.Where("customer_id = ?", customerID).Preload("Customer").Find(&tickets).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customer tickets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets})
}
