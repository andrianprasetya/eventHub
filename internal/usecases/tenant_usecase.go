package usecases

import (
	"github.com/andrianprasetya/eventHub/internal/models"
	"github.com/andrianprasetya/eventHub/internal/repositories"
	"github.com/andrianprasetya/eventHub/internal/utils"
)

type TenantUsecase struct {
	TenantRepo *repositories.TenantRepository
}

func NewTenantUsecase(tenantRepo *repositories.TenantRepository) *TenantUsecase {
	return &TenantUsecase{TenantRepo: tenantRepo}
}

func (u *TenantUsecase) CreateTenant(name, email, logo, plan string) (*models.Tenant, error) {
	tenant := &models.Tenant{
		ID:               utils.GenerateID(),
		Name:             name,
		Email:            email,
		Logo:             logo,
		SubscriptionPlan: plan,
	}

	err := u.TenantRepo.Create(tenant)
	return tenant, err
}

func (u *TenantUsecase) GetTenantByID(id string) (*models.Tenant, error) {
	return u.TenantRepo.GetByID(id)
}

func (u *TenantUsecase) UpdateTenant(tenant *models.Tenant) error {
	return u.TenantRepo.Update(tenant)
}

func (u *TenantUsecase) DeleteTenant(id string) error {
	return u.TenantRepo.Delete(id)
}
