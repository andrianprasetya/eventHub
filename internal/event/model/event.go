package model

import (
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	"time"
)

type Event struct {
	ID          string             `gorm:"type:varchar(50);primary_key:true"`
	TenantID    string             `gorm:"type:varchar(50);not null;index"`
	title       string             `gorm:"type:varchar(100);not null;index"`
	Description string             `gorm:"type:text"`
	Location    string             `gorm:"type:varchar(255);not null"`
	StartDate   time.Time          `gorm:"type:date;not null"`
	EndDate     time.Time          `gorm:"type:date;not null"`
	Status      string             `gorm:"type:varchar(15);not null;comment: draft | published | cancelled"`
	Tenant      modelTenant.Tenant `gorm:"foreignKey:TenantID;references:ID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
