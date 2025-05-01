package model

import "time"

type Subscription struct {
	ID        string    `gorm:"type:varchar(50);primary_key:true"`
	TenantID  string    `gorm:"type:varchar(50);not null"`
	PlanID    string    `gorm:"type:varchar(50);not null"`
	StartDate time.Time `gorm:"type:date;not null"`
	EndDate   time.Time `gorm:"type:date;not null"`
	IsActive  int       `gorm:"type:smallint;default:0;comment: 0 => in-active | 1 => active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
