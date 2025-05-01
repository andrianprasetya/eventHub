package model

import "time"

type Permission struct {
	ID          string `gorm:"type:varchar(50);primary_key:true"`
	Name        string `gorm:"type:varchar(100);not null;index"`
	Description string `gorm:"type:text;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
