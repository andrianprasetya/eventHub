package repository

import (
	"github.com/andrianprasetya/eventHub/internal/order/model"
	"gorm.io/gorm"
)

type InvoiceRepository interface {
	Create(invoice *model.Invoice) error
}

type invoiceRepository struct {
	DB *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{DB: db}
}

func (r *invoiceRepository) Create(invoice *model.Invoice) error {
	return nil
}
