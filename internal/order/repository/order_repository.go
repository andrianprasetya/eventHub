package repository

import (
	"github.com/andrianprasetya/eventHub/internal/order/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order) error
}

type orderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{DB: db}
}

func (r *orderRepository) Create(order *model.Order) error {
	return nil
}
