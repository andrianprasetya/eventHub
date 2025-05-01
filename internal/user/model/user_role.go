package model

import "time"

type UserRole struct {
	ID        string `gorm:"type:varchar(50);primary_key:true"`
	UserID    string `gorm:"type:varchar(100);not null;index"`
	RoleID    string `gorm:"type:varchar(100);not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
