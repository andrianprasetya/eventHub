package model

import (
	"github.com/andrianprasetya/eventHub/internal/user/model"
	"time"
)

type LoginHistory struct {
	ID        string     `gorm:"type:varchar(50);primary_key:true"`
	UserID    string     `gorm:"type:varchar(50);not null;index"`
	LoginTime time.Time  `gorm:"type:timestamp;not null"`
	IpAddress string     `gorm:"type:varchar(50);not null"`
	User      model.User `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
