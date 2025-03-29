package handlers

import (
	"github.com/andrianprasetya/eventHub/internal/models"
	"github.com/andrianprasetya/eventHub/internal/usecases"
	"github.com/labstack/echo/v4"
	"net/http"
)

type TicketHandler struct {
	TicketUC *usecases.TicketUsecase
}

func NewTicketHandler(ticketUC *usecases.TicketUsecase) *TicketHandler {
	return &TicketHandler{TicketUC: ticketUC}
}

// CreateTicket creates a new ticket
func (h *TicketHandler) CreateTicket(c echo.Context) error {
	var ticket models.Ticket
	if err := c.Bind(&ticket); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	newTicket, err := h.TicketUC.CreateTicket(ticket.OrderID, ticket.QrCode, ticket.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create ticket"})
	}

	return c.JSON(http.StatusCreated, newTicket)
}

// GetTicket retrieves a ticket by ID
func (h *TicketHandler) GetTicket(c echo.Context) error {
	ticketID := c.Param("id")
	ticket, err := h.TicketUC.GetTicket(ticketID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Ticket not found"})
	}

	return c.JSON(http.StatusOK, ticket)
}

// UpdateTicket updates a ticket's status
func (h *TicketHandler) UpdateTicket(c echo.Context) error {
	ticketID := c.Param("id")
	var updateData struct {
		Status string `json:"status"`
	}

	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	updatedTicket, err := h.TicketUC.UpdateTicket(ticketID, updateData.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update ticket"})
	}

	return c.JSON(http.StatusOK, updatedTicket)
}

// DeleteTicket deletes a ticket by ID
func (h *TicketHandler) DeleteTicket(c echo.Context) error {
	ticketID := c.Param("id")
	if err := h.TicketUC.DeleteTicket(ticketID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete ticket"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Ticket deleted successfully"})
}
