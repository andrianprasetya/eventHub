package model

import "time"

type Report struct {
	ID         string `gorm:"type:varchar(50);primary_key:true"`
	TenantID   string `gorm:"type:varchar(50);not null;index"`
	Name       string `gorm:"type:varchar(100);not null"`
	ReportData string `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
