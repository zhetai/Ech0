package repository

import (
	"gorm.io/gorm"
)

type SettingRepository struct {
	db func() *gorm.DB
}

func NewSettingRepository(dbProvider func() *gorm.DB) SettingRepositoryInterface {
	return &SettingRepository{
		db: dbProvider,
	}
}

// func (settingRepository *SettingRepository) getDB(ctx context.Context) *gorm.DB {
// 	if tx, ok := ctx.Value(transaction.TxKey).(*gorm.DB); ok {
// 		return tx
// 	}
// 	return settingRepository.db()
// }
