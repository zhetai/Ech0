package repository

import (
	"gorm.io/gorm"

	model "github.com/lin-snow/ech0/internal/model/setting"
)

type SettingRepository struct {
	db func() *gorm.DB
}

func NewSettingRepository(dbProvider func() *gorm.DB) SettingRepositoryInterface {
	return &SettingRepository{
		db: dbProvider,
	}
}

// CreateWebhook 创建一个webhook
func (settingRepository *SettingRepository) CreateWebhook(webhook *model.Webhook) error {
	return nil
}
