package handlers

import (
	"net/http"
	"strconv"

	"sistem-internal/database"
	"sistem-internal/models"

	"github.com/gin-gonic/gin"
)

// GetAllTickets handles GET /api/tickets with pagination
func GetAllTickets(c *gin.Context) {
	var tickets []models.TroubleTicket
	var total int64

	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := database.DB.Model(&models.TroubleTicket{}).Preload("Customer")

	// Add search filter if provided
	if search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	query.Count(&total)

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&tickets).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tickets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tickets,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetTroubleTypeStats handles GET /api/reports/trouble-types
func GetTroubleTypeStats(c *gin.Context) {
	type TroubleTypeStat struct {
		Name       string  `json:"name"`
		Count      int64   `json:"count"`
		Percentage float64 `json:"percentage"`
	}

	var stats []TroubleTypeStat
	var total int64

	// Get total count
	database.DB.Model(&models.TroubleTicket{}).Count(&total)

	if total > 0 {
		// Get stats by type
		rows, err := database.DB.Model(&models.TroubleTicket{}).
			Select("type as name, count(*) as count").
			Group("type").
			Rows()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trouble type stats"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var stat TroubleTypeStat
			rows.Scan(&stat.Name, &stat.Count)
			stat.Percentage = float64(stat.Count) / float64(total) * 100
			stats = append(stats, stat)
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetGeographicTroubleData handles GET /api/reports/geographic
func GetGeographicTroubleData(c *gin.Context) {
	type GeographicData struct {
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		TicketCount int64   `json:"ticket_count"`
		Area        string  `json:"area"`
	}

	var data []GeographicData

	// Get geographic trouble data
	rows, err := database.DB.Table("trouble_tickets").
		Select("customers.gps_lat as latitude, customers.gps_long as longitude, count(*) as ticket_count").
		Joins("JOIN customers ON trouble_tickets.customer_id = customers.id").
		Where("customers.gps_lat IS NOT NULL AND customers.gps_long IS NOT NULL").
		Group("customers.gps_lat, customers.gps_long").
		Order("ticket_count DESC").
		Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch geographic data"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item GeographicData
		rows.Scan(&item.Latitude, &item.Longitude, &item.TicketCount)
		item.Area = "Area " + strconv.FormatFloat(item.Latitude, 'f', 2, 64) + ", " + strconv.FormatFloat(item.Longitude, 'f', 2, 64)
		data = append(data, item)
	}

	// Find most affected area
	var mostAffectedArea string
	if len(data) > 0 {
		mostAffectedArea = data[0].Area
	}

	c.JSON(http.StatusOK, gin.H{
		"data":             data,
		"mostAffectedArea": mostAffectedArea,
	})
}

// GetTroubleSummary handles GET /api/reports/summary
func GetTroubleSummary(c *gin.Context) {
	var totalTickets, resolvedTickets, inProgressTickets, openTickets int64

	// Get total tickets
	database.DB.Model(&models.TroubleTicket{}).Count(&totalTickets)

	// Get resolved tickets
	database.DB.Model(&models.TroubleTicket{}).Where("status = ?", "resolved").Count(&resolvedTickets)

	// Get in progress tickets
	database.DB.Model(&models.TroubleTicket{}).Where("status = ?", "in_progress").Count(&inProgressTickets)

	// Get open tickets
	database.DB.Model(&models.TroubleTicket{}).Where("status = ?", "open").Count(&openTickets)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"totalTickets":      totalTickets,
			"resolvedTickets":   resolvedTickets,
			"inProgressTickets": inProgressTickets,
			"openTickets":       openTickets,
		},
	})
}
