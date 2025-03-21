package services

import (
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/repository"
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
	if err != nil || !user.IsAdmin { // TODO : 目前只有管理员可以发布留言
		return err
	}
	message.Username = user.Username // 获取用户名
	return repository.CreateMessage(message)
}

// DeleteMessage 根据 ID 删除留言
func DeleteMessage(id uint) error {
	return repository.DeleteMessage(id)
}
