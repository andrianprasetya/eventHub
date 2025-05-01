package repository

import (
	"github.com/andrianprasetya/eventHub/internal/user"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
	"gorm.io/gorm"
)

type roleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) user.RoleRepository {
	return &roleRepository{DB: db}
}

func (r roleRepository) CreateUserRole(userRole *modelUser.UserRole) error {
	return r.DB.Create(&userRole).Error
}

func (r roleRepository) GetRole(name string) (*modelUser.Role, error) {
	var role modelUser.Role
	if err := r.DB.First(&role, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
