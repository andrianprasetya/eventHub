package models

import (
	"gorm.io/gorm"
)

// Webhook model
type Webhook struct {
	ID       string `gorm:"type:varchar(50);primary_key:true" json:"id"`
	TenantID string `gorm:"type:varchar(100);not null" json:"tenant_id"`
	Url      string `gorm:"type:varchar(100)" json:"url"`
	Event    string `gorm:"type:varchar(100)" json:"event"`
	gorm.Model
}
