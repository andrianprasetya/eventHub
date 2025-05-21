package repository

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/ticket/model"
	"gorm.io/gorm"
)

type TicketRepository interface {
	CreateBulkWithTx(ctx context.Context, tx *gorm.DB, ticket []*model.EventTicket) error
}

type ticketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{DB: db}
}

func (r *ticketRepository) CreateBulkWithTx(ctx context.Context, tx *gorm.DB, ticket []*model.EventTicket) error {
	return tx.WithContext(ctx).Create(ticket).Error
}
