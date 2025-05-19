package repository

import (
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	"gorm.io/gorm"
)

type TenantSettingRepository interface {
	CreateBulkWithTx(tx *gorm.DB, tenantSettings []*modelTenant.TenantSetting) error
	GetByTenantID(tenantID, key string) (*modelTenant.TenantSetting, error)
}

type tenantSettingRepository struct {
	DB *gorm.DB
}

func NewTenantSettingRepository(db *gorm.DB) TenantSettingRepository {
	return &tenantSettingRepository{DB: db}
}

func (r *tenantSettingRepository) CreateBulkWithTx(tx *gorm.DB, tenantSettings []*modelTenant.TenantSetting) error {
	return tx.Create(tenantSettings).Error
}

func (r *tenantSettingRepository) GetByTenantID(tenantID, key string) (*modelTenant.TenantSetting, error) {
	var tenantSetting modelTenant.TenantSetting
	err := r.DB.Find(&tenantSetting).Where("tenant_id = ?", tenantSetting).Where("key = ?", key).Error
	return &tenantSetting, err
}
