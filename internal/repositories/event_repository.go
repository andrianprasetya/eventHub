package repositories

import (
	"github.com/andrianprasetya/eventHub/internal/models"
	"gorm.io/gorm"
)

type EventRepository struct {
	DB *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{DB: db}
}

func (r *EventRepository) Create(event *models.Event) error {
	return r.DB.Create(event).Error
}

func (r *EventRepository) GetByID(id string) (*models.Event, error) {
	var event models.Event
	err := r.DB.First(&event, "id = ?", id).Error
	return &event, err
}

func (r *EventRepository) GetByTenantID(tenantID string) ([]models.Event, error) {
	var events []models.Event
	err := r.DB.Where("tenant_id = ?", tenantID).Find(&events).Error
	return events, err
}

func (r *EventRepository) Update(event *models.Event) error {
	return r.DB.Save(event).Error
}

func (r *EventRepository) Delete(id string) error {
	return r.DB.Delete(&models.Event{}, "id = ?", id).Error
}
