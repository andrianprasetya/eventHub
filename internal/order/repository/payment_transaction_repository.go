package repository

import (
	"github.com/andrianprasetya/eventHub/internal/order/model"
	"gorm.io/gorm"
)

type PaymentTransactionRepository interface {
	Create(paymentTransaction *model.PaymentTransaction) error
}

type paymentTransactionRepository struct {
	DB *gorm.DB
}

func NewPaymentTransactionRepository(db *gorm.DB) PaymentTransactionRepository {
	return &paymentTransactionRepository{DB: db}
}

func (r *paymentTransactionRepository) Create(paymentTransaction *model.PaymentTransaction) error {
	return nil
}
