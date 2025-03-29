package usecases

import (
	"github.com/andrianprasetya/eventHub/internal/models"
	"github.com/andrianprasetya/eventHub/internal/repositories"
	"github.com/andrianprasetya/eventHub/internal/utils"
)

type TicketUsecase struct {
	TicketRepo *repositories.TicketRepository
}

func NewTicketUsecase(ticketRepo *repositories.TicketRepository) *TicketUsecase {
	return &TicketUsecase{
		TicketRepo: ticketRepo,
	}
}

// CreateTicket generates a new ticket for an order
func (uc *TicketUsecase) CreateTicket(orderID, qrCode, status string) (*models.Ticket, error) {
	ticket := &models.Ticket{
		ID:      utils.GenerateID(),
		OrderID: orderID,
		QrCode:  qrCode,
		Status:  status,
	}

	err := uc.TicketRepo.Create(ticket)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

// GetTicket retrieves a ticket by its ID
func (uc *TicketUsecase) GetTicket(ticketID string) (*models.Ticket, error) {
	return uc.TicketRepo.GetByID(ticketID)
}

// UpdateTicket updates ticket status
func (uc *TicketUsecase) UpdateTicket(ticketID, status string) (*models.Ticket, error) {
	ticket, err := uc.TicketRepo.GetByID(ticketID)
	if err != nil {
		return nil, err
	}

	ticket.Status = status
	err = uc.TicketRepo.Update(ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

// DeleteTicket removes a ticket
func (uc *TicketUsecase) DeleteTicket(ticketID string) error {
	return uc.TicketRepo.Delete(ticketID)
}
