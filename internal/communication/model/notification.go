package model

import "time"

type Notification struct {
	ID        string `gorm:"type:varchar(50);primary_key:true"`
	UserID    string `gorm:"type:varchar(50);not null;index"`
	Message   string `gorm:"type:varchar(255);not null"`
	IsRead    string `gorm:"type:smallint;default:0;comment: 0 => false | 1 => true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
