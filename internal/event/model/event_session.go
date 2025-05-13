package model

import "time"

type EventSession struct {
	ID            string    `gorm:"type:varchar(50);primary_key:true"`
	EventID       string    `gorm:"type:varchar(50);not null"`
	Title         string    `gorm:"type:varchar(255);not null"`
	StartDateTime time.Time `gorm:"type:timestamp;not null"`
	EndDateTime   time.Time `gorm:"type:timestamp;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
