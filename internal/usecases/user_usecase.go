package usecases

import (
	"github.com/andrianprasetya/eventHub/internal/models"
	"github.com/andrianprasetya/eventHub/internal/repositories"
	"github.com/andrianprasetya/eventHub/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	UserRepo *repositories.UserRepository
}

func NewUserUsecase(userRepo *repositories.UserRepository) *UserUsecase {
	return &UserUsecase{UserRepo: userRepo}
}

func (u *UserUsecase) RegisterUser(tenantID, name, email, password, role string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       utils.GenerateID(),
		TenantID: tenantID,
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		IsActive: true,
	}

	err = u.UserRepo.Create(user)
	return user, err
}

func (u *UserUsecase) GetUserByEmail(email string) (*models.User, error) {
	return u.UserRepo.GetByEmail(email)
}

func (u *UserUsecase) GetUserByID(id string) (*models.User, error) {
	return u.UserRepo.GetByID(id)
}
