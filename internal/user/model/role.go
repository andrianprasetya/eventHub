package model

import "time"

type Role struct {
	ID          string `gorm:"type:varchar(50);primary_key:true"`
	Name        string `gorm:"type:varchar(100);not null;index"`
	Slug        string `gorm:"type:varchar(100);not null;unique"`
	Description string `gorm:"type:text;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RoleChannel struct {
	Role *Role
	Err  error
}
