package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"sistem-internal/database"
	"sistem-internal/models"
)

type createEmployeeRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleName string `json:"role"` // one of: owner, customer_service, noc, technician
}

type updateEmployeeRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleName string `json:"role"`
	Status   string `json:"status"`
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

// GetAllStaff handles GET /api/staff with pagination and search
func GetAllStaff(c *gin.Context) {
	var employees []models.Employee
	var total int64

	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := database.DB.Model(&models.Employee{}).Preload("Role")

	// Add search filter
	if search != "" {
		searchLower := strings.ToLower(search)
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", "%"+searchLower+"%", "%"+searchLower+"%")
	}

	// Get total count
	query.Count(&total)

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&employees).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch staff members"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  employees,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetStaffById handles GET /api/staff/:id
func GetStaffById(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	err := database.DB.Preload("Role").First(&employee, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Staff member not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch staff member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": employee})
}

// CreateEmployee allows Owner to create a new employee and assign a role
func CreateEmployee(c *gin.Context) {
	var req createEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" || req.Email == "" || req.Password == "" || req.RoleName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// Check if email already exists
	var existingEmployee models.Employee
	if err := database.DB.Where("email = ?", req.Email).First(&existingEmployee).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	roleID, ok := roleNameToID(req.RoleName)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	emp := models.Employee{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
		RoleID:   roleID,
		Status:   "active",
	}

	if err := database.DB.Create(&emp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create employee"})
		return
	}

	// Return employee without password
	emp.Password = ""
	c.JSON(http.StatusCreated, gin.H{"data": emp})
}

// UpdateStaff handles PUT /api/staff/:id
func UpdateStaff(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	// Check if employee exists
	if err := database.DB.First(&employee, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Staff member not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch staff member"})
		return
	}

	// Bind update data
	var req updateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if req.Name != "" {
		employee.Name = req.Name
	}
	if req.Email != "" {
		// Check if email already exists (excluding current employee)
		var existingEmployee models.Employee
		if err := database.DB.Where("email = ? AND id != ?", req.Email, id).First(&existingEmployee).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		employee.Email = req.Email
	}
	if req.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		employee.Password = string(hashed)
	}
	if req.RoleName != "" {
		roleID, ok := roleNameToID(req.RoleName)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
			return
		}
		employee.RoleID = roleID
	}
	if req.Status != "" {
		employee.Status = req.Status
	}

	err := database.DB.Save(&employee).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff member"})
		return
	}

	// Return employee without password
	employee.Password = ""
	c.JSON(http.StatusOK, gin.H{"data": employee})
}

// DeleteStaff handles DELETE /api/staff/:id
func DeleteStaff(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	// Check if employee exists
	if err := database.DB.First(&employee, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Staff member not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch staff member"})
		return
	}

	err := database.DB.Delete(&employee).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete staff member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff member deleted successfully"})
}

// GetRoles handles GET /api/roles
func GetRoles(c *gin.Context) {
	var roles []models.Role

	err := database.DB.Find(&roles).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// ListEmployees shows all employees (Owner only) - keeping for backward compatibility
func ListEmployees(c *gin.Context) {
	var employees []models.Employee
	if err := database.DB.Preload("Role").Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list employees"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"employees": employees})
}
