package model

import "time"

type EventCustomField struct {
	ID        string `gorm:"type:varchar(50);primary_key:true"`
	TicketID  string `gorm:"type:varchar(50);not null;index"`
	FieldName string `gorm:"type:varchar(100);not null"`
	FieldType string `gorm:"type:varchar(20);not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
