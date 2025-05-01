package model

import "time"

type ActivityLog struct {
	ID         string `gorm:"type:varchar(50);primary_key:true"`
	UserID     string `gorm:"type:varchar(50);not null;index"`
	Action     string `gorm:"type:varchar(100);not null"`
	ObjectType string `gorm:"type:varchar(100);not null"`
	ObjectID   string `gorm:"type:varchar(50);not null;index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
