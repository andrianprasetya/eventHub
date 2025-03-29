package models

import (
	"gorm.io/gorm"
)

// Ticket model
type Ticket struct {
	ID      string `gorm:"type:varchar(50);primary_key:true" json:"id"`
	OrderID string `gorm:"type:varchar(100);not null" json:"order_id"`
	QrCode  string `gorm:"type:varchar(255);not null" json:"qr_code"`
	Status  string `gorm:"type:varchar(25)" json:"status"`
	gorm.Model
}
