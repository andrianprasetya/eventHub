package model

import (
	modelEvent "github.com/andrianprasetya/eventHub/internal/event/model"
	"time"
)

type Order struct {
	ID            string           `gorm:"type:varchar(50);primary_key:true"`
	EventID       string           `gorm:"type:varchar(50);not null;index"`
	Name          string           `gorm:"type:varchar(100);not null;index"`
	Email         string           `gorm:"type:varchar(50);not null;index"`
	TotalAmount   int              `gorm:"type:integer;not null"`
	Status        string           `gorm:"type:varchar(15);not null;comment: pending | paid | cancelled"`
	PaymentMethod string           `gorm:"type:varchar(50);not null"`
	Event         modelEvent.Event `gorm:"foreignKey:EventID;references:ID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
