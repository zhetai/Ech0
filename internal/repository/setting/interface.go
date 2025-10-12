package repository

import (
	"context"

	model "github.com/lin-snow/ech0/internal/model/setting"
)

type SettingRepositoryInterface interface {
	// CreateWebhook 创建一个webhook
	CreateWebhook(ctx context.Context, webhook *model.Webhook) error
}
