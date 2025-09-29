package repository

import (
	model "github.com/lin-snow/ech0/internal/model/common"
	echoModel "github.com/lin-snow/ech0/internal/model/echo"
	userModel "github.com/lin-snow/ech0/internal/model/user"
)

type CommonRepositoryInterface interface {
	// GetUserByUserId 根据用户ID获取用户信息
	GetUserByUserId(userid uint) (userModel.User, error)

	// GetSysAdmin 获取系统管理员信息
	GetSysAdmin() (userModel.User, error)

	// GetAllUsers 获取所有用户信息
	GetAllUsers() ([]userModel.User, error)

	// GetAllEchos 获取所有Echo
	GetAllEchos(showPrivate bool) ([]echoModel.Echo, error)

	// GetHeatMap 获取热力图数据
	GetHeatMap(startDate, endDate string) ([]model.Heatmap, error)

	// SaveTempFile 保存临时文件记录
	SaveTempFile(file model.TempFile) error

	// DeleteTempFile 删除临时文件记录
	DeleteTempFile(id uint) error

	// DeleteTempFilePermanently 永久删除临时文件记录
	DeleteTempFilePermanently(id uint) error

	// DeleteTempFileByObjectKey 根据对象键删除临时文件记录
	DeleteTempFileByObjectKey(objectKey string) error

	// GetAllTempFiles 获取所有未删除的临时文件
	GetAllTempFiles() ([]model.TempFile, error)

	// UpdateTempFileAccessTime 更新临时文件的最后访问时间
	UpdateTempFileAccessTime(id uint, accessTime int64) error
}
