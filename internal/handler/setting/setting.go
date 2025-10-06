package handler

import (
	"github.com/gin-gonic/gin"
	res "github.com/lin-snow/ech0/internal/handler/response"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/setting"
	service "github.com/lin-snow/ech0/internal/service/setting"
)

type SettingHandler struct {
	settingService service.SettingServiceInterface
}

// NewSettingHandler SettingHandler 的构造函数
func NewSettingHandler(settingService service.SettingServiceInterface) *SettingHandler {
	return &SettingHandler{
		settingService: settingService,
	}
}

// GetSettings 获取设置
//
// @Summary 获取设置
// @Description 获取系统的全局设置
// @Tags 系统设置
// @Accept json
// @Produce json
// @Success 200 {object} res.Response{data=model.SystemSetting} "获取设置成功"
// @Failure 200 {object} res.Response "获取设置失败"
// @Router /settings [get]
func (settingHandler *SettingHandler) GetSettings() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		var settings model.SystemSetting
		if err := settingHandler.settingService.GetSetting(&settings); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: settings,
			Msg:  commonModel.GET_SETTINGS_SUCCESS,
		}
	})

}

// UpdateSettings 更新设置
//
// @Summary 更新设置
// @Description 更新系统的全局设置
// @Tags 系统设置
// @Accept json
// @Produce json
// @Param settings body model.SystemSettingDto true "新的系统设置"
// @Success 200 {object} res.Response "更新设置成功"
// @Failure 200 {object} res.Response "更新设置失败"
// @Router /settings [put]
func (settingHandler *SettingHandler) UpdateSettings() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		// 解析请求体中的参数
		var newSettings model.SystemSettingDto
		if err := ctx.ShouldBindJSON(&newSettings); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		if err := settingHandler.settingService.UpdateSetting(userid, &newSettings); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.UPDATE_SETTINGS_SUCCESS,
		}
	})
}

// GetCommentSettings 获取评论设置
//
// @Summary 获取评论设置
// @Description 获取系统的评论相关设置
// @Tags 系统设置
// @Accept json
// @Produce json
// @Success 200 {object} res.Response{data=model.CommentSetting} "获取评论设置成功"
// @Failure 200 {object} res.Response "获取评论设置失败"
// @Router /comment/settings [get]
func (settingHandler *SettingHandler) GetCommentSettings() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		var commentSetting model.CommentSetting
		if err := settingHandler.settingService.GetCommentSetting(&commentSetting); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: commentSetting,
			Msg:  commonModel.GET_COMMENT_SETTINGS_SUCCESS,
		}
	})
}

// UpdateCommentSettings 更新评论设置
//
// @Summary 更新评论设置
// @Description 更新系统的评论相关设置
// @Tags 系统设置
// @Accept json
// @Produce json
// @Param commentSettings body model.CommentSettingDto true "新的评论设置"
// @Success 200 {object} res.Response "更新评论设置成功"
// @Failure 200 {object} res.Response "更新评论设置失败"
// @Router /comment/settings [put]
func (settingHandler *SettingHandler) UpdateCommentSettings() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		// 解析请求体中的参数
		var newCommentSettings model.CommentSettingDto
		if err := ctx.ShouldBindJSON(&newCommentSettings); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		if err := settingHandler.settingService.UpdateCommentSetting(userid, &newCommentSettings); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.UPDATE_COMMENT_SETTINGS_SUCCESS,
		}
	})
}

// GetS3Settings 获取 S3 存储设置
//
// @Summary 获取 S3 存储设置
// @Description 获取系统的 S3 存储相关设置
// @Tags 系统设置
// @Accept json
// @Produce json
// @Success 200 {object} res.Response{data=model.S3Setting} "获取 S3 存储设置成功"
// @Failure 200 {object} res.Response "获取 S3 存储设置失败"
// @Router /s3/settings [get]
func (settingHandler *SettingHandler) GetS3Settings() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		var s3Setting model.S3Setting
		if err := settingHandler.settingService.GetS3Setting(userid, &s3Setting); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: s3Setting,
			Msg:  commonModel.GET_S3_SETTINGS_SUCCESS,
		}
	})

}

// UpdateS3Settings 更新 S3 存储设置
//
// @Summary 更新 S3 存储设置
// @Description 更新系统的 S3 存储相关设置
// @Tags 系统设置
// @Accept json
// @Produce json
// @Param s3Settings body model.S3SettingDto true "新的 S3 存储设置"
// @Success 200 {object} res.Response "更新 S3 存储设置成功"
// @Failure 200 {object} res.Response "更新 S3 存储设置失败"
// @Router /s3/settings [put]
func (settingHandler *SettingHandler) UpdateS3Settings() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		// 解析请求体中的参数
		var newS3Settings model.S3SettingDto
		if err := ctx.ShouldBindJSON(&newS3Settings); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		if err := settingHandler.settingService.UpdateS3Setting(userid, &newS3Settings); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.UPDATE_S3_SETTINGS_SUCCESS,
		}
	})
}

// GetOAuth2Settings 获取 OAuth2 设置
//
// @Summary 获取 OAuth2 设置
// @Description 获取系统的 OAuth2 相关设置
// @Tags 系统设置
// @Accept json
// @Produce json
// @Success 200 {object} res.Response{data=model.OAuth2Setting} "获取 OAuth2 设置成功"
// @Failure 200 {object} res.Response "获取 OAuth2 设置失败"
// @Router /oauth2/settings [get]
func (settingHandler *SettingHandler) GetOAuth2Settings() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		var oauthSetting model.OAuth2Setting
		if err := settingHandler.settingService.GetOAuth2Setting(userid, &oauthSetting); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: oauthSetting,
			Msg:  commonModel.GET_OAUTH_SETTINGS_SUCCESS,
		}
	})
}

// UpdateOAuth2Settings 更新 OAuth2 设置
//
// @Summary 更新 OAuth2 设置
// @Description 更新系统的 OAuth2 相关设置
// @Tags 系统设置
// @Accept json
// @Produce json
// @Param oauthSettings body model.OAuth2SettingDto true "新的 OAuth 设置"
// @Success 200 {object} res.Response "更新 OAuth 设置成功"
// @Failure 200 {object} res.Response "更新 OAuth 设置失败"
// @Router /oauth2/settings [put]
func (settingHandler *SettingHandler) UpdateOAuth2Settings() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		// 解析请求体中的参数
		var newOAuthSettings model.OAuth2SettingDto
		if err := ctx.ShouldBindJSON(&newOAuthSettings); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		if err := settingHandler.settingService.UpdateOAuth2Setting(userid, &newOAuthSettings); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.UPDATE_OAUTH_SETTINGS_SUCCESS,
		}
	})
}