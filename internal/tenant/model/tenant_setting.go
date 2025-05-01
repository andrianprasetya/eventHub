package model

import "time"

type TenantSetting struct {
	ID        string `gorm:"type:varchar(50);primary_key:true"`
	TenantID  string `gorm:"type:varchar(50);not null;index"`
	Key       string `gorm:"type:varchar(50);not null;index"`
	Value     string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
