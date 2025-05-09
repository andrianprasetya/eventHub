package model

import (
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"

	"time"
)

type User struct {
	ID        string             `gorm:"type:varchar(50);primary_key:true"`
	TenantID  string             `gorm:"type:varchar(50);default:null"`
	RoleID    string             `gorm:"type:varchar(50);not null"`
	Name      string             `gorm:"type:varchar(100);not null;index"`
	Email     string             `gorm:"type:varchar(50);not null;unique"`
	Password  string             `gorm:"type:varchar(100)"`
	IsActive  int                `gorm:"type:smallint;default:0;comment: 0 => in-active | 1 => active"`
	Tenant    modelTenant.Tenant `gorm:"foreignKey:TenantID;references:ID"`
	Role      Role               `gorm:"foreignKey:RoleID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
