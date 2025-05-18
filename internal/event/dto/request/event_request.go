package request

import (
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"time"
)

type CreateEventRequest struct {
	Title       string             `json:"title" validate:"required"`
	Description *string            `json:"description"`
	EventType   string             `json:"event_type" validate:"required"`
	Location    *string            `json:"location" validate:"required"`
	StartDate   time.Time          `json:"start_date" validate:"required,not_past_date"`
	EndDate     time.Time          `json:"end_date" validate:"required"`
	CategoryID  string             `json:"category_id" validate:"required"`
	Tags        *utils.StringArray `json:"tags" validate:"is_array,required"`
	IsTicket    int                `json:"is_ticket" validate:"required,smallint"`
	Status      string             `json:"status" validate:"required"`
	Tickets     []EventTicket      `json:"tickets" validate:"omitempty,dive"`
	Sessions    []EventSession     `json:"sessions" validate:"omitempty,dive"`
	Discounts   []EventDiscount    `json:"discounts" validate:"omitempty,dive"`
}

type EventTicket struct {
	Type     string `json:"type"`
	Price    int    `json:"price" validate:"required,numeric"` // atau pakai int jika harga tidak pakai desimal
	Quantity int    `json:"quantity" validate:"required,numeric"`
}

type EventSession struct {
	Title         string    `json:"title" validate:"required"`
	StartDateTime time.Time `json:"start_date_time" validate:"required,not_past_date"`
	EndDateTime   time.Time `json:"end_date_time" validate:"required"`
}

type EventDiscount struct {
	Code               string    `json:"code" validate:"required"`
	DiscountPercentage int       `json:"discount_percentage" validate:"required"`
	StartDate          time.Time `json:"start_date" validate:"required,not_past_date"`
	EndDate            time.Time `json:"end_date" validate:"required"`
}

type EventPaginateRequest struct {
	Page     int     `query:"page"`
	PageSize int     `query:"pageSize"`
	Name     *string `json:"name"`
}

type EventTagPaginateRequest struct {
	Page     int     `query:"page"`
	PageSize int     `query:"pageSize"`
	Name     *string `json:"name"`
}
