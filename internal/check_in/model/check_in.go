package model

import "time"

type CheckIn struct {
	ID        string    `gorm:"type:varchar(50);primary_key:true"`
	TicketID  string    `gorm:"type:varchar(50);not null;index"`
	DueDate   time.Time `gorm:"type:timestamp;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
