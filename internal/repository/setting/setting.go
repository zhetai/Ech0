package repository

import (
	"context"

	"gorm.io/gorm"

	model "github.com/lin-snow/ech0/internal/model/setting"
	"github.com/lin-snow/ech0/internal/transaction"
)

type SettingRepository struct {
	db func() *gorm.DB
}

func NewSettingRepository(dbProvider func() *gorm.DB) SettingRepositoryInterface {
	return &SettingRepository{
		db: dbProvider,
	}
}

func (settingRepository *SettingRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transaction.TxKey).(*gorm.DB); ok {
		return tx
	}
	return settingRepository.db()
}

// CreateWebhook 创建一个webhook
func (settingRepository *SettingRepository) CreateWebhook(ctx context.Context, webhook *model.Webhook) error {
	if err := settingRepository.getDB(ctx).Create(webhook).Error; err != nil {
		return err
	}

	return nil
}

// GetAllWebhooks 获取所有webhooks
func (settingRepository *SettingRepository) GetAllWebhooks() ([]model.Webhook, error) {
	var webhooks []model.Webhook
	if err := settingRepository.db().Find(&webhooks).Error; err != nil {
		return nil, err
	}

	return webhooks, nil
}

// DeleteWebhookByID 根据ID删除webhook
func (settingRepository *SettingRepository) DeleteWebhookByID(ctx context.Context, id uint) error {
	if err := settingRepository.getDB(ctx).Delete(&model.Webhook{}, id).Error; err != nil {
		return err
	}

	return nil
}
