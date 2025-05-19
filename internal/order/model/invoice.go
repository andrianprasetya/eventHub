package model

import "time"

type Invoice struct {
	ID            string    `gorm:"type:varchar(50);primary_key:true"`
	OrderID       string    `gorm:"type:varchar(50);not null"`
	InvoiceNumber string    `gorm:"type:varchar(100);not null"`
	IssueDate     time.Time `gorm:"type:date;not null"`
	DueDate       time.Time `gorm:"type:date;not null"`
	Total         int       `gorm:"type:integer;default:0"`
	Order         Order     `gorm:"foreignKey:OrderID;references:ID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
