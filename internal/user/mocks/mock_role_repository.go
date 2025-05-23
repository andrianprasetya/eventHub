package mocks

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/model"
)

type MockRoleRepository struct {
	GetAllFunc    func(qctx context.Context, uery request.RolePaginateParams, role string) ([]*model.Role, int64, error)
	GetByIDFunc   func(ctx context.Context, id string) (*model.Role, error)
	GetBySlugFunc func(ctx context.Context, slug string) (*model.Role, error)
}

func (m *MockRoleRepository) GetAll(ctx context.Context, query request.RolePaginateParams, role string) ([]*model.Role, int64, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, query, role)
	}
	return nil, 0, nil
}

func (m *MockRoleRepository) GetByID(ctx context.Context, id string) (*model.Role, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockRoleRepository) GetBySlug(ctx context.Context, slug string) (*model.Role, error) {
	if m.GetBySlugFunc != nil {
		return m.GetBySlugFunc(ctx, slug)
	}
	return nil, nil
}
