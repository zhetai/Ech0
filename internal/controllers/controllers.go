package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/services"
	"gorm.io/gorm"
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

// GetMessages 处理 GET /messages 请求，返回所有留言
func GetMessages(c *gin.Context) {
	messages, err := services.GetAllMessages(false)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetAllMessagesFailMessage))
		return
	}
	c.JSON(http.StatusOK, dto.OK(messages, models.GetAllMessagesSuccess))
}

// GetMessage 处理 GET /messages/:id 请求，获取留言详情
func GetMessage(c *gin.Context) {
	// 从 URL 参数获取留言 ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidIDMessage))
		return
	}

	// 获取当前用户 ID (为0表示未登录，不返回隐私数据,不为零时则必须为管理员)
	userID := c.MustGet("userid").(uint)
	showPrivate := false
	if userID != 0 {
		// 获取当前用户信息
		user, err := services.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusOK, dto.Fail[string](models.UserNotFoundMessage))
			return
		}
		// 如果是管理员，则可以查看所有留言
		if user.IsAdmin {
			showPrivate = true
		}
	}

	// 调用 Service 层根据 ID 获取留言
	message, err := services.GetMessageByID(uint(id), showPrivate)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetMessageByIDFailMessage))
		return
	}

	if message == nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.MessageNotFoundMessage))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, dto.OK(message, models.GetMessageByIDSuccess))
}

// GetMessagesByPage 处理 POST /messages/page 请求，分页获取留言
func GetMessagesByPage(c *gin.Context) {
	// 从请求体获取分页参数
	var pageRequest dto.PageQueryDto
	if err := c.ShouldBindJSON(&pageRequest); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	// 获取当前用户 ID (为0表示未登录，不返回隐私数据,不为零时则必须为管理员)
	userID := c.MustGet("userid").(uint)
	showPrivate := false

	if userID != 0 {
		// 获取当前用户信息
		user, err := services.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusOK, dto.Fail[string](models.UserNotFoundMessage))
			return
		}
		// 如果是管理员，则可以查看所有留言
		if user.IsAdmin {
			showPrivate = true
		}
	}

	// 调用 Service 层获取分页留言
	pageQueryResult, err := services.GetMessagesByPage(pageRequest.Page, pageRequest.PageSize, showPrivate)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, dto.OK(pageQueryResult, models.GetMessagesByPageSuccess))
}

// PostMessage 处理 POST /messages 请求，发布留言
func PostMessage(c *gin.Context) {
	var message models.Message

	// 绑定请求体到 message 对象
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	// 调用 Service 层保存留言
	message.UserID = c.MustGet("userid").(uint)
	if err := services.CreateMessage(&message); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, dto.OK(message, models.PostMessageSuccess))
}

// GetStatus 处理 GET /status 请求，获取服务器状态
func GetStatus(c *gin.Context) {
	// 调用 Service 层获取状态
	status, err := services.GetStatus()
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetStatusFailMessage))
		return
	}

	c.JSON(http.StatusOK, dto.OK(status, models.GetStatusSuccessMessage))
}

// DeleteMessage 处理 DELETE /messages/:id 请求，删除留言
func DeleteMessage(c *gin.Context) {
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
		c.JSON(http.StatusOK, dto.Fail[string]("models.InvalidIDMessage"))
		return
	}

	// 调用 Service 层删除留言
	if err := services.DeleteMessage(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, dto.Fail[string](models.MessageNotFoundMessage))
		} else {
			c.JSON(http.StatusOK, dto.Fail[string](models.DeleteFailMessage))
		}
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, dto.OK[interface{}](nil, models.DeleteSuccessMessage))
}

func GenerateRSS(c *gin.Context) {

	atom, err := services.GenerateRSS(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail[string](models.GenerateRSSFailMessage))
		return
	}

	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", []byte(atom))
}

// 更新用户信息 （用户名 、 密码）
func UpdateUser(c *gin.Context) {
	// 检查用户是否为管理员
	user, err := services.GetUserByID(c.MustGet("userid").(uint))
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.UserNotFoundMessage))
		return
	}

	// 解析请求体中的参数
	var userdto dto.UserInfoDto
	if err := c.ShouldBindJSON(&userdto); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	// 调用 Service 层更新用户信息
	if err := services.UpdateUser(user, userdto); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.UpdateUserSuccessMessage))
}

// 更改密码
func ChangePassword(c *gin.Context) {
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

	// 调用 Service 层更改密码
	if err := services.ChangePassword(user, userdto); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.ChangePasswordSuccessMessage))
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

// 更改系统设置 （是否允许注册）
// func UpdateSetting(c *gin.Context) {
// 	// 解析请求体中的设置参数

// }
