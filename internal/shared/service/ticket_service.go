package service

import (
	"fmt"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	"github.com/andrianprasetya/eventHub/internal/ticket/dto/request"
	modelTicket "github.com/andrianprasetya/eventHub/internal/ticket/model"
)

func MapEventTicket(eventID, unlimitedTicket string, tickets []request.EventTicket, limitTicket int) ([]*modelTicket.EventTicket, error) {
	eventTickets := make([]*modelTicket.EventTicket, 0, len(tickets))
	var allTicketQuantity int
	for _, ticket := range tickets {
		eventTickets = append(eventTickets, &modelTicket.EventTicket{
			ID:         utils.GenerateID(),
			EventID:    eventID,
			TicketType: ticket.Type,
			Price:      ticket.Price,
			Quantity:   ticket.Quantity,
		})
		allTicketQuantity = allTicketQuantity + ticket.Quantity
	}
	if unlimitedTicket == "false" {
		if allTicketQuantity > limitTicket {
			return nil, validation.ValidationError{
				"event_limit": fmt.Sprintf("number of tickets: %d, your subscription package limit: %d", allTicketQuantity, limitTicket),
			}
		}
	}
	return eventTickets, nil
}
