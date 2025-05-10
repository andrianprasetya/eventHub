package repository

import (
	"github.com/andrianprasetya/eventHub/internal/user/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(tx *gorm.DB, user *model.User) error
	GetByEmail(email string) (*model.User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r userRepository) Create(tx *gorm.DB, user *model.User) error {
	return tx.Create(user).Error
}

func (r userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.DB.Preload("Tenant").Preload("Role").First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
