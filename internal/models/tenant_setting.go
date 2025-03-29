package models

import (
	"gorm.io/gorm"
)

// Tenant Setting model
type TenantSetting struct {
	ID       string `gorm:"type:varchar(50);primary_key:true" json:"id"`
	TenantID string `gorm:"type:varchar(100);not null" json:"tenant_id"`
	Key      string `gorm:"type:varchar(100);not null" json:"key"`
	Value    string `gorm:"type:text" json:"value"`
	gorm.Model
}
