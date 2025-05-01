package model

import "time"

type ApiKey struct {
	ID          string `gorm:"type:varchar(50);primary_key:true"`
	TenantID    string `gorm:"type:varchar(50);not null;index"`
	ApiKey      string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Status      string `gorm:"type:varchar(15);not null;comment: inactive | active"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
