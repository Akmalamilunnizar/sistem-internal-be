package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sistem-internal/database"
	"sistem-internal/models"
)

// Create ticket by Customer Service on behalf of a customer
func CreateTicket(c *gin.Context) {
	var req struct {
		CustomerID  uint   `json:"customer_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.CustomerID == 0 || req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	ticket := models.TroubleTicket{
		CustomerID:          req.CustomerID,
		Title:               req.Title,
		Description:         req.Description,
		Status:              "received",
		CurrentAssigneeRole: "customer_service",
	}
	if err := database.DB.Create(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ticket"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"ticket": ticket})
}

// Forward to NOC by Customer Service
func ForwardToNOC(c *gin.Context) {
	id := c.Param("id")
	var ticket models.TroubleTicket
	if err := database.DB.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}
	ticket.Status = "forwarded_to_noc"
	ticket.CurrentAssigneeRole = "noc"
	if err := database.DB.Save(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to forward"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ticket": ticket})
}

// NOC diagnosis and resolution/forward decision
func NOCDiagnose(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Note              string `json:"note"`
		IsPhysicalProblem bool   `json:"is_physical_problem"`
		ResolvedByNOC     bool   `json:"resolved_by_noc"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	var ticket models.TroubleTicket
	if err := database.DB.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}
	ticket.NOCNote = req.Note
	if req.ResolvedByNOC {
		ticket.Status = "resolved"
		ticket.CurrentAssigneeRole = "customer_service"
	} else if req.IsPhysicalProblem {
		ticket.Status = "forwarded_to_technician"
		ticket.CurrentAssigneeRole = "technician"
	} else {
		ticket.Status = "diagnosed"
		ticket.CurrentAssigneeRole = "customer_service"
	}
	if err := database.DB.Save(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ticket": ticket})
}

// Technician resolves
func TechnicianResolve(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Note string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	var ticket models.TroubleTicket
	if err := database.DB.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}
	ticket.TechnicianNote = req.Note
	ticket.Status = "resolved"
	ticket.CurrentAssigneeRole = "customer_service"
	if err := database.DB.Save(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ticket": ticket})
}

// Basic list endpoint for convenience
func ListTickets(c *gin.Context) {
	var tickets []models.TroubleTicket
	if err := database.DB.Order("id desc").Find(&tickets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tickets": tickets})
}
