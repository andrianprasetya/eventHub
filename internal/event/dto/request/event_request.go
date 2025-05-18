package request

import (
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"time"
)

type CreateEventRequest struct {
	Title       string             `json:"title" validate:"required"`
	Description *string            `json:"description"`
	Location    string             `json:"location" validate:"required"`
	StartDate   time.Time          `json:"start_date" validate:"required,not_past_date"`
	EndDate     time.Time          `json:"end_date" validate:"required"`
	CategoryID  string             `json:"category_id" validate:"required"`
	Tags        *utils.StringArray `json:"tags" validate:"is_array,required"`
	Status      string             `json:"status" validate:"required"`
	Tickets     []EventTicket      `json:"tickets" validate:"omitempty,dive"`
	Sessions    []EventSession     `json:"sessions" validate:"omitempty,dive"`
	Discounts   []EventDiscount    `json:"discounts" validate:"omitempty,dive"`
}

type EventTicket struct {
	Type     string `json:"type"`
	Price    int    `json:"price"` // atau pakai int jika harga tidak pakai desimal
	Quantity int    `json:"quantity"`
}

type EventSession struct {
	Title         string    `json:"title"`
	StartDateTime time.Time `json:"start_date_time"`
	EndDateTime   time.Time `json:"end_date_time"`
}

type EventDiscount struct {
	Code               string    `json:"code"`
	DiscountPercentage int       `json:"discount_percentage"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
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
