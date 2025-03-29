package models

import (
	"gorm.io/gorm"
	"time"
)

// Order model
type Order struct {
	ID            string    `gorm:"type:varchar(50);primary_key:true" json:"id"`
	EventID       string    `gorm:"type:varchar(100);not null" json:"event_id"`
	CustomerName  string    `gorm:"type:varchar(100);not null" json:"customer_name"`
	CustomerEmail string    `gorm:"type:varchar(100);not null" json:"customer_email"`
	TotalAmount   float64   `gorm:"type:decimal(10,2)" json:"total_amount"`
	PaymentMethod string    `gorm:"type:varchar(100)" json:"payment_method"`
	Status        string    `gorm:"type:varchar(25)" json:"status"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	gorm.Model
}
