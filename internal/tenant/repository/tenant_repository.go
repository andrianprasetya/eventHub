package repository

import (
	"github.com/andrianprasetya/eventHub/internal/tenant"
	"github.com/andrianprasetya/eventHub/internal/tenant/model"
	"gorm.io/gorm"
)

type tenantRepository struct {
	DB *gorm.DB
}

func NewTenantRepository(db *gorm.DB) tenant.TenantRepository {
	return &tenantRepository{DB: db}
}

func (r tenantRepository) Create(tenant *model.Tenant) error {
	return r.DB.Create(tenant).Error
}
