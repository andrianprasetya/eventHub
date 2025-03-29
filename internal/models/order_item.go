package models

import (
	"gorm.io/gorm"
)

// Order Item model
type OrderItem struct {
	ID       string  `gorm:"type:varchar(50);primary_key:true" json:"id"`
	OrderID  string  `gorm:"type:varchar(100);not null" json:"order_id"`
	TicketID string  `gorm:"type:varchar(100);not null" json:"ticket_id"`
	Price    float64 `gorm:"type:decimal(10,2)" json:"price"`
	Quantity int     `gorm:"type:integer" json:"quantity"`
	gorm.Model
}
