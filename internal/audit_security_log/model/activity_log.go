package model

import (
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
	"time"
)

type ActivityLog struct {
	ID         string             `gorm:"type:varchar(50);primary_key:true"`
	TenantID   string             `gorm:"type:varchar(50);not null;index"`
	UserID     string             `gorm:"type:varchar(50);not null;index"`
	URL        string             `gorm:"type:varchar(100);not null"`
	Action     string             `gorm:"type:varchar(50);not null"`
	ObjectData string             `gorm:"type:text;not null"`
	ObjectType string             `gorm:"type:varchar(50);not null"`
	ObjectID   string             `gorm:"type:varchar(50);not null;index"`
	User       modelUser.User     `gorm:"foreignKey:UserID;references:ID"`
	Tenant     modelTenant.Tenant `gorm:"foreignKey:TenantID;references:ID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
