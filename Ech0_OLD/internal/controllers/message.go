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
	pageQueryResult, err := services.GetMessagesByPage(pageRequest.Page, pageRequest.PageSize, pageRequest.Search, showPrivate)
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
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidIDMessage))
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
