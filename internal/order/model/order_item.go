package model

import "time"

type OrderItem struct {
	ID        string `gorm:"type:varchar(50);primary_key:true"`
	OrderID   string `gorm:"type:varchar(50);not null;index"`
	TicketID  string `gorm:"type:varchar(50);not null;index"`
	Quantity  int    `gorm:"type:integer;default:1"`
	Price     int    `gorm:"type:integer;default:0"`
	Order     Order  `gorm:"foreignKey:OrderID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
