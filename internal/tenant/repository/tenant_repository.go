package repository

import (
	"github.com/andrianprasetya/eventHub/internal/tenant/model"
	"gorm.io/gorm"
)

type TenantRepository interface {
	Create(tenant *model.Tenant) error
}

type tenantRepository struct {
	DB *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{DB: db}
}

func (r tenantRepository) Create(tenant *model.Tenant) error {
	return r.DB.Create(tenant).Error
}
