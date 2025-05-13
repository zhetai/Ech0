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
	c.JSON(http.StatusOK, services.UploadImage(c))
}

// GetConnect 处理 GET /connect 请求，获取 Connect 信息
func GetConnect(c *gin.Context) {
	// 调用 Service 层获取 Connect 信息
	connect, err := services.GetConnect()
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetConnectFailMessage))
		return
	}

	c.JSON(http.StatusOK, dto.OK(connect, models.GetConnectSuccessMessage))
}
