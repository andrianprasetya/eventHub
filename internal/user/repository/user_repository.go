package repository

import (
	"github.com/andrianprasetya/eventHub/internal/user/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	CreateWithTx(tx *gorm.DB, user *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetByID(id string) (*model.User, error)
	GetAll(page, pageSize int, tenantID *string) ([]*model.User, int64, error)
	Update(user *model.User) error
	Delete(id string) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *userRepository) CreateWithTx(tx *gorm.DB, user *model.User) error {
	return tx.Create(user).Error
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.DB.Preload("Tenant").Preload("Role").First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAll(page, pageSize int, tenantID *string) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	db := r.DB.Model(&model.User{})
	if tenantID != nil {
		db = db.Where("tenant_id = ?", tenantID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := db.Preload("Role").Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	if err := r.DB.Preload("Tenant").Preload("Role").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.DB.Save(user).Error
}

func (r *userRepository) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&model.User{}).Error
}
