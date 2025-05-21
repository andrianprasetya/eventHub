package repository

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/tenant/model"
	"gorm.io/gorm"
)

type TenantRepository interface {
	CreateWithTx(ctx context.Context, tx *gorm.DB, tenant *model.Tenant) error
	GetByID(ctx context.Context, id string) (*model.Tenant, error)
	Update(ctx context.Context, tenant *model.Tenant) error
}

type tenantRepository struct {
	DB *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{DB: db}
}

func (r tenantRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, tenant *model.Tenant) error {
	return tx.WithContext(ctx).Create(tenant).Error
}

func (r tenantRepository) GetByID(ctx context.Context, id string) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.DB.WithContext(ctx).First(&tenant, "id = ?", id).Error
	return &tenant, err
}

func (r tenantRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	return r.DB.WithContext(ctx).Save(tenant).Error
}
