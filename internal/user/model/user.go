package model

import "time"

type User struct {
	ID        string `gorm:"type:varchar(50);primary_key:true"`
	TenantID  string `gorm:"type:varchar(50);not null"`
	Name      string `gorm:"type:varchar(100);not null;index"`
	Email     string `gorm:"type:varchar(50);not null;unique"`
	Password  string `gorm:"type:varchar(100)"`
	IsActive  int    `gorm:"type:smallint;default:0;comment: 0 => in-active | 1 => active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
