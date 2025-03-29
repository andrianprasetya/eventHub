package handlers

import (
	"github.com/andrianprasetya/eventHub/internal/usecases"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type EventHandler struct {
	EventUC *usecases.EventUsecase
}

func NewEventHandler(eventUC *usecases.EventUsecase) *EventHandler {
	return &EventHandler{EventUC: eventUC}
}

func (h *EventHandler) CreateEvent(c echo.Context) error {
	req := struct {
		TenantID    string    `json:"tenant_id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Location    string    `json:"location"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
		Status      string    `json:"status"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	event, err := h.EventUC.CreateEvent(req.TenantID, req.Title, req.Description, req.Location, req.StartTime, req.EndTime, req.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create event"})
	}

	return c.JSON(http.StatusCreated, event)
}

func (h *EventHandler) GetEventByID(c echo.Context) error {
	id := c.Param("id")
	event, err := h.EventUC.GetEventByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Event not found"})
	}
	return c.JSON(http.StatusOK, event)
}
