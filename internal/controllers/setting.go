package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/services"
)

// GetSettings 处理 GET /settings 请求，获取系统设置
func GetSettings(c *gin.Context) {
	settings, err := services.GetSetting()
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetSettingsFailMessage))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, dto.OK(settings, models.GetSettingsSuccessMessage))
}

// 更新系统设置
func UpdateSettings(c *gin.Context) {
	// 获取当前用户 ID
	userID := c.MustGet("userid").(uint)

	// 检查用户是否为管理员
	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.UserNotFoundMessage))
		return
	}
	if !user.IsAdmin {
		c.JSON(http.StatusOK, dto.Fail[string](models.NoPermissionMessage))
		return
	}

	// 解析请求体中的参数
	var settings dto.SystemSettingDto
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	// 调用 Service 层更新系统设置
	var newSettings models.SystemSetting
	newSettings.SiteTitle = settings.SiteTitle
	newSettings.ServerName = settings.ServerName
	newSettings.ServerURL = settings.ServerURL
	newSettings.AllowRegister = settings.AllowRegister
	newSettings.ICPNumber = settings.ICPNumber
	newSettings.MetingAPI = settings.MetingAPI
	newSettings.CustomCSS = settings.CustomCSS
	newSettings.CustomJS = settings.CustomJS
	if err := services.UpdateSetting(newSettings); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, dto.OK[any](nil, models.UpdateSettingsSuccessMessage))
}
