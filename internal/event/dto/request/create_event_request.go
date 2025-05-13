package request

import (
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"time"
)

type CreateEventRequest struct {
	Title       string             `json:"title" validate:"required"`
	Description *string            `json:"description"`
	Location    string             `json:"location" validate:"required"`
	StartDate   time.Time          `json:"start_date" validate:"required"`
	EndDate     time.Time          `json:"end_date" validate:"required"`
	CategoryID  string             `json:"category_id" validate:"required"`
	Tags        *utils.StringArray `json:"tags" validate:"omitempty,dive,required"`
	Status      string             `json:"status"`
}
