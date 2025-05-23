package repository

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/ticket/dto/request"
	"github.com/andrianprasetya/eventHub/internal/ticket/model"
	"gorm.io/gorm"
)

type TicketRepository interface {
	CreateBulkWithTx(ctx context.Context, tx *gorm.DB, ticket []*model.EventTicket) error
	GetAll(ctx context.Context, query request.TicketPaginateParams, tenantID *string) ([]*model.EventTicket, int64, error)
	GetByID(ctx context.Context, id string) (*model.EventTicket, error)
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

func (r *ticketRepository) GetAll(ctx context.Context, query request.TicketPaginateParams, tenantID *string) ([]*model.EventTicket, int64, error) {
	var tickets []*model.EventTicket
	var total int64

	db := r.DB.WithContext(ctx).Preload("Event").Model(&model.EventTicket{})
	if tenantID != nil {
		db = db.Joins("JOIN events ON events.id = discounts.event_id").Where("events.tenant_id = ?", *tenantID)
	}

	if query.Name != nil {
		db = db.Where("name ILIKE ?", "%"+*query.Name+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize

	if err := db.Limit(query.PageSize).Offset(offset).Find(&tickets).Error; err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}

func (r *ticketRepository) GetByID(ctx context.Context, id string) (*model.EventTicket, error) {
	var ticket model.EventTicket
	if err := r.DB.WithContext(ctx).First(&ticket, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}
