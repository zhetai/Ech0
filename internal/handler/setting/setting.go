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

func NewSettingHandler(settingService service.SettingServiceInterface) *SettingHandler {
	return &SettingHandler{
		settingService: settingService,
	}
}

// GetSettings 获取设置
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
