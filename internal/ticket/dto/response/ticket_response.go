package response

type TicketResponse struct {
	ID         string `json:"id"`
	EventID    string `json:"event_id"`
	TicketType string `json:"ticket_type"`
	Price      int    `json:"price"`
	Quantity   int    `json:"quantity"`
	Sold       int    `json:"sold"`
}

type TicketListItemResponse struct {
	ID         string `json:"id"`
	TicketType string `json:"ticket_type"`
	Price      int    `json:"price"`
	Quantity   int    `json:"quantity"`
}
