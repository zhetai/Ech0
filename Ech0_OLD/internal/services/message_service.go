package services

import (
	"errors"

	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/repository"
	"github.com/lin-snow/ech0/pkg"
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
func GetMessagesByPage(page, pageSize int, search string, showPrivate bool) (dto.PageQueryResult, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 调用持久化层获取数据
	messages, total := repository.GetMessagesByPage(page, pageSize, search, showPrivate)

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

	// 检查Extension内容
	if message.Extension != "" && message.ExtensionType != "" {
		if message.ExtensionType == models.Extension_MUSIC {

		} else if message.ExtensionType == models.Extension_VIDEO {

		} else if message.ExtensionType == models.Extension_GITHUBPROJ {
			// 处理GitHub项目的链接
			message.Extension = pkg.TrimURL(message.Extension)
		} else if message.ExtensionType == models.Extension_WEBSITE {

		}
	} else {
		message.Extension = ""
		message.ExtensionType = ""
	}

	if message.ImageURL == "" {
		message.ImageSource = ""
	}

	message.Username = user.Username // 获取用户名
	for i := range message.Images {
		if message.Images[i].ImageURL == "" {
			message.Images[i].ImageSource = ""
		}
	}

	return repository.CreateMessage(message)
}

// DeleteMessage 根据 ID 删除留言
func DeleteMessage(id uint) error {
	// 检查该留言是否存在图片
	message, err := GetMessageByID(id, true)
	if err != nil {
		return err
	}
	if message == nil {
		return errors.New(models.MessageNotFoundMessage)
	}

	// 旧版本的图片删除
	// if message.ImageURL != "" {
	// 	// 构造图片 DTO
	// 	image := dto.ImageDto{
	// 		URL:    message.ImageURL,
	// 		SOURCE: message.ImageSource,
	// 	}
	// 	// 调用图片服务删除图片
	// 	if err := DeleteImage(image); err != nil {
	// 		return err
	// 	}
	// }

	// 删除留言中的图片
	if len(message.Images) > 0 {
		for _, img := range message.Images {
			if err := DeleteImage(dto.ImageDto{
				URL:    img.ImageURL,
				SOURCE: img.ImageSource,
			}); err != nil {
				return err
			}
		}
	}

	return repository.DeleteMessage(id)
}
