package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/services"
)

// GetConnect 处理 GET /connect 请求，获取 Connect 信息
func GetConnect(c *gin.Context) {
	//

	// 调用 Service 层获取 Connect 信息
	connect, err := services.GetConnect()
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetConnectFailMessage))
		return
	}

	c.JSON(http.StatusOK, dto.OK(connect, models.GetConnectSuccessMessage))
}

func AddConnect(c *gin.Context) {
	// 检查用户是否为管理员
	user, err := services.GetUserByID(c.MustGet("userid").(uint))
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.UserNotFoundMessage))
		return
	}
	if !user.IsAdmin {
		c.JSON(http.StatusOK, dto.Fail[string](models.NoPermissionMessage))
		return
	}

	var connected models.Connected
	if err := c.ShouldBindJSON(&connected); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	// 调用 Service 层添加 Connect 信息
	if err := services.AddConnect(connected); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.AddConnectSuccessMessage))
}

func GetConnects(c *gin.Context) {
	// 调用 Service 层获取 Connect 列表
	connects, err := services.GetConnects()
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetConnectsFailMessage))
		return
	}

	c.JSON(http.StatusOK, dto.OK(connects, models.GetConnectsSuccessMessage))
}

func DeleteConnect(c *gin.Context) {
	// 检查用户是否为管理员
	user, err := services.GetUserByID(c.MustGet("userid").(uint))
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.UserNotFoundMessage))
		return
	}
	if !user.IsAdmin {
		c.JSON(http.StatusOK, dto.Fail[string](models.NoPermissionMessage))
		return
	}

	// 从 URL 参数获取留言 ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidIDMessage))
		return
	}

	// 调用 Service 层删除 Connect 信息
	if err := services.DeleteConnect(uint(id)); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.DeleteConnectSuccessMessage))
}

// GetConnectsInfo
func GetConnectsInfo(c *gin.Context) {
	// 调用 Service 层获取 Connect 信息
	connects, err := services.GetConnectsInfo()
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK(connects, models.GetConnectsInfoSuccessMessage))
}
