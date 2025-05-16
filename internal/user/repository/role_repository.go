package repository

import (
	"github.com/andrianprasetya/eventHub/internal/user/model"
	"gorm.io/gorm"
)

type RoleRepository interface {
	GetAll(page, pageSize int) ([]*model.Role, int64, error)
	GetByID(id string) (*model.Role, error)
	GetBySlug(slug string) (*model.Role, error)
}

type roleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{DB: db}
}

func (r *roleRepository) GetAll(page, pageSize int) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64
	if err := r.DB.Model(&model.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := r.DB.Limit(pageSize).Offset(offset).Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func (r *roleRepository) GetByID(id string) (*model.Role, error) {
	var role model.Role
	if err := r.DB.First(&role, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetBySlug(slug string) (*model.Role, error) {
	var role model.Role
	if err := r.DB.First(&role, "slug = ?", slug).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
