package mocks

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/audit_security_log/model"
)

type MockLogActivityRepository struct {
	CreateFunc func(ctx context.Context, history *model.ActivityLog) error
}

func (m *MockLogActivityRepository) Create(ctx context.Context, history *model.ActivityLog) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, history)
	}
	return nil
}
