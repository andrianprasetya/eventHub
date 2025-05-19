package repository

import (
	"github.com/andrianprasetya/eventHub/internal/order/model"
	"gorm.io/gorm"
)

type OrderItemRepository interface {
	Create(orderItem *model.OrderItem) error
}

type orderItemRepository struct {
	DB *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) OrderItemRepository {
	return &orderItemRepository{DB: db}
}

func (r *orderItemRepository) Create(orderItem *model.OrderItem) error {
	return nil
}
