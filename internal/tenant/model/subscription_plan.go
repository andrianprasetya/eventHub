package model

import "time"

type SubscriptionPlan struct {
	ID          string `gorm:"type:varchar(50);primary_key:true"`
	Name        string `gorm:"type:varchar(100);not null;index"`
	Price       int    `gorm:"type:integer;default:0"`
	Feature     string `gorm:"type:text"`
	DurationDay int    `gorm:"type:integer;default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
