package model

import "time"

type Tenant struct {
	ID        string `gorm:"type:varchar(50);primary_key:true"`
	Name      string `gorm:"type:varchar(100);not null;index"`
	Email     string `gorm:"type:varchar(50);not null;unique"`
	LogoUrl   string `gorm:"type:text"`
	IsActive  int    `gorm:"type:smallint;default:0;comment: 0 => in-active | 1 => active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
