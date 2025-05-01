package model

import "time"

type Export struct {
	ID        string `gorm:"type:varchar(50);primary_key:true"`
	TenantID  string `gorm:"type:varchar(50);not null;index"`
	FilePath  string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
