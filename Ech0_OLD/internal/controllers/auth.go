package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/services"
)

// Login 处理 POST /login 请求，用户登录
func Login(c *gin.Context) {
	// 从请求体获取用户名和密码
	var user dto.LoginDto
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	// 调用 Service 层登录验证
	token, err := services.Login(user)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	// 返回成功响应，包含 JWT Token
	c.JSON(http.StatusOK, dto.OK(token, models.LoginSuccessMessage))
}

// Register 处理用户注册
func Register(c *gin.Context) {
	// 从请求体获取用户名和密码
	var user dto.RegisterDto
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	// 调用 Service 层注册用户
	if err := services.Register(user); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, dto.OK[any](nil, models.RegisterSuccessMessage))
}
