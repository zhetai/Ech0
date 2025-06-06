package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/config"
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

	// 调用 Service 层删除音频
	if err := services.DeleteAudio(); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.DeleteAudioSuccessMessage))
}

// 获取可播放的音乐的URL
func GetPlayMusic(c *gin.Context) {
	// 调用 Service 层获取可播放的音乐的URL
	musicURLs := services.GetPlayMusic()

	c.JSON(http.StatusOK, dto.OK(musicURLs, models.GetPlayMusicSuccessMessage))
}

// PlayMusic 控制器
func PlayMusic(c *gin.Context) {
	// 以文件流的形式返回音乐文件
	musicURL := services.GetPlayMusic()
	musicName := ""
	if musicURL != "" {
		// 只保留最后的文件名
		musicName = musicURL[len("/audios/"):]
	}

	// 获取音乐文件的路径
	musicPath := config.Config.Upload.AudioPath + musicName

	// 获取 Content-Type
	contentType := "audio/mpeg"
	if musicName[len(musicName)-4:] == ".flac" {
		contentType = "audio/flac"
	} else if musicName[len(musicName)-4:] == ".m4a" {
		contentType = "audio/mp4"
	}

	// 读取文件内容
	data, err := os.ReadFile(musicPath)
	if err != nil {
		c.String(500, "读取音乐文件失败")
		return
	}

	// 设置响应头
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	// 直接写文件内容，Gin 会自动关闭连接，不会长时间占用文件
	c.Data(http.StatusOK, contentType, data)
}
