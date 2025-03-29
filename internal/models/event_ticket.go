package models

import (
	"gorm.io/gorm"
)

// Event Ticket model
type EventTicket struct {
	ID         string  `gorm:"type:varchar(50);primary_key:true" json:"id"`
	EventID    string  `gorm:"type:varchar(100);not null" json:"event_id"`
	TicketType string  `gorm:"type:varchar(25);not null" json:"ticket_type"`
	Price      float64 `gorm:"type:decimal(10,2)" json:"price"`
	Quantity   int     `gorm:"type:integer" json:"quantity"`
	Sold       int     `gorm:"type:integer" json:"sold"`
	gorm.Model
}
