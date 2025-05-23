package mapper

import (
	"github.com/andrianprasetya/eventHub/internal/ticket/dto/response"
	"github.com/andrianprasetya/eventHub/internal/ticket/model"
)

func FromTicketModel(ticket *model.EventTicket) *response.TicketResponse {
	return &response.TicketResponse{
		EventID:    ticket.EventID,
		TicketType: ticket.TicketType,
		Price:      ticket.Price,
		Quantity:   ticket.Quantity,
		Sold:       ticket.Sold,
	}
}

func FromTicketToListItem(ticket *model.EventTicket) *response.TicketListItemResponse {
	return &response.TicketListItemResponse{
		TicketType: ticket.TicketType,
		Price:      ticket.Price,
		Quantity:   ticket.Quantity,
	}
}

func FromTicketToList(tickets []*model.EventTicket) []*response.TicketListItemResponse {
	result := make([]*response.TicketListItemResponse, 0, len(tickets))
	for _, ticket := range tickets {
		result = append(result, FromTicketToListItem(ticket))
	}
	return result
}
