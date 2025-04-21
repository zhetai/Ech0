package services

import (
	"github.com/lin-snow/ech0/config"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/repository"
)

// GetSettings 获取系统设置
func GetSetting() (models.SystemSetting, error) {
	var setting models.SystemSetting
	setting, err := repository.GetKeyValue[models.SystemSetting](models.SystemSettingsKey)
	if err != nil {
		// 未获取到设置，将默认值加入到数据库中
		setting.ServerName = config.Config.Setting.Servername
		setting.AllowRegister = config.Config.Setting.AllowRegister
		error := repository.AddKeyValue(models.SystemSettingsKey, setting)
		if error != nil {
			return setting, error
		}
	}

	return setting, nil
}

// 更新系统设置
func UpdateSetting(newSetting models.SystemSetting) error {
	// 更新数据库中的设置
	err := repository.UpdateKeyValue(models.SystemSettingsKey, newSetting)
	if err != nil {
		return err
	}
	return nil
}
