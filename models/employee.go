package models

import (
	"time"
)

type Employee struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	RoleID    uint      `json:"role_id" gorm:"not null"`
	Role      Role      `json:"role" gorm:"foreignKey:RoleID"`
	Status    string    `json:"status" gorm:"not null;default:active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
