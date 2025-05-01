package model

import "time"

type EventTicket struct {
	ID         string `gorm:"type:varchar(50);primary_key:true"`
	EventID    string `gorm:"type:varchar(50);not null;index"`
	TicketType string `gorm:"type:varchar(100);not null;index"`
	Price      int    `gorm:"type:integer;default:0"`
	Quantity   int    `gorm:"type:integer;default:1"`
	Sold       int    `gorm:"type:integer;default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
