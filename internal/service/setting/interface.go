package service

import model "github.com/lin-snow/ech0/internal/model/setting"

type SettingServiceInterface interface {
	// GetSetting 获取设置
	GetSetting(setting *model.SystemSetting) error

	// UpdateSetting 更新设置
	UpdateSetting(userid uint, newSetting *model.SystemSettingDto) error

	// GetCommentSetting 获取评论设置
	GetCommentSetting(setting *model.CommentSetting) error

	// UpdateCommentSetting 更新评论设置
	UpdateCommentSetting(userid uint, newSetting *model.CommentSettingDto) error

	// GetS3Setting 获取 S3 存储设置
	GetS3Setting(userid uint, setting *model.S3Setting) error

	// UpdateS3Setting 更新 S3 存储设置
	UpdateS3Setting(userid uint, newSetting *model.S3SettingDto) error
}
