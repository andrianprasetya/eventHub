package mocks

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/model"
	"gorm.io/gorm"
)

// MockUserRepository implements UserRepository interface for testing purpose
type MockUserRepository struct {
	// function fields supaya bisa flexible set return value di test
	CreateFunc           func(ctx context.Context, user *model.User) error
	CreateWithTxFunc     func(ctx context.Context, tx *gorm.DB, user *model.User) error
	GetByEmailFunc       func(ctx context.Context, email string) (*model.User, error)
	GetByIDFunc          func(ctx context.Context, id string) (*model.User, error)
	GetAllFunc           func(ctx context.Context, query request.UserPaginateParams, tenantID *string) ([]*model.User, int64, error)
	UpdateFunc           func(ctx context.Context, user *model.User) error
	DeleteFunc           func(ctx context.Context, id string) error
	CountCreatedUserFunc func(ctx context.Context, tenantID string) (int, error)
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, user *model.User) error {
	if m.CreateWithTxFunc != nil {
		return m.CreateWithTxFunc(ctx, tx, user)
	}
	return nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) GetAll(ctx context.Context, query request.UserPaginateParams, tenantID *string) ([]*model.User, int64, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, query, tenantID)
	}
	return nil, 0, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockUserRepository) CountCreatedUser(ctx context.Context, tenantID string) (int, error) {
	if m.CountCreatedUserFunc != nil {
		return m.CountCreatedUserFunc(ctx, tenantID)
	}
	return 0, nil
}
