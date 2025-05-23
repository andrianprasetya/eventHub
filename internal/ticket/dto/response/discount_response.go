package response

import "time"

type DiscountResponse struct {
	ID                 string    `json:"id"`
	EventID            string    `json:"event_id"`
	Code               string    `json:"code"`
	DiscountPercentage int       `json:"discount_percentage"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date" `
}

type DiscountListItemResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date" `
}
