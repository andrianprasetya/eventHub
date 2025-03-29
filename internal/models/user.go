package models

import (
	"gorm.io/gorm"
	"time"
)

// User model
type User struct {
	ID        string    `gorm:"type:varchar(50);primary_key:true" json:"id"`
	TenantID  string    `gorm:" type:varchar(100);not null;index" json:"tenant_id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Role      string    `gorm:"type:varchar(100);not null" json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	gorm.Model
}
