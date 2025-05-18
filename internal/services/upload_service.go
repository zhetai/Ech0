package services

import (
	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/pkg"
)

// UploadFile 上传文件
func UploadFile(c *gin.Context, fileType models.FileType) dto.Result[string] {
	user, err := GetUserByID(c.MustGet("userid").(uint))
	if err != nil {
		return dto.Fail[string](err.Error())
	}

	if !user.IsAdmin {
		return dto.Fail[string](models.NoPermissionMessage)
	}

	// 调用 pkg 中的文件上传方法
	fileURL, err := pkg.UploadFile(c, fileType)
	if err != nil {
		return dto.Fail[string](err.Error())
	}

	return dto.OK(fileURL)
}
