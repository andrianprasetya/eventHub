package model

import "time"

type EmailTemplate struct {
	ID        string `gorm:"type:varchar(50);primary_key:true"`
	Name      string `gorm:"type:varchar(100);not null"`
	Subject   string `gorm:"type:varchar(100);not null"`
	Body      string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
