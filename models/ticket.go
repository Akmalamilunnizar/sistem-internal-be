package models

import "time"

// TroubleTicket represents a customer trouble workflow
type TroubleTicket struct {
	ID                  uint     `json:"id" gorm:"primaryKey"`
	CustomerID          uint     `json:"customer_id" gorm:"not null"`
	Customer            Customer `json:"customer" gorm:"foreignKey:CustomerID"`
	Title               string   `json:"title" gorm:"not null"`
	Type                string   `json:"type" gorm:"not null"` // e.g., "connection_issue", "billing", "technical"
	Description         string   `json:"description" gorm:"type:text"`
	Status              string   `json:"status" gorm:"not null;default:open"` // open, in_progress, resolved
	CurrentAssigneeRole string   `json:"current_assignee_role" gorm:"not null;default:customer_service"`

	// Optional relationships
	CustomerNote   string `json:"customer_note" gorm:"type:text"`
	NOCNote        string `json:"noc_note" gorm:"type:text"`
	TechnicianNote string `json:"technician_note" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
