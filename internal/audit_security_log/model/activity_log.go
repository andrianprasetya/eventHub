package model

import (
	"github.com/andrianprasetya/eventHub/internal/user/model"
	"time"
)

type ActivityLog struct {
	ID         string     `gorm:"type:varchar(50);primary_key:true"`
	UserID     string     `gorm:"type:varchar(50);not null;index"`
	URL        string     `gorm:"type:varchar(100);not null"`
	Action     string     `gorm:"type:varchar(50);not null"`
	ObjectData string     `gorm:"type:text;not null"`
	ObjectType string     `gorm:"type:varchar(50);not null"`
	ObjectID   string     `gorm:"type:varchar(50);not null;index"`
	User       model.User `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
