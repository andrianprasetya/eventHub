package repositories

import (
	"github.com/andrianprasetya/eventHub/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.DB.Create(order).Error
}

func (r *OrderRepository) GetByID(id string) (*models.Order, error) {
	var order models.Order
	err := r.DB.First(&order, "id = ?", id).Error
	return &order, err
}

func (r *OrderRepository) Update(order *models.Order) error {
	return r.DB.Save(order).Error
}

func (r *OrderRepository) Delete(id string) error {
	return r.DB.Delete(&models.Order{}, "id = ?", id).Error
}
