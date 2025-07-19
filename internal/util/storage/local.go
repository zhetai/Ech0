package util

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/lin-snow/ech0/internal/config"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
)

// UploadFileToLocal 根据文件类型上传文件到本地存储
func UploadFileToLocal(file *multipart.FileHeader, fileType commonModel.UploadFileType) (string, error) {
	switch fileType {
	case commonModel.ImageType:
		return UploadImageToLocal(file)
	case commonModel.AudioType:
		return UploadAudioToLocal(file)
	default:
		return "", errors.New(commonModel.FILE_TYPE_NOT_ALLOWED)
	}
}

// UploadImageToLocal 将图片上传到本地存储
func UploadImageToLocal(file *multipart.FileHeader) (string, error) {
	// 创建图片存储目录
	if err := createDirIfNotExist(config.Config.Upload.ImagePath); err != nil {
		return "", err
	}

	// 获取原始文件名和扩展名
	ext := filepath.Ext(file.Filename)
	baseName := strings.TrimSuffix(file.Filename, ext)

	// 使用 UUID 和原始文件名生成新的文件名
	newFileName := fmt.Sprintf("%s_%s%s", baseName, uuid.New().String(), ext)
	// 保存文件到指定目录
	savePath := filepath.Join(config.Config.Upload.ImagePath, newFileName)
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := src.Close(); closeErr != nil {
			log.Println("Failed to close file source:", closeErr)
		}
	}()

	if err = os.MkdirAll(filepath.Dir(savePath), 0750); err != nil {
		return "", err
	}

	out, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil {
			log.Println("Failed to close destination file:", closeErr)
		}
	}()

	if _, err = io.Copy(out, src); err != nil {
		return "", err
	}

	imageURL := fmt.Sprintf("/images/%s", newFileName)
	return imageURL, nil
}

// UploadAudioToLocal 将音频上传到本地存储
func UploadAudioToLocal(file *multipart.FileHeader) (string, error) {
	// 创建音频存储目录
	if err := createDirIfNotExist(config.Config.Upload.AudioPath); err != nil {
		return "", err
	}

	// 获取扩展名
	ext := filepath.Ext(file.Filename)

	// 重名音频文件名（暂时使用固定名字 music + 扩展名）
	newFileName := fmt.Sprintf("music%s", ext)
	savePath := filepath.Join(config.Config.Upload.AudioPath, newFileName)
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := src.Close(); closeErr != nil {
			log.Println("Failed to close file source:", closeErr)
		}
	}()

	if err = os.MkdirAll(filepath.Dir(savePath), 0750); err != nil {
		return "", err
	}

	out, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil {
			log.Println("Failed to close destination file:", closeErr)
		}
	}()

	if _, err = io.Copy(out, src); err != nil {
		return "", err
	}

	// 返回音频的 URL
	audioURL := fmt.Sprintf("/audios/%s", newFileName)
	return audioURL, nil
}

// DeleteFileFromLocal 删除本地文件
func DeleteFileFromLocal(filePath string) error {
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		// 只有当错误不是"文件不存在"时才返回错误
		return err
	}
	return nil
}
