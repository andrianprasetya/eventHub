package model

import "time"

type ApiKey struct {
	ID          string `gorm:"type:varchar(50);primary_key:true"`
	TenantID    string `gorm:"type:varchar(50);not null;index"`
	ApiKey      string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Status      int    `gorm:"type:smallint;default:0;comment: 0 => in-active | 1 => active"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
