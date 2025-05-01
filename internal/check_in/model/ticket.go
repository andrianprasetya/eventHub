package model

import "time"

type Ticket struct {
	ID        string `gorm:"type:varchar(50);primary_key:true"`
	OrderID   string `gorm:"type:varchar(50);not null;index"`
	QRCode    string `gorm:"type:varchar(255);not null"`
	Status    string `gorm:"type:varchar(15);not null;comment: valid | used | cancelled"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
