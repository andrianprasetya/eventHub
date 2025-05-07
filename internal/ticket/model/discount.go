package model

import (
	modelEvent "github.com/andrianprasetya/eventHub/internal/event/model"
	"time"
)

type Discount struct {
	ID                 string           `gorm:"type:varchar(50);primary_key:true"`
	EventID            string           `gorm:"type:varchar(50);not null;index"`
	Code               string           `gorm:"type:varchar(100);not null;index"`
	DiscountPercentage int              `gorm:"type:integer;not null"`
	StartDate          time.Time        `gorm:"type:date;not null"`
	EndDate            time.Time        `gorm:"type:date;not null"`
	Event              modelEvent.Event `gorm:"foreignKey:EventID;references:ID"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
