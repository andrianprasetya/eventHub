package usecase

import (
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

type MockTenantRepo struct {
	mock.Mock
}

func (m *MockTenantRepo) Create(tx *gorm.DB, tenant *modelTenant.Tenant) error {
	args := m.Called(tx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepo) GetByID(id string) (*modelTenant.Tenant, error) {
	args := m.Called(id)
	return args.Get(0).(*modelTenant.Tenant), args.Error(1)
}

func (m *MockTenantRepo) Update(tenant *modelTenant.Tenant) error {
	args := m.Called(tenant)
	return args.Error(0)
}

func TestTenantUsecaseCreate_Success(t *testing.T) {

}
