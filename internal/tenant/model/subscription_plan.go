package model

import (
	"gorm.io/gorm"
	"time"
)

type SubscriptionPlan struct {
	ID          string `gorm:"type:varchar(50);primary_key:true"`
	Name        string `gorm:"type:varchar(100);not null;index"`
	Price       int    `gorm:"type:integer;default:0"`
	Feature     string `gorm:"type:text;not null"`
	DurationDay int    `gorm:"type:integer;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
