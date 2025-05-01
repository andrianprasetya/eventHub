package repository

import (
	"github.com/andrianprasetya/eventHub/internal/user"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{DB: db}
}

func (r userRepository) Create(user *modelUser.User) error {
	return r.DB.Create(user).Error
}

func (r userRepository) GetByEmail(email string) (*modelUser.User, error) {
	var user modelUser.User
	if err := r.DB.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
