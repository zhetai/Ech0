package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/services"
)

// 更新用户信息 (用户名、密码、头像)
func UpdateUser(c *gin.Context) {
	// 解析用户请求体中的参数
	var userdto dto.UserInfoDto
	if err := c.ShouldBindJSON(&userdto); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	// 获取当前用户 ID
	user, err := services.GetUserByID(c.MustGet("userid").(uint))
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.UserNotFoundMessage))
		return
	}

	// 调用 Service 层更新用户信息
	if err := services.UpdateUser(user, userdto); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.UpdateUserSuccessMessage))
}

// 更新用户权限
func UpdateUserAdmin(c *gin.Context) {
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

	// 从Query参数获取用户 ID
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	// 不能和当前用户 ID 一样，不能是系统管理员 ID 1，不能为空
	if err != nil || id == uint64(user.ID) || id == 1 {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidIDMessage))
		return
	}

	// 调用 Service 层更新用户权限
	if err := services.UpdateUserAdmin(uint(id)); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.UpdateUserSuccessMessage))
}

// 获取当前登录用户的信息
func GetUserInfo(c *gin.Context) {
	// 获取当前用户 ID
	userID := c.MustGet("userid").(uint)

	// 调用 Service 层获取用户信息
	user, err := services.GetUserByID(userID)
	user.Password = "" // 不返回密码
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.UserNotFoundMessage))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, dto.OK(user, models.QuerySuccessMessage))
}
