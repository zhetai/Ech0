package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
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
	messages, err := services.GetAllMessages()
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

	// 调用 Service 层根据 ID 获取留言
	message, err := services.GetMessageByID(uint(id))
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

	// 调用 Service 层获取分页留言
	pageQueryResult, err := services.GetMessagesByPage(pageRequest.Page, pageRequest.PageSize)
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
	var messages []models.Message
	// 调用 Service 层获取所有留言
	messages, err := services.GetAllMessages()
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[any](models.GetAllMessagesFailMessage))
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	feed := &feeds.Feed{
		Title:       "Latest Messages",
		Link:        &feeds.Link{Href: fmt.Sprintf("%s://%s/", scheme, host)},
		Description: "RSS feed for the latest messages",
		Author:      &feeds.Author{Name: "Ech0s~"},
		Created:     time.Now(),
	}

	for _, msg := range messages {
		title := truncate(msg.Content, 50)
		item := &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s://%s/api/messages/%d", scheme, host, msg.ID)},
			Description: msg.Content,
			Author:      &feeds.Author{Name: msg.Username},
			Created:     msg.CreatedAt,
		}
		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate RSS feed"})
		return
	}

	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", []byte(rss))
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}
