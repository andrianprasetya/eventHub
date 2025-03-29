package repositories

import (
	"github.com/andrianprasetya/eventHub/internal/models"
	"gorm.io/gorm"
)

type TenantRepository struct {
	DB *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{DB: db}
}

func (r *TenantRepository) Create(tenant *models.Tenant) error {
	return r.DB.Create(tenant).Error
}

func (r *TenantRepository) GetByID(id string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.DB.First(&tenant, "id = ?", id).Error
	return &tenant, err
}

func (r *TenantRepository) Update(tenant *models.Tenant) error {
	return r.DB.Save(tenant).Error
}

func (r *TenantRepository) Delete(id string) error {
	return r.DB.Delete(&models.Tenant{}, "id = ?", id).Error
}
