package models

import (
	"gorm.io/gorm"
	"time"
)

// Activity Log model
type ActivityLog struct {
	ID         string    `gorm:"type:varchar(50);primary_key:true" json:"id"`
	UserID     string    `gorm:"type:varchar(100);not null" json:"user_id"`
	Action     string    `gorm:"type:varchar(25);not null" json:"action"`
	ObjectID   float64   `gorm:"type:varchar(100);not null" json:"object_id"`
	ObjectType int       `gorm:"type:varchar(100);not null" json:"object_type"`
	Timestamp  time.Time `gorm:"autoCreateTime"json:"change_at"`
	gorm.Model
}
