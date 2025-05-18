package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lin-snow/ech0/config"
	"github.com/lin-snow/ech0/internal/models"
)

// 通用上传函数
func UploadFile(c *gin.Context, fileType models.FileType) (string, error) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		return "", errors.New(models.NotUploadFileErrorMessage)
	}

	// 检查文件扩展名
	allowedExtensions := config.Config.Upload.AllowedTypes
	// 检查文件类型是否合法
	if !isAllowedType(file.Header.Get("Content-Type"), allowedExtensions) {
		return "", errors.New(models.NotSupportedFileTypeErrorMessage)
	}

	// 获取原始文件名和扩展名
	ext := filepath.Ext(file.Filename)
	baseName := strings.TrimSuffix(file.Filename, ext)

	// 根据文件类型分为音频和图片
	var savePath string
	if strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		// 创建存储图片的目录（如果没有的话）
		if err := createImageDirIfNotExist(config.Config.Upload.ImagePath); err != nil {
			return "", err
		}

		// 检查文件大小
		if file.Size > int64(config.Config.Upload.ImageMaxSize) {
			return "", errors.New(models.ImageSizeLimitErrorMessage + strconv.Itoa(config.Config.Upload.ImageMaxSize/1024/1024) + "MB")
		}

		// 使用 UUID 和原始文件名生成新的文件名
		newFileName := fmt.Sprintf("%s_%s%s", baseName, uuid.New().String(), ext)
		// 保存文件到指定目录
		savePath = filepath.Join(config.Config.Upload.ImagePath, newFileName)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			fmt.Println(savePath)
			return "", errors.New(models.ImageUploadErrorMessage)
		}

		// 返回图片的 URL
		imageURL := fmt.Sprintf("/images/%s", newFileName)
		return imageURL, nil
	} else if strings.HasPrefix(file.Header.Get("Content-Type"), "audio/") {
		// 创建存储音频的目录（如果没有的话）
		if err := createImageDirIfNotExist(config.Config.Upload.AudioPath); err != nil {
			return "", err
		}

		// 检查文件大小
		if file.Size > int64(config.Config.Upload.AudioMaxSize) {
			return "", errors.New(models.AudioSizeLimitErrorMessage + strconv.Itoa(config.Config.Upload.AudioMaxSize/1024/1024) + "MB")
		}

		// 重名音频文件名（暂时使用固定名字 music + 扩展名）
		newFileName := fmt.Sprintf("music%s", ext)
		savePath = filepath.Join(config.Config.Upload.AudioPath, newFileName)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			fmt.Println(savePath)
			return "", errors.New(models.AudioUploadErrorMessage)
		}

		// 返回音频的 URL
		audioURL := fmt.Sprintf("/audios/%s", newFileName)
		return audioURL, nil
	} else {
		return "", errors.New(models.NotSupportedFileTypeErrorMessage)
	}
}

// 检查文件类型是否合法
func isAllowedType(contentType string, allowedTypes []string) bool {
	fmt.Println(contentType)
	fmt.Println(allowedTypes)
	for _, allowed := range allowedTypes {
		if contentType == allowed {
			return true
		}
	}
	return false
}

// 创建存储图片的目录
func createImageDirIfNotExist(imagePath string) error {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		if err := os.MkdirAll(imagePath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
