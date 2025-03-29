package models

import (
	"gorm.io/gorm"
	"time"
)

// Event model
type Event struct {
	ID          string    `gorm:"type:varchar(50);primary_key:true" json:"id"`
	TenantID    string    `gorm:"type:varchar(100);not null" json:"tenant_id"`
	Title       string    `gorm:"type:varchar(100);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Location    string    `gorm:"type:varchar(255)" json:"location"`
	StartTime   time.Time `gorm:"type:timestamp" json:"start_time"`
	EndTime     time.Time `gorm:"type:timestamp" json:"end_time"`
	Status      string    `gorm:"type:varchar(255)" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	gorm.Model
}
