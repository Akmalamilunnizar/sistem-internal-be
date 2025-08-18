package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"sistem-internal/database"
	"sistem-internal/models"
)

type employeeLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type customerLoginRequest struct {
	EmailOrPhone string `json:"email"` // UI may send phone in email field
	Phone        string `json:"phone"` // optional explicit phone field
	Password     string `json:"password"`
}

// EmployeeLogin authenticates an employee by email and password
func EmployeeLogin(c *gin.Context) {
	var req employeeLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	var employee models.Employee
	if err := database.DB.Where("email = ?", req.Email).First(&employee).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// In a real app, generate JWT. For now, return a placeholder token
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token": "employee-token",
			"user": gin.H{
				"role": gin.H{"id": employee.RoleID},
			},
		},
	})
}

// CustomerLogin authenticates a customer by email or phone and password
func CustomerLogin(c *gin.Context) {
	var req customerLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "credentials are required"})
		return
	}

	identifier := req.EmailOrPhone
	if identifier == "" {
		identifier = req.Phone
	}
	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or phone is required"})
		return
	}

	var customer models.Customer
	if err := database.DB.Where("email = ? OR phone = ?", identifier, identifier).First(&customer).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token": "customer-token",
		},
	})
}
