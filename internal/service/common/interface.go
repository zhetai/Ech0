package service

import (
	"github.com/gin-gonic/gin"
	model "github.com/lin-snow/ech0/internal/model/common"
	userModel "github.com/lin-snow/ech0/internal/model/user"
	"mime/multipart"
)

type CommonServiceInterface interface {
	CommonGetUserByUserId(userId uint) (userModel.User, error)
	UploadImage(userid uint, file *multipart.FileHeader) (string, error)
	DeleteImage(userid uint, url, source string) error
	DirectDeleteImage(url, source string) error
	GetSysAdmin() (userModel.User, error)
	GetStatus() (model.Status, error)
	GetHeatMap() ([]model.Heapmap, error)
	GenerateRSS(ctx *gin.Context) (string, error)
	UploadMusic(userId uint, file *multipart.FileHeader) (string, error)
	DeleteMusic(userid uint) error
	GetPlayMusicUrl() string
	PlayMusic(ctx *gin.Context)
}
