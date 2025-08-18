package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"sistem-internal/database"
	"sistem-internal/models"
)

type createEmployeeRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleName string `json:"role"` // one of: owner, customer_service, noc, technician
}

func roleNameToID(roleName string) (uint, bool) {
	switch roleName {
	case "owner":
		return 1, true
	case "customer_service":
		return 2, true
	case "noc":
		return 3, true
	case "technician":
		return 4, true
	default:
		return 0, false
	}
}

// CreateEmployee allows Owner to create a new employee and assign a role
func CreateEmployee(c *gin.Context) {
	var req createEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" || req.Email == "" || req.Password == "" || req.RoleName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	roleID, ok := roleNameToID(req.RoleName)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	emp := models.Employee{
		Name: req.Name, Email: req.Email, Password: string(hashed), RoleID: roleID,
	}
	if err := database.DB.Create(&emp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create employee"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"employee": gin.H{"id": emp.ID, "name": emp.Name, "email": emp.Email, "role_id": emp.RoleID}})
}

// ListEmployees shows all employees (Owner only)
func ListEmployees(c *gin.Context) {
	var employees []models.Employee
	if err := database.DB.Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list employees"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"employees": employees})
}
