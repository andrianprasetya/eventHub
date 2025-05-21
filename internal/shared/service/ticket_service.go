package service

import (
	"github.com/andrianprasetya/eventHub/internal/event/dto/request"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	modelTicket "github.com/andrianprasetya/eventHub/internal/ticket/model"
)

func MapEventTicket(eventID string, tickets []request.EventTicket) []*modelTicket.EventTicket {
	eventTickets := make([]*modelTicket.EventTicket, 0, len(tickets))
	for _, ticket := range tickets {
		eventTickets = append(eventTickets, &modelTicket.EventTicket{
			ID:         utils.GenerateID(),
			EventID:    eventID,
			TicketType: ticket.Type,
			Price:      ticket.Price,
			Quantity:   ticket.Quantity,
		})
	}
	return eventTickets
}
