package handler

import (
	"strings"

	"github.com/gin-gonic/gin"

	res "github.com/lin-snow/ech0/internal/handler/response"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	service "github.com/lin-snow/ech0/internal/service/backup"
	jwtUtil "github.com/lin-snow/ech0/internal/util/jwt"
)

type BackupHandler struct {
	backupService service.BackupServiceInterface
}

// NewBackupHandler BackupHandler 的构造函数
func NewBackupHandler(backupService service.BackupServiceInterface) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
	}
}

// Backup 执行数据备份
//
// @Summary 执行数据备份
// @Description 用户触发数据备份操作，成功后返回备份成功信息
// @Tags 系统备份
// @Accept json
// @Produce json
// @Success 200 {object} res.Response "备份成功"
// @Failure 200 {object} res.Response "备份失败"
// @Router /backup [get]
func (backupHandler *BackupHandler) Backup() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		userId := ctx.MustGet("userid").(uint)
		if err := backupHandler.backupService.Backup(userId); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.BACKUP_SUCCESS,
		}
	})
}

// ExportBackup 导出数据备份
//
// @Summary 导出数据备份
// @Description 用户导出备份文件，成功后触发文件下载
// @Tags 系统备份
// @Accept json
// @Produce application/octet-stream
// @Success 200 {object} res.Response "导出备份成功，返回文件下载"
// @Failure 200 {object} res.Response "导出备份失败"
// @Router /backup/export [get]
func (backupHandler *BackupHandler) ExportBackup() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		token := ctx.Query("token")
		if token == "" {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
			}
		}

		token = strings.Trim(token, `"`) // 去掉可能的双引号

		// 使用 JWT Util进行处理
		if _, err := jwtUtil.ParseToken(token); err != nil {
			return res.Response{
				Msg: commonModel.TOKEN_NOT_VALID,
				Err: err,
			}
		}

		// userId := ctx.MustGet("userid").(uint)
		if err := backupHandler.backupService.ExportBackup(ctx); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.EXPORT_BACKUP_SUCCESS,
		}
	})
}
