package models

import (
	"gorm.io/gorm"
	"time"
)

// Password Reset model
type PasswordReset struct {
	ID        string    `gorm:"type:varchar(50);primary_key:true" json:"id"`
	Token     string    `gorm:"type:varchar(255);not null" json:"token"`
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	gorm.Model
}
