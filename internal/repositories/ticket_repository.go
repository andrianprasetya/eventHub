package repositories

import (
	"github.com/andrianprasetya/eventHub/internal/models"
	"gorm.io/gorm"
)

type TicketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{DB: db}
}

func (r *TicketRepository) Create(ticket *models.Ticket) error {
	return r.DB.Create(ticket).Error
}

func (r *TicketRepository) GetByID(id string) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.DB.First(&ticket, "id = ?", id).Error
	return &ticket, err
}

func (r *TicketRepository) Update(ticket *models.Ticket) error {
	return r.DB.Save(ticket).Error
}

func (r *TicketRepository) Delete(id string) error {
	return r.DB.Delete(&models.Ticket{}, "id = ?", id).Error
}
