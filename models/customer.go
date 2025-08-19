package models

import (
	"time"
)

type Customer struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Phone     string    `json:"phone_number" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Address   string    `json:"address" gorm:"type:text"`
	GPSLat    float64   `json:"gps_lat" gorm:"type:decimal(10,8)"`
	GPSLong   float64   `json:"gps_long" gorm:"type:decimal(11,8)"`
	Password  string    `json:"-" gorm:"not null"`
	Status    string    `json:"status" gorm:"not null;default:active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
