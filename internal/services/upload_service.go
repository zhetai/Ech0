package services

import (
	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/pkg"
)

// UploadImage 上传图片
func UploadImage(c *gin.Context) dto.Result[string] {
	user, err := GetUserByID(c.MustGet("userid").(uint))
	if err != nil {
		return dto.Fail[string](err.Error())
	}

	if !user.IsAdmin {
		return dto.Fail[string](models.NoPermissionMessage)
	}

	// 调用 pkg 中的图片上传方法
	imageURL, err := pkg.UploadImage(c)
	if err != nil {
		return dto.Fail[string](err.Error())
	}

	return dto.OK(imageURL)
}
