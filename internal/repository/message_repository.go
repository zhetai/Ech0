package repository

import (
	"errors"
	"strings"

	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/models"
	"gorm.io/gorm"
)

// GetAllMessages 从数据库获取所有留言
func GetAllMessages(showPrivate bool) ([]models.Message, error) {
	var messages []models.Message

	// 是否将私密内容也查询出来
	if showPrivate {
		if err := database.DB.Order("created_at DESC").Find(&messages).Error; err != nil {
			return nil, err
		}
	} else {
		if err := database.DB.Where("private = ?", false).Find(&messages).Error; err != nil {
			return nil, err
		}
	}

	return messages, nil
}

// GetMessageByID 根据 ID 获取留言
func GetMessageByID(id uint, showPrivate bool) (*models.Message, error) {
	var message models.Message
	result := database.DB.First(&message, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 如果未找到记录，则返回 nil
		}
		return nil, result.Error // 其他错误返回
	}

	if !showPrivate && message.Private {
		return nil, nil
	}

	return &message, nil
}

// CreateMessage 保存一条留言
func CreateMessage(message *models.Message) error {
	// 防止 XSS 攻击，使用 bluemonday 清理器来清理 HTML 标签
	// p := bluemonday.UGCPolicy()                   // 创建一个新的 HTML 清理器
	// message.Content = p.Sanitize(message.Content) // 清理内容中的 HTML 标签

	message.Content = strings.TrimSpace(message.Content) // 去除内容前后的空格
	// message.Username = strings.TrimSpace(message.Username) // 去除用户名前后的空格

	if message.Content == "" && message.ImageURL == "" {
		return errors.New(models.CannotBeEmptyMessage) // 如果内容和图片链接都为空，则返回错误
	}

	result := database.DB.Create(message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteMessage 根据 ID 删除留言
func DeleteMessage(id uint) error {
	var message models.Message
	result := database.DB.Delete(&message, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // 如果没有找到记录
	}
	return nil
}
