package model

import "time"

type TicketCoupon struct {
	ID                 string `gorm:"type:varchar(50);primary_key:true"`
	EventTicketID      string `gorm:"type:varchar(50);not null;index"`
	Code               string `gorm:"type:varchar(100);not null;index"`
	DiscountPercentage int    `gorm:"type:integer;not null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
