package request

import "time"

type CreateDiscountRequest struct {
	EventID   string          `json:"event_id" validate:"required"`
	Discounts []EventDiscount `json:"discount" validate:"required,dive"`
}

type EventDiscount struct {
	Code               string    `json:"code" validate:"required"`
	DiscountPercentage int       `json:"discount_percentage" validate:"required"`
	StartDate          time.Time `json:"start_date" validate:"required,not_past_date"`
	EndDate            time.Time `json:"end_date" validate:"required"`
}

type DiscountPaginateParams struct {
	Page     int     `query:"page"`
	PageSize int     `query:"pageSize"`
	Name     *string `query:"name"`
}
