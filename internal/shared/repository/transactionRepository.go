package repository

import (
	"context"
	"gorm.io/gorm"
)

type TxManager interface {
	Begin(ctx context.Context) *gorm.DB
}

type GormTxManager struct {
	DB *gorm.DB
}

func (tm *GormTxManager) Begin(ctx context.Context) *gorm.DB {
	return tm.DB.WithContext(ctx).Begin()
}
