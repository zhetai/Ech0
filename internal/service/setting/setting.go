package service

import (
	"errors"

	"github.com/lin-snow/ech0/internal/config"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/setting"
	keyvalueRepository "github.com/lin-snow/ech0/internal/repository/keyvalue"
	commonService "github.com/lin-snow/ech0/internal/service/common"
	httpUtil "github.com/lin-snow/ech0/internal/util/http"
	jsonUtil "github.com/lin-snow/ech0/internal/util/json"
)

type SettingService struct {
	commonService      commonService.CommonServiceInterface
	keyvalueRepository keyvalueRepository.KeyValueRepositoryInterface
}

func NewSettingService(commonService commonService.CommonServiceInterface, keyvalueRepository keyvalueRepository.KeyValueRepositoryInterface) SettingServiceInterface {
	return &SettingService{
		commonService:      commonService,
		keyvalueRepository: keyvalueRepository,
	}
}

// GetSetting 获取设置
func (settingService *SettingService) GetSetting(setting *model.SystemSetting) error {
	settingValue, err := settingService.keyvalueRepository.GetKeyValue(commonModel.SystemSettingsKey)
	if err != nil {
		// 数据库中不存在数据，手动添加初始数据
		setting.SiteTitle = config.Config.Setting.SiteTitle
		setting.ServerName = config.Config.Setting.Servername
		setting.ServerURL = config.Config.Setting.Serverurl
		setting.AllowRegister = config.Config.Setting.AllowRegister
		setting.ICPNumber = config.Config.Setting.Icpnumber
		setting.MetingAPI = config.Config.Setting.MetingAPI
		setting.CommentAPI = config.Config.Setting.CommentAPI
		setting.CustomCSS = config.Config.Setting.CustomCSS
		setting.CustomJS = config.Config.Setting.CustomJS

		// 处理 URL
		setting.ServerURL = httpUtil.TrimURL(setting.ServerURL)
		setting.MetingAPI = httpUtil.TrimURL(setting.MetingAPI)
		setting.CommentAPI = httpUtil.TrimURL(setting.CommentAPI)

		// 序列化为 JSON
		settingToJSON, err := jsonUtil.JSONMarshal(setting)
		if err != nil {
			return err
		}
		if err := settingService.keyvalueRepository.AddKeyValue(commonModel.SystemSettingsKey, string(settingToJSON)); err != nil {
			return err
		}
	}

	if err := jsonUtil.JSONUnmarshal([]byte(settingValue.(string)), setting); err != nil {
		return err
	}

	return nil
}

// UpdateSetting 更新设置
func (settingService *SettingService) UpdateSetting(userid uint, newSetting *model.SystemSettingDto) error {
	user, err := settingService.commonService.CommonGetUserByUserId(userid)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	var setting model.SystemSetting
	setting.SiteTitle = newSetting.SiteTitle
	setting.ServerName = newSetting.ServerName
	setting.ServerURL = httpUtil.TrimURL(newSetting.ServerURL)
	setting.AllowRegister = newSetting.AllowRegister
	setting.ICPNumber = newSetting.ICPNumber
	setting.MetingAPI = httpUtil.TrimURL(newSetting.MetingAPI)
	setting.CommentAPI = httpUtil.TrimURL(newSetting.CommentAPI)
	setting.CustomCSS = newSetting.CustomCSS
	setting.CustomJS = newSetting.CustomJS

	// 序列化为 JSON
	settingToJSON, err := jsonUtil.JSONMarshal(setting)
	if err != nil {
		return err
	}

	// 将字节切片转换为字符串
	settingToJSONString := string(settingToJSON)
	if err := settingService.keyvalueRepository.UpdateKeyValue(commonModel.SystemSettingsKey, settingToJSONString); err != nil {
		return err
	}

	return nil
}
