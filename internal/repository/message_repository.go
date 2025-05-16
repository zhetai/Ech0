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

// GetMessagesByPage 分页获取留言
func GetMessagesByPage(page, pageSize int, search string, showPrivate bool) ([]models.Message, int64) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询数据库
	var messages []models.Message
	var total int64

	query := database.DB.Model(&models.Message{})

	// 如果 search 不为空，添加模糊查询条件
	if search != "" || len(search) > 0 {
		searchPattern := "%" + search + "%" // 模糊匹配模式
		query = query.Where("content LIKE ?", searchPattern)
	}

	// 如果不是管理员，过滤私密留言
	if !showPrivate {
		query = query.Where("private = ?", false)
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	query.Limit(pageSize).Offset(offset).Order("created_at DESC").Find(&messages)

	// 返回结果
	return messages, total
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

	if message.Content == "" && message.ImageURL == "" && (message.Extension == "" || message.ExtensionType == "") {
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
