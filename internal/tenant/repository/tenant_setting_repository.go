package repository

import (
	"context"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	"gorm.io/gorm"
)

type TenantSettingRepository interface {
	CreateBulkWithTx(ctx context.Context, tx *gorm.DB, tenantSettings []*modelTenant.TenantSetting) error
	GetByTenantID(ctx context.Context, tenantID, key string) (*modelTenant.TenantSetting, error)
}

type tenantSettingRepository struct {
	DB *gorm.DB
}

func NewTenantSettingRepository(db *gorm.DB) TenantSettingRepository {
	return &tenantSettingRepository{DB: db}
}

func (r *tenantSettingRepository) CreateBulkWithTx(ctx context.Context, tx *gorm.DB, tenantSettings []*modelTenant.TenantSetting) error {
	return tx.WithContext(ctx).Create(tenantSettings).Error
}

func (r *tenantSettingRepository) GetByTenantID(ctx context.Context, tenantID, key string) (*modelTenant.TenantSetting, error) {
	var tenantSetting modelTenant.TenantSetting
	err := r.DB.WithContext(ctx).Where("tenant_id = ?", tenantID).Where("key = ?", key).First(&tenantSetting).Error
	return &tenantSetting, err
}
