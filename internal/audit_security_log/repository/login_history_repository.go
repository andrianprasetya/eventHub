package repository

import (
	"github.com/andrianprasetya/eventHub/internal/audit_security_log/model"
	"gorm.io/gorm"
)

type LoginHistoryRepository interface {
	Create(history *model.LoginHistory) error
}

type loginHistoryRepository struct {
	DB *gorm.DB
}

func NewLoginHistoryRepository(db *gorm.DB) LoginHistoryRepository {
	return &loginHistoryRepository{DB: db}
}

func (r loginHistoryRepository) Create(history *model.LoginHistory) error {
	return r.DB.Create(history).Error
}
