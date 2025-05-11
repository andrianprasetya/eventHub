package repository

import (
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	"gorm.io/gorm"
)

type TenantSettingRepository interface {
	CreateBulk(tx *gorm.DB, tenantSettings []*modelTenant.TenantSetting) error
}

type tenantSettingRepository struct {
	DB *gorm.DB
}

func NewTenantSettingRepository(db *gorm.DB) TenantSettingRepository {
	return &tenantSettingRepository{DB: db}
}

func (r *tenantSettingRepository) CreateBulk(tx *gorm.DB, tenantSettings []*modelTenant.TenantSetting) error {
	return tx.Create(tenantSettings).Error
}
