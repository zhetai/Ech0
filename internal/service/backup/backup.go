package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/backup"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	commonService "github.com/lin-snow/ech0/internal/service/common"
)

type BackupService struct {
	commonService commonService.CommonServiceInterface
}

func NewBackupService(commonService commonService.CommonServiceInterface) BackupServiceInterface {
	return &BackupService{
		commonService: commonService,
	}
}

// Backup 执行备份
func (backupService *BackupService) Backup(userid uint) error {
	user, err := backupService.commonService.CommonGetUserByUserId(userid)
	if err != nil {
		return err
	}

	if !user.IsAdmin {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 执行备份
	if _, _, err := backup.ExecuteBackup(); err != nil {
		return err
	}

	return nil
}

// ExportBackup 导出备份
func (backupService *BackupService) ExportBackup(ctx *gin.Context) error {
	// 导出备份
	// 1. 先备份
	var backupFilePath string // 备份文件路径
	var err error
	backupFilePath, _, err = backup.ExecuteBackup()
	if err != nil {
		return err
	}

	// 2. 计算文件大小
	fileInfo, err := os.Stat(backupFilePath)
	if err != nil {
		return err
	}

	// 设置响应头
	filename := fmt.Sprintf("ech0-backup-%s.zip", time.Now().Format("2006-01-02-150405"))

	// 设置响应头的顺序很重要
	ctx.Writer.Header().Set("Content-Type", "application/zip")
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	ctx.Writer.Header().Set("Accept-Ranges", "bytes")
	ctx.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// ✅ 立即刷新响应头到客户端
	ctx.Writer.WriteHeader(200)

	// 使用 Gin 的内置方法，支持 Range 请求
	ctx.File(backupFilePath)
	return nil
}
