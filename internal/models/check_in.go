package models

import (
	"gorm.io/gorm"
	"time"
)

// CheckIn model
type CheckIn struct {
	ID          string    `gorm:"type:varchar(50);primary_key:true" json:"id"`
	TicketID    string    `gorm:"type:varchar(100);not null" json:"ticket_id"`
	CheckInTime time.Time `gorm:"type:timestamp" json:"check_in_time"`
	gorm.Model
}
