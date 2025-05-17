package response

import (
	"time"
)

type EventResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Category    Category  `json:"category"`
	Tags        []string  `json:"tags"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EventListItemResponse struct {
	ID           string    `json:"ID"`
	Title        string    `json:"title"`
	CategoryName string    `json:"category_name"`
	Tags         []string  `json:"tags"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Status       string    `json:"status"`
}

type EventTagListItemResponse struct {
	Name string `json:"name"`
}

type EventCategoryListItemResponse struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}
