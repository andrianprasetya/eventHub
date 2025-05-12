package request

import "time"

type CreateEventRequest struct {
	Title        string       `json:"title"`
	Description  *string      `json:"description"`
	Location     string       `json:"location"`
	StartDate    time.Time    `json:"start_date"`
	EndDate      time.Time    `json:"end_date"`
	CategoryID   string       `json:"category_id"`
	Tags         []string     `json:"tags" validate:"omitempty,dive,required"`
	CustomFields CustomFields `json:"custom_fields" validate:"required,dive"`
}

type CustomFields struct {
	SpeakerCount int    `json:"speaker_count" validate:"required,min=1"`
	Theme        string `json:"theme" validate:"required"`
}
