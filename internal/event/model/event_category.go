package model

import "time"

type EventCategory struct {
	ID        string `gorm:"type:varchar(50);primary_key:true"`
	TenantID  string `gorm:"type:varchar(50);not null;index"`
	Name      string `gorm:"type:varchar(100);not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
