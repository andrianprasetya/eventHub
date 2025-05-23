package request

type CreateTicketRequest struct {
	EventID string        `json:"event_id" validate:"required"`
	Tickets []EventTicket `json:"tickets" validate:"required,dive"`
}

type EventTicket struct {
	Type     string `json:"type"`
	Price    int    `json:"price" validate:"required,numeric"`
	Quantity int    `json:"quantity" validate:"required,numeric,min=1"`
}

type UpdateTicketRequest struct {
	Type     string `json:"type"`
	Price    int    `json:"price" validate:"required,numeric"`
	Quantity int    `json:"quantity" validate:"required,numeric,min=1"`
}

type TicketPaginateParams struct {
	Page     int     `query:"page"`
	PageSize int     `query:"pageSize"`
	Name     *string `query:"name"`
}
