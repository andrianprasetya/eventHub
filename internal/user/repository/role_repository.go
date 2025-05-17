package repository

import (
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/model"
	"gorm.io/gorm"
)

type RoleRepository interface {
	GetAll(query request.RolePaginateParams) ([]*model.Role, int64, error)
	GetByID(id string) (*model.Role, error)
	GetBySlug(slug string) (*model.Role, error)
}

type roleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{DB: db}
}

func (r *roleRepository) GetAll(query request.RolePaginateParams) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	db := r.DB.Model(&model.Role{})
	if query.Name != nil {
		db = db.Where("name ILIKE ?", "%"+*query.Name+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize

	if err := db.Limit(query.PageSize).Offset(offset).Find(&roles).Error; err != nil {
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
