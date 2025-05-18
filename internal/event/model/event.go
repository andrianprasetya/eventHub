package model

import (
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
	"time"
)

type Event struct {
	ID          string             `gorm:"type:varchar(50);primary_key:true"`
	TenantID    string             `gorm:"type:varchar(50);not null;index"`
	Title       string             `gorm:"type:varchar(100);not null;index"`
	EventType   string             `gorm:"type:varchar(25);not null;comment: online || offline"`
	CategoryID  string             `gorm:"type:varchar(50);not null;index"`
	Tags        *utils.StringArray `gorm:"type:text[]"`
	Description *string            `gorm:"type:text"`
	Location    *string            `gorm:"type:varchar(255)"`
	StartDate   time.Time          `gorm:"type:date;not null"`
	EndDate     time.Time          `gorm:"type:date;not null"`
	CreatedBy   string             `gorm:"type:varchar(50);not null;index"`
	IsTicket    int                `gorm:"type:smallint;default:0;comment:0 => without ticket || 1 => with ticket"`
	Status      string             `gorm:"type:varchar(15);not null;comment: draft | published | cancelled"`
	Tenant      modelTenant.Tenant `gorm:"foreignKey:TenantID;references:ID"`
	User        modelUser.User     `gorm:"foreignKey:CreatedBy;references:ID"`
	Category    EventCategory      `gorm:"foreignKey:CategoryID;references:ID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
