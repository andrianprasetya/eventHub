package repository

import (
	"github.com/andrianprasetya/eventHub/internal/tenant/model"
	"gorm.io/gorm"
)

type TenantRepository interface {
	CreateWithTx(tx *gorm.DB, tenant *model.Tenant) error
	GetByID(id string) (*model.Tenant, error)
	Update(tenant *model.Tenant) error
}

type tenantRepository struct {
	DB *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{DB: db}
}

func (r tenantRepository) CreateWithTx(tx *gorm.DB, tenant *model.Tenant) error {
	return tx.Create(tenant).Error
}

func (r tenantRepository) GetByID(id string) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.DB.First(&tenant, "id = ?", id).Error
	return &tenant, err
}

func (r tenantRepository) Update(tenant *model.Tenant) error {
	return r.DB.Save(tenant).Error
}
