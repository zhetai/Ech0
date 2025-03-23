package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/repository"
	"github.com/lin-snow/ech0/pkg"
)

// GetAllMessages 封装业务逻辑，获取所有留言
func GetAllMessages() ([]models.Message, error) {
	return repository.GetAllMessages()
}

// GetMessageByID 根据 ID 获取留言
func GetMessageByID(id uint) (*models.Message, error) {
	return repository.GetMessageByID(id)
}

// GetMessagesByPage 分页获取留言
func GetMessagesByPage(page, pageSize int) (dto.PageQueryResult, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 查询数据库
	var messages []models.Message
	var total int64

	database.DB.Model(&models.Message{}).Count(&total)
	database.DB.Limit(pageSize).Offset(offset).Order("created_at DESC").Find(&messages)

	// 返回结果
	var PageQueryResult dto.PageQueryResult
	PageQueryResult.Total = total
	PageQueryResult.Items = messages

	return PageQueryResult, nil
}

// CreateMessage 发布一条留言
func CreateMessage(message *models.Message) error {
	user, err := GetUserByID(message.UserID)
	if err != nil {
		return err
	}

	if !user.IsAdmin {
		return errors.New(models.NoPermissionMessage)
	}

	message.Username = user.Username // 获取用户名
	return repository.CreateMessage(message)
}

// DeleteMessage 根据 ID 删除留言
func DeleteMessage(id uint) error {
	return repository.DeleteMessage(id)
}

func GenerateRSS(c *gin.Context) (string, error) {
	// 获取所有留言
	messages, err := GetAllMessages()
	if err != nil {
		return "", err
	}

	// 生成 RSS 订阅链接
	schema := "http"
	if c.Request.TLS != nil {
		schema = "https"
	}
	host := c.Request.Host
	feed := &feeds.Feed{
		Title: "Ech0s~",
		Link: &feeds.Link{
			Href: fmt.Sprintf("%s://%s/", schema, host),
		},
		Image: &feeds.Image{
			Url: fmt.Sprintf("%s://%s/favicon.ico", schema, host),
		},
		Description: "Ech0s~",
		Author: &feeds.Author{
			Name: "Ech0s~",
		},
		Updated: time.Now(),
	}

	for _, msg := range messages {
		renderedContent := pkg.MdToHTML([]byte(msg.Content))

		title := msg.Username + " - " + msg.CreatedAt.Format("2006-01-02")

		// 添加图片链接到正文前(scheme://host/api/ImageURL)
		if msg.ImageURL != "" {
			image := fmt.Sprintf("%s://%s/api%s", schema, host, msg.ImageURL)
			renderedContent = append([]byte(fmt.Sprintf("<img src=\"%s\" alt=\"Image\" style=\"max-width:100%%;height:auto;\" />", image)), renderedContent...)
		}

		item := &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s://%s/api/messages/%d", schema, host, msg.ID)},
			Description: string(renderedContent),
			Author: &feeds.Author{
				Name: msg.Username,
			},
			Created: msg.CreatedAt,
		}
		feed.Items = append(feed.Items, item)
	}

	atom, err := feed.ToAtom()
	if err != nil {
		return "", err
	}

	return atom, nil
}
