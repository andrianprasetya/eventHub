package repository

import (
	"github.com/andrianprasetya/eventHub/internal/ticket/model"
	"gorm.io/gorm"
)

type TicketRepository interface {
	CreateBulkWithTx(tx *gorm.DB, ticket []*model.EventTicket) error
}

type ticketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{DB: db}
}

func (r *ticketRepository) CreateBulkWithTx(tx *gorm.DB, ticket []*model.EventTicket) error {
	return tx.Create(ticket).Error
}
