package usecases

import (
	"github.com/andrianprasetya/eventHub/internal/models"
	"github.com/andrianprasetya/eventHub/internal/repositories"
	"github.com/andrianprasetya/eventHub/internal/utils"
	"time"
)

type EventUsecase struct {
	EventRepo *repositories.EventRepository
}

func NewEventUsecase(eventRepo *repositories.EventRepository) *EventUsecase {
	return &EventUsecase{EventRepo: eventRepo}
}

func (u *EventUsecase) CreateEvent(tenantID, title, description, location string, startTime, endTime time.Time, status string) (*models.Event, error) {
	event := &models.Event{
		ID:          utils.GenerateID(),
		TenantID:    tenantID,
		Title:       title,
		Description: description,
		Location:    location,
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      status,
	}

	err := u.EventRepo.Create(event)
	return event, err
}

func (u *EventUsecase) GetEventByID(id string) (*models.Event, error) {
	return u.EventRepo.GetByID(id)
}

func (u *EventUsecase) GetEventsByTenant(tenantID string) ([]models.Event, error) {
	return u.EventRepo.GetByTenantID(tenantID)
}

func (u *EventUsecase) UpdateEvent(event *models.Event) error {
	return u.EventRepo.Update(event)
}

func (u *EventUsecase) DeleteEvent(id string) error {
	return u.EventRepo.Delete(id)
}
