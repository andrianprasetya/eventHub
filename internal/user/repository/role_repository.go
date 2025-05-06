package repository

import (
	"github.com/andrianprasetya/eventHub/internal/user/model"
	"gorm.io/gorm"
)

type RoleRepository interface {
	GetRole(slug string) (*model.Role, error)
}

type roleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{DB: db}
}

func (r roleRepository) GetRole(slug string) (*model.Role, error) {
	var role model.Role
	if err := r.DB.First(&role, "slug = ?", slug).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
