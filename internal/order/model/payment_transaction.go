package model

import "time"

type PaymentTransaction struct {
	ID              string `gorm:"type:varchar(50);primary_key:true"`
	OrderID         string `gorm:"type:varchar(50);not null;index"`
	PaymentGateway  string `gorm:"type:varchar(100);not null"`
	Amount          int    `gorm:"type:integer;default:0"`
	Status          string `gorm:"type:varchar(15);not null"`
	TransactionData string `gorm:"type:text;not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
