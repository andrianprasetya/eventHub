package models

import (
	"gorm.io/gorm"
	"time"
)

// Notification model
type Notification struct {
	ID        string    `gorm:"type:varchar(50);primary_key:true" json:"id"`
	UserID    string    `gorm:"type:varchar(100);not null" json:"user_id"`
	Message   string    `gorm:"type:text;" json:"message"`
	IsRead    bool      `gorm:"type:boolean" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	gorm.Model
}
