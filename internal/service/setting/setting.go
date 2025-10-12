package service

import (
	"context"
	"errors"

	"github.com/lin-snow/ech0/internal/config"
	authModel "github.com/lin-snow/ech0/internal/model/auth"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/setting"
	keyvalueRepository "github.com/lin-snow/ech0/internal/repository/keyvalue"
	commonService "github.com/lin-snow/ech0/internal/service/common"
	"github.com/lin-snow/ech0/internal/transaction"
	httpUtil "github.com/lin-snow/ech0/internal/util/http"
	jsonUtil "github.com/lin-snow/ech0/internal/util/json"
)

type SettingService struct {
	txManager          transaction.TransactionManager
	commonService      commonService.CommonServiceInterface
	keyvalueRepository keyvalueRepository.KeyValueRepositoryInterface
}

func NewSettingService(
	tm transaction.TransactionManager,
	commonService commonService.CommonServiceInterface,
	keyvalueRepository keyvalueRepository.KeyValueRepositoryInterface,
) SettingServiceInterface {
	return &SettingService{
		txManager:          tm,
		commonService:      commonService,
		keyvalueRepository: keyvalueRepository,
	}
}

// GetSetting 获取设置
func (settingService *SettingService) GetSetting(setting *model.SystemSetting) error {
	return settingService.txManager.Run(func(ctx context.Context) error {
		systemSetting, err := settingService.keyvalueRepository.GetKeyValue(commonModel.SystemSettingsKey)
		if err != nil {
			// 数据库中不存在数据，手动添加初始数据
			setting.SiteTitle = config.Config.Setting.SiteTitle
			setting.ServerName = config.Config.Setting.Servername
			setting.ServerURL = config.Config.Setting.Serverurl
			setting.AllowRegister = config.Config.Setting.AllowRegister
			setting.ICPNumber = config.Config.Setting.Icpnumber
			setting.MetingAPI = config.Config.Setting.MetingAPI
			setting.CustomCSS = config.Config.Setting.CustomCSS
			setting.CustomJS = config.Config.Setting.CustomJS

			// 处理 URL
			setting.ServerURL = httpUtil.TrimURL(setting.ServerURL)
			setting.MetingAPI = httpUtil.TrimURL(setting.MetingAPI)

			// 序列化为 JSON
			settingToJSON, err := jsonUtil.JSONMarshal(setting)
			if err != nil {
				return err
			}
			if err := settingService.keyvalueRepository.AddKeyValue(ctx, commonModel.SystemSettingsKey, string(settingToJSON)); err != nil {
				return err
			}

			// 处理 ServerURL
			if err := settingService.keyvalueRepository.AddKeyValue(ctx, commonModel.ServerURLKey, setting.ServerURL); err != nil {
				return err
			}

			return nil
		}

		if err := jsonUtil.JSONUnmarshal([]byte(systemSetting.(string)), setting); err != nil {
			return err
		}

		return nil
	})
}

// UpdateSetting 更新设置
func (settingService *SettingService) UpdateSetting(userid uint, newSetting *model.SystemSettingDto) error {
	return settingService.txManager.Run(func(ctx context.Context) error {
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
		setting.CustomCSS = newSetting.CustomCSS
		setting.CustomJS = newSetting.CustomJS

		// 序列化为 JSON
		settingToJSON, err := jsonUtil.JSONMarshal(setting)
		if err != nil {
			return err
		}

		// 将字节切片转换为字符串
		settingToJSONString := string(settingToJSON)
		if err := settingService.keyvalueRepository.UpdateKeyValue(ctx, commonModel.SystemSettingsKey, settingToJSONString); err != nil {
			return err
		}

		// 更新 ServerURL
		if err := settingService.keyvalueRepository.UpdateKeyValue(ctx, commonModel.ServerURLKey, setting.ServerURL); err != nil {
			return err
		}

		return nil
	})
}

// GetCommentSetting 获取评论设置
func (settingService *SettingService) GetCommentSetting(setting *model.CommentSetting) error {
	return settingService.txManager.Run(func(ctx context.Context) error {
		commentSetting, err := settingService.keyvalueRepository.GetKeyValue(commonModel.CommentSettingKey)
		if err != nil {
			// 数据库中不存在数据，手动添加初始数据
			setting.EnableComment = config.Config.Comment.EnableComment
			setting.Provider = config.Config.Comment.Provider
			setting.CommentAPI = config.Config.Comment.CommentAPI

			// 处理 URL
			setting.CommentAPI = httpUtil.TrimURL(setting.CommentAPI)

			// 序列化为 JSON
			settingToJSON, err := jsonUtil.JSONMarshal(setting)
			if err != nil {
				return err
			}
			if err := settingService.keyvalueRepository.AddKeyValue(ctx, commonModel.CommentSettingKey, string(settingToJSON)); err != nil {
				return err
			}

			return nil
		}

		if err := jsonUtil.JSONUnmarshal([]byte(commentSetting.(string)), setting); err != nil {
			return err
		}

		return nil
	})
}

// UpdateCommentSetting 更新评论设置
func (settingService *SettingService) UpdateCommentSetting(userid uint, newSetting *model.CommentSettingDto) error {
	return settingService.txManager.Run(func(ctx context.Context) error {
		user, err := settingService.commonService.CommonGetUserByUserId(userid)
		if err != nil {
			return err
		}
		if !user.IsAdmin {
			return errors.New(commonModel.NO_PERMISSION_DENIED)
		}

		// 检查评论服务提供者是否有效
		if newSetting.Provider != string(commonModel.TWIKOO) &&
			newSetting.Provider != string(commonModel.ARTALK) &&
			newSetting.Provider != string(commonModel.WALINE) &&
			newSetting.Provider != string(commonModel.GISCUS) {
			return errors.New(commonModel.NO_SUCH_COMMENT_PROVIDER)
		}

		commentSetting := &model.CommentSetting{
			EnableComment: newSetting.EnableComment,
			Provider:      newSetting.Provider,
			CommentAPI:    httpUtil.TrimURL(newSetting.CommentAPI),
		}

		// 序列化为 JSON
		settingToJSON, err := jsonUtil.JSONMarshal(commentSetting)
		if err != nil {
			return err
		}

		if err := settingService.keyvalueRepository.UpdateKeyValue(ctx, commonModel.CommentSettingKey, string(settingToJSON)); err != nil {
			return err
		}

		return nil
	})
}

// GetS3Setting 获取 S3 存储设置
func (settingService *SettingService) GetS3Setting(userid uint, setting *model.S3Setting) error {
	return settingService.txManager.Run(func(ctx context.Context) error {
		s3Setting, err := settingService.keyvalueRepository.GetKeyValue(commonModel.S3SettingKey)
		if err != nil {
			// 数据库中不存在数据，手动添加初始数据
			setting.Enable = false
			setting.Provider = string(commonModel.AWS)
			setting.Endpoint = ""
			setting.AccessKey = ""
			setting.SecretKey = ""
			setting.BucketName = ""
			setting.Region = ""
			setting.UseSSL = false
			setting.CDNURL = ""
			setting.PathPrefix = ""
			setting.PublicRead = true

			// 序列化为 JSON
			settingToJSON, err := jsonUtil.JSONMarshal(setting)
			if err != nil {
				return err
			}
			if err := settingService.keyvalueRepository.AddKeyValue(ctx, commonModel.S3SettingKey, string(settingToJSON)); err != nil {
				return err
			}

			return nil
		}

		if err := jsonUtil.JSONUnmarshal([]byte(s3Setting.(string)), setting); err != nil {
			return err
		}

		// 如果用户未登录且不为管理员,则屏蔽 S3 设置的敏感信息
		if userid == authModel.NO_USER_LOGINED {
			setting.AccessKey = "******"
			setting.SecretKey = "******"
			setting.BucketName = "******"
			setting.Endpoint = "******"
		}

		return nil
	})
}

// UpdateS3Setting 更新 S3 存储设置
func (settingService *SettingService) UpdateS3Setting(userid uint, newSetting *model.S3SettingDto) error {
	return settingService.txManager.Run(func(ctx context.Context) error {
		user, err := settingService.commonService.CommonGetUserByUserId(userid)
		if err != nil {
			return err
		}
		if !user.IsAdmin {
			return errors.New(commonModel.NO_PERMISSION_DENIED)
		}

		s3Setting := &model.S3Setting{
			Enable:     newSetting.Enable,
			Provider:   newSetting.Provider,
			Endpoint:   httpUtil.TrimURL(newSetting.Endpoint),
			AccessKey:  newSetting.AccessKey,
			SecretKey:  newSetting.SecretKey,
			BucketName: newSetting.BucketName,
			Region:     newSetting.Region,
			UseSSL:     newSetting.UseSSL,
			CDNURL:     httpUtil.TrimURL(newSetting.CDNURL),
			PathPrefix: httpUtil.TrimURL(newSetting.PathPrefix),
			PublicRead: newSetting.PublicRead,
		}

		// 序列化为 JSON
		settingToJSON, err := jsonUtil.JSONMarshal(s3Setting)
		if err != nil {
			return err
		}

		if err := settingService.keyvalueRepository.UpdateKeyValue(ctx, commonModel.S3SettingKey, string(settingToJSON)); err != nil {
			return err
		}

		return nil
	})
}

// GetOAuth2Setting 获取 OAuth2 设置
func (settingService *SettingService) GetOAuth2Setting(
	userid uint,
	setting *model.OAuth2Setting,
	forInternal bool,
) error {
	return settingService.txManager.Run(func(ctx context.Context) error {
		if !forInternal {
			user, err := settingService.commonService.CommonGetUserByUserId(userid)
			if err != nil {
				return err
			}
			if !user.IsAdmin {
				return errors.New(commonModel.NO_PERMISSION_DENIED)
			}
		}

		oauthSetting, err := settingService.keyvalueRepository.GetKeyValue(commonModel.OAuth2SettingKey)
		if err != nil {
			// 数据库中不存在数据，手动添加初始数据
			setting.Enable = false
			setting.Provider = string(commonModel.OAuth2GITHUB)
			setting.ClientID = ""
			setting.ClientSecret = ""
			setting.AuthURL = "https://github.com/login/oauth/authorize"
			setting.TokenURL = "https://github.com/login/oauth/access_token"
			setting.UserInfoURL = "https://api.github.com/user"
			setting.RedirectURI = ""
			setting.Scopes = []string{
				"read:user",
			}

			// 序列化为 JSON
			settingToJSON, err := jsonUtil.JSONMarshal(setting)
			if err != nil {
				return err
			}
			if err := settingService.keyvalueRepository.AddKeyValue(ctx, commonModel.OAuth2SettingKey, string(settingToJSON)); err != nil {
				return err
			}

			return nil
		}

		if err := jsonUtil.JSONUnmarshal([]byte(oauthSetting.(string)), setting); err != nil {
			return err
		}

		return nil
	})
}

// UpdateOAuth2Setting 更新 OAuth2 设置
func (settingService *SettingService) UpdateOAuth2Setting(userid uint, newSetting *model.OAuth2SettingDto) error {
	return settingService.txManager.Run(func(ctx context.Context) error {
		user, err := settingService.commonService.CommonGetUserByUserId(userid)
		if err != nil {
			return err
		}
		if !user.IsAdmin {
			return errors.New(commonModel.NO_PERMISSION_DENIED)
		}

		oauthSetting := &model.OAuth2Setting{
			Enable:       newSetting.Enable,
			Provider:     newSetting.Provider,
			ClientID:     newSetting.ClientID,
			ClientSecret: newSetting.ClientSecret,
			AuthURL:      httpUtil.TrimURL(newSetting.AuthURL),
			TokenURL:     httpUtil.TrimURL(newSetting.TokenURL),
			UserInfoURL:  httpUtil.TrimURL(newSetting.UserInfoURL),
			RedirectURI:  httpUtil.TrimURL(newSetting.RedirectURI),
			Scopes:       newSetting.Scopes,
		}

		// 序列化为 JSON
		settingToJSON, err := jsonUtil.JSONMarshal(oauthSetting)
		if err != nil {
			return err
		}

		if err := settingService.keyvalueRepository.UpdateKeyValue(ctx, commonModel.OAuth2SettingKey, string(settingToJSON)); err != nil {
			return err
		}

		return nil
	})
}

// GetOAuth2Status 获取 OAuth2 状态
func (settingService *SettingService) GetOAuth2Status(status *model.OAuth2Status) error {
	var oauthSetting model.OAuth2Setting
	if err := settingService.GetOAuth2Setting(authModel.NO_USER_LOGINED, &oauthSetting, true); err != nil {
		return err
	}

	status.Enabled = oauthSetting.Enable
	status.Provider = oauthSetting.Provider

	return nil
}

// GetAllWebhooks 获取所有 Webhook
func (settingService *SettingService) GetAllWebhooks(userid uint) ([]model.Webhook, error) {
	return nil, nil
}

// DeleteWebhook 删除 Webhook
func (settingService *SettingService) DeleteWebhook(userid, id uint) error {
	return nil
}

// CreateWebhook 创建 Webhook
func (settingService *SettingService) CreateWebhook(userid uint, newWebhook *model.WebhookDto) error {
	return nil
}
