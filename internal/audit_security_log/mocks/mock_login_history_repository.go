package mocks

import (
	"github.com/andrianprasetya/eventHub/internal/audit_security_log/model"
)

type MockLoginHistoryRepository struct {
	CreateFunc func(history *model.LoginHistory) error
}

func (m *MockLoginHistoryRepository) Create(history *model.LoginHistory) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(history)
	}
	return nil
}
