package request

import (
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/ticket/dto/request"
	"time"
)

type CreateEventRequest struct {
	Title       string                  `json:"title" validate:"required"`
	Description *string                 `json:"description"`
	EventType   string                  `json:"event_type" validate:"required"`
	Location    *string                 `json:"location" validate:"required"`
	StartDate   time.Time               `json:"start_date" validate:"required,not_past_date"`
	EndDate     time.Time               `json:"end_date" validate:"required"`
	CategoryID  string                  `json:"category_id" validate:"required"`
	Tags        *utils.StringArray      `json:"tags" validate:"is_array,required"`
	IsTicket    int                     `json:"is_ticket" validate:"required,smallint"`
	Status      string                  `json:"status" validate:"required"`
	Tickets     []request.EventTicket   `json:"tickets" validate:"omitempty,dive"`
	Sessions    []EventSession          `json:"sessions" validate:"omitempty,dive"`
	Discounts   []request.EventDiscount `json:"discounts" validate:"omitempty,dive"`
}

type EventSession struct {
	Title         string    `json:"title" validate:"required"`
	StartDateTime time.Time `json:"start_date_time" validate:"required,not_past_datetime"`
	EndDateTime   time.Time `json:"end_date_time" validate:"required"`
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

type UpdateEventRequest struct {
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
}
