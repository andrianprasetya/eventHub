package models

import (
	"gorm.io/gorm"
	"time"
)

// Tenant model
type Tenant struct {
	ID               string    `gorm:"type:varchar(50);primary_key:true" json:"id"`
	Name             string    `gorm:"type:varchar(100);not null" json:"name"`
	Email            string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	SubscriptionPlan string    `gorm:"type:varchar(50)" json:"subscription_plan"`
	Logo             string    `gorm:"type:varchar(255)" json:"logo"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	gorm.Model
}
