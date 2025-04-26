package services

import (
	"errors"

	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/repository"
)

// GetAllMessages 封装业务逻辑，获取所有留言
func GetAllMessages(showPrivate bool) ([]models.Message, error) {
	return repository.GetAllMessages(showPrivate)
}

// GetMessageByID 根据 ID 获取留言
func GetMessageByID(id uint, showPrivate bool) (*models.Message, error) {
	return repository.GetMessageByID(id, showPrivate)
}

// GetMessagesByPage 分页获取留言
func GetMessagesByPage(page, pageSize int, showPrivate bool) (dto.PageQueryResult, error) {
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

	if showPrivate {
		// 如果是管理员，则不需要过滤私密留言
		database.DB.Model(&models.Message{}).Count(&total)
		database.DB.Limit(pageSize).Offset(offset).Order("created_at DESC").Find(&messages)
	} else {
		// 如果不是管理员，则只查询公开的留言
		database.DB.Model(&models.Message{}).Where("private = ?", false).Count(&total)
		database.DB.Limit(pageSize).Offset(offset).Where("private = ?", false).Order("created_at DESC").Find(&messages)
	}

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
