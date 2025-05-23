package mocks

import (
	"context"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	"gorm.io/gorm"
)

type MockTenantSettingRepository struct {
	CreateBulkWithTxFunc func(ctx context.Context, tx *gorm.DB, tenantSettings []*modelTenant.TenantSetting) error
	GetByTenantIDFunc    func(ctx context.Context, tenantID, key string) (*modelTenant.TenantSetting, error)
}

func (m *MockTenantSettingRepository) CreateBulkWithTx(ctx context.Context, tx *gorm.DB, tenantSettings []*modelTenant.TenantSetting) error {
	if m.CreateBulkWithTxFunc != nil {
		return m.CreateBulkWithTxFunc(ctx, tx, tenantSettings)
	}
	return nil
}

func (m *MockTenantSettingRepository) GetByTenantID(ctx context.Context, tenantID, key string) (*modelTenant.TenantSetting, error) {
	if m.GetByTenantIDFunc != nil {
		return m.GetByTenantIDFunc(ctx, tenantID, key)
	}
	return nil, nil
}
