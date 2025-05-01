package model

import "time"

type Webhook struct {
	ID        string `gorm:"type:varchar(50);primary_key:true"`
	TenantID  string `gorm:"type:varchar(50);not null;index"`
	Url       string `gorm:"type:varchar(255);not null"`
	Event     string `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
