package model

import "time"

type EventSession struct {
	ID        string    `gorm:"type:varchar(50);primary_key:true"`
	EventID   string    `gorm:"type:varchar(50);not null"`
	Title    string    `gorm:"type:varchar(255);not null"`
	StartDate time.Time `gorm:"type:date;not null"`
	EndDate   time.Time `gorm:"type:date;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
