package repository

import "gorm.io/gorm"

type TxManager interface {
	Begin() *gorm.DB
}

type GormTxManager struct {
	DB *gorm.DB
}

func (tm *GormTxManager) Begin() *gorm.DB {
	return tm.DB.Begin()
}
