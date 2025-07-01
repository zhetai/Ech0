package service

import model "github.com/lin-snow/ech0/internal/model/setting"

type SettingServiceInterface interface {
	// GetSetting 获取设置
	GetSetting(setting *model.SystemSetting) error

	// UpdateSetting 更新设置
	UpdateSetting(userid uint, newSetting *model.SystemSettingDto) error
}
