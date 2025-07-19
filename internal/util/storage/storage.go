package util

import (
	"errors"
	"mime/multipart"
	"os"

	commonModel "github.com/lin-snow/ech0/internal/model/common"
)

// UploadFile 根据文件类型和存储类型上传文件
func UploadFile(file *multipart.FileHeader, fileType commonModel.UploadFileType, storageType commonModel.FileStorageType) (string, error) {
	if file == nil {
		return "", errors.New(commonModel.NO_FILE_UPLOAD_ERROR)
	}

	switch storageType {
	case commonModel.LOCAL_FILE:
		return UploadFileToLocal(file, fileType)
	case commonModel.S3_FILE:
		// TODO: Implement S3 file upload
	case commonModel.R2_FILE:
		// TODO: Implement R2 file upload
	default:
		return "", errors.New(commonModel.NO_FILE_STORAGE_ERROR)
	}

	return "", errors.New(commonModel.NO_FILE_STORAGE_ERROR)
}

// IsAllowedType 检查Content-Type是否在允许的类型列表中
func IsAllowedType(contentType string, allowedTypes []string) bool {
	for _, allowed := range allowedTypes {
		if contentType == allowed {
			return true
		}
	}
	return false
}

// createDirIfNotExist 创建目录如果不存在
func createDirIfNotExist(imagePath string) error {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		if err := os.MkdirAll(imagePath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// FileExists 文件是否存在
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
