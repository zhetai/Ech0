package repository

import model "github.com/lin-snow/ech0/internal/model/setting"

type SettingRepositoryInterface interface {
	// CreateWebhook 创建一个webhook
	CreateWebhook(webhook *model.Webhook) error
}
