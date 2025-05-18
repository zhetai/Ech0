package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/services"
)

// GetStatus 处理 GET /status 请求，获取服务器状态
func GetStatus(c *gin.Context) {
	// 检查系统是否存在管理员（第一次安装时）
	_, err := services.GetSysAdmin()
	if err != nil {
		c.JSON(http.StatusOK, dto.OKWithCode[any](nil, models.InitInstallCode, models.PleaseSignUpFirstMessage))
		return
	}

	// 调用 Service 层获取状态
	status, err := services.GetStatus()
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetStatusFailMessage))
		return
	}

	c.JSON(http.StatusOK, dto.OK(status, models.GetStatusSuccessMessage))
}

// GetHeapMap 处理 GET /heapmap 请求，获取热力图数据
func GetHeatMap(c *gin.Context) {
	// 调用 Service 层获取热力图数据
	heatMap, err := services.GetHeatMap()
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetHeatMapFailMessage))
		return
	}

	c.JSON(http.StatusOK, dto.OK(heatMap, models.GetHeatMapSuccessMessage))
}

// GenerateRSS 处理 GET /rss 请求，生成 RSS 订阅链接
func GenerateRSS(c *gin.Context) {

	atom, err := services.GenerateRSS(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail[string](models.GenerateRSSFailMessage))
		return
	}

	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", []byte(atom))
}

// UploadImage 控制器，调用 service 上传图片
func UploadImage(c *gin.Context) {
	c.JSON(http.StatusOK, services.UploadFile(c, models.ImageType))
}

// UploadAudio 控制器，调用 service 上传音频
func UploadAudio(c *gin.Context) {
	c.JSON(http.StatusOK, services.UploadFile(c, models.AudioType))
}

// DeleteImage 控制器，调用 service 删除图片
func DeleteImage(c *gin.Context) {
	// 检查用户是否为管理员
	user, err := services.GetUserByID(c.MustGet("userid").(uint))
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.UserNotFoundMessage))
		return
	}
	if !user.IsAdmin {
		c.JSON(http.StatusOK, dto.Fail[string](models.NoPermissionMessage))
		return
	}

	var image dto.ImageDto
	if err := c.ShouldBindJSON(&image); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	// 调用 Service 层删除图片
	if err := services.DeleteImage(image); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.DeleteImageSuccessMessage))
}

// DeleteAudio 控制器，调用 service 删除音频
func DeleteAudio(c *gin.Context) {

}
