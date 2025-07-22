package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	backup "github.com/lin-snow/ech0/internal/backup"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	commonService "github.com/lin-snow/ech0/internal/service/common"
)

type BackupService struct {
	commonService  commonService.CommonServiceInterface
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
func (backupService *BackupService) ExportBackup(userid uint, ctx *gin.Context) error {
	user, err := backupService.commonService.CommonGetUserByUserId(userid)
	if err != nil {
		return err
	}

	if !user.IsAdmin {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 导出备份
	// 1. 先备份
	var backupFilePath string // 备份文件路径
	var backupFileName string // 备份文件名
	backupFilePath, backupFileName, err = backup.ExecuteBackup();
	if err != nil {
		return err
	}

	// 2. 再导出
	ctx.Header("Content-Disposition", "attachment; filename="+backupFileName)
	ctx.Header("Content-Type", "application/zip")
	ctx.File(backupFilePath)

	return nil
}