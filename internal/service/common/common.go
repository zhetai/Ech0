package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/lin-snow/ech0/internal/config"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	echoModel "github.com/lin-snow/ech0/internal/model/echo"
	userModel "github.com/lin-snow/ech0/internal/model/user"
	repository "github.com/lin-snow/ech0/internal/repository/common"
	mdUtil "github.com/lin-snow/ech0/internal/util/md"
	storageUtil "github.com/lin-snow/ech0/internal/util/storage"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type CommonService struct {
	commonRepository repository.CommonRepositoryInterface
}

func NewCommonService(commonRepository repository.CommonRepositoryInterface) CommonServiceInterface {
	return &CommonService{
		commonRepository: commonRepository,
	}
}

func (commonService *CommonService) CommonGetUserByUserId(userId uint) (userModel.User, error) {
	return commonService.commonRepository.GetUserByUserId(userId)
}

func (commonService *CommonService) UploadImage(userId uint, file *multipart.FileHeader) (string, error) {
	user, err := commonService.commonRepository.GetUserByUserId(userId)
	if err != nil {
		return "", err
	}
	if !user.IsAdmin {
		return "", errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 检查文件类型是否合法
	if !storageUtil.IsAllowedType(file.Header.Get("Content-Type"), config.Config.Upload.AllowedTypes) {
		return "", errors.New(commonModel.FILE_TYPE_NOT_ALLOWED)
	}

	// 检查文件大小是否合法
	if file.Size > int64(config.Config.Upload.ImageMaxSize) {
		return "", errors.New(commonModel.FILE_SIZE_EXCEED_LIMIT)
	}

	// 调用存储函数存储图片
	imageUrl, err := storageUtil.UploadFile(file, commonModel.ImageType, commonModel.LOCAL_FILE)
	if err != nil {
		return "", err
	}

	return imageUrl, nil
}

func (commonService *CommonService) DeleteImage(userid uint, url, source string) error {
	user, err := commonService.commonRepository.GetUserByUserId(userid)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 检查图片是否存在
	if url == "" {
		return errors.New(commonModel.IMAGE_NOT_FOUND)
	}

	if source == echoModel.ImageSourceLocal {
		// 获取图片名字（去除前面的/images/)
		imageName := url[len("/images/"):]

		// 构造图片路径
		imagePath := fmt.Sprintf("data/images/%s", imageName)

		// 删除图片
		return storageUtil.DeleteFileFromLocal(imagePath)
	} else if source == echoModel.ImageSourceURL {
		// 无需处理
	} else if source == echoModel.ImageSourceS3 {

	} else if source == echoModel.ImageSourceR2 {

	} else {
		// 未知图片来源按本地图片处理
		// 获取图片名字（去除前面的/images/)
		imageName := url[len("/images/"):]

		// 构造图片路径
		imagePath := fmt.Sprintf("data/images/%s", imageName)

		// 删除图片
		return storageUtil.DeleteFileFromLocal(imagePath)
	}

	return nil
}

func (commonService *CommonService) DirectDeleteImage(url, source string) error {
	// 检查图片是否存在
	if url == "" {
		return errors.New(commonModel.IMAGE_NOT_FOUND)
	}

	if source == echoModel.ImageSourceLocal {
		// 获取图片名字（去除前面的/images/)
		imageName := url[len("/images/"):]

		// 构造图片路径
		imagePath := fmt.Sprintf("data/images/%s", imageName)

		// 删除图片
		return storageUtil.DeleteFileFromLocal(imagePath)
	} else if source == echoModel.ImageSourceURL {
		// 无需处理
	} else if source == echoModel.ImageSourceS3 {

	} else if source == echoModel.ImageSourceR2 {

	} else {
		// 未知图片来源按本地图片处理
		// 获取图片名字（去除前面的/images/)
		imageName := url[len("/images/"):]

		// 构造图片路径
		imagePath := fmt.Sprintf("data/images/%s", imageName)

		// 删除图片
		return storageUtil.DeleteFileFromLocal(imagePath)
	}

	return nil
}

func (commonService *CommonService) GetSysAdmin() (userModel.User, error) {
	return commonService.commonRepository.GetSysAdmin()
}

func (commonService *CommonService) GetStatus() (commonModel.Status, error) {
	// 获取系统管理员信息
	sysuser, err := commonService.commonRepository.GetSysAdmin()
	if err != nil {
		return commonModel.Status{}, err
	}

	// 获取所有用户状态信息
	var users []commonModel.UserStatus
	allusers, err := commonService.commonRepository.GetAllUsers()
	if err != nil {
		return commonModel.Status{}, err
	}
	for _, user := range allusers {
		users = append(users, commonModel.UserStatus{
			UserID:   user.ID,
			UserName: user.Username,
			IsAdmin:  user.IsAdmin,
		})
	}

	status := commonModel.Status{}

	echos, err := commonService.commonRepository.GetAllEchos(true)
	if err != nil {
		return status, err
	}

	status.SysAdminID = sysuser.ID
	status.Username = sysuser.Username
	status.Logo = sysuser.Avatar
	status.Users = users
	status.TotalMessages = len(echos)

	return status, nil
}

func (commonService *CommonService) GetHeatMap() ([]commonModel.Heapmap, error) {
	// 获取当前日期
	today := time.Now()

	// 获取一个月前的日期
	oneMonthAgo := today.AddDate(0, -1, 0)

	// 格式化为YYYY-MM-DD
	startDate := oneMonthAgo.Format("2006-01-02") // 一个月前的日期
	endDate := today.Format("2006-01-02")         // 当前日期

	// 数据库查询 （只返回某天count >= 1的item）
	heatmapData, err := commonService.commonRepository.GetHeatMap(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 如果不足30天，补齐数据（date为缺的日期，count为0）
	// Create a map for quick lookup of existing heatmap data
	heatmapMap := make(map[string]commonModel.Heapmap)
	for _, item := range heatmapData {
		heatmapMap[item.Date] = item
	}

	var results [30]commonModel.Heapmap
	for i := 0; i < 30; i++ {
		// 计算日期 (from today back to 29 days ago)
		date := today.AddDate(0, 0, -i).Format("2006-01-02")
		resultIndex := 29 - i

		if item, ok := heatmapMap[date]; ok {
			// 找到数据，填充结果
			results[resultIndex] = item
		} else {
			// 未找到数据，填充默认值
			results[resultIndex] = commonModel.Heapmap{
				Date:  date,
				Count: 0,
			}
		}
	}

	return results[:], nil
}

func (commonService *CommonService) GenerateRSS(ctx *gin.Context) (string, error) {
	// 获取所有留言
	echos, err := commonService.commonRepository.GetAllEchos(false)
	if err != nil {
		return "", err
	}

	// 生成 RSS 订阅链接
	schema := "http"
	if ctx.Request.TLS != nil {
		schema = "https"
	}
	host := ctx.Request.Host
	feed := &feeds.Feed{
		Title: "Ech0s~",
		Link: &feeds.Link{
			Href: fmt.Sprintf("%s://%s/", schema, host),
		},
		Image: &feeds.Image{
			Url: fmt.Sprintf("%s://%s/favicon.ico", schema, host),
		},
		Description: "Ech0s~",
		Author: &feeds.Author{
			Name: "Ech0s~",
		},
		Updated: time.Now(),
	}

	for _, msg := range echos {
		renderedContent := mdUtil.MdToHTML([]byte(msg.Content))

		title := msg.Username + " - " + msg.CreatedAt.Format("2006-01-02")

		// 添加图片链接到正文前(scheme://host/api/ImageURL)
		if len(msg.Images) > 0 {
			for _, image := range msg.Images {
				// 根据图片来源生成链接
				if image.ImageSource == echoModel.ImageSourceLocal {
					imageURL := fmt.Sprintf("%s://%s/api%s", schema, host, image.ImageURL)
					renderedContent = append([]byte(fmt.Sprintf("<img src=\"%s\" alt=\"Image\" style=\"max-width:100%%;height:auto;\" />", imageURL)), renderedContent...)
				}
			}
		}

		item := &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s://%s/api/echo/%d", schema, host, msg.ID)},
			Description: string(renderedContent),
			Author: &feeds.Author{
				Name: msg.Username,
			},
			Created: msg.CreatedAt,
		}
		feed.Items = append(feed.Items, item)
	}

	atom, err := feed.ToAtom()
	if err != nil {
		return "", err
	}

	return atom, nil
}

func (commonService *CommonService) UploadMusic(userId uint, file *multipart.FileHeader) (string, error) {
	user, err := commonService.commonRepository.GetUserByUserId(userId)
	if err != nil {
		return "", err
	}
	if !user.IsAdmin {
		return "", errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 检查文件类型是否合法
	if !storageUtil.IsAllowedType(file.Header.Get("Content-Type"), config.Config.Upload.AllowedTypes) {
		return "", errors.New(commonModel.FILE_TYPE_NOT_ALLOWED)
	}

	// 检查文件大小是否合法
	if file.Size > int64(config.Config.Upload.AudioMaxSize) {
		return "", errors.New(commonModel.FILE_SIZE_EXCEED_LIMIT)
	}

	// 调用存储函数存储图片
	audioUrl, err := storageUtil.UploadFile(file, commonModel.AudioType, commonModel.LOCAL_FILE)
	if err != nil {
		return "", err
	}

	return audioUrl, nil
}

func (commonService *CommonService) DeleteMusic(userid uint) error {
	user, err := commonService.commonRepository.GetUserByUserId(userid)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 支持的音频格式
	audioFiles := []string{"music.flac", "music.m4a", "music.mp3"}

	for _, file := range audioFiles {
		audioPath := fmt.Sprintf("data/audios/%s", file)
		if storageUtil.FileExists(audioPath) {
			return storageUtil.DeleteFileFromLocal(audioPath)
		}
	}

	return nil
}

func (commonService *CommonService) GetPlayMusicUrl() string {
	// 支持的音频格式
	audioFiles := []string{"music.flac", "music.m4a", "music.mp3"}

	for _, file := range audioFiles {
		audioPath := fmt.Sprintf("data/audios/%s", file)
		if storageUtil.FileExists(audioPath) {
			return fmt.Sprintf("/audios/%s", file)
		}
	}

	// 没有找到音频文件
	return ""
}

func (commonService *CommonService) PlayMusic(ctx *gin.Context) {
	// 以文件流的形式返回音乐文件
	musicURL := commonService.GetPlayMusicUrl()
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
		ctx.String(500, "读取音乐文件失败")
		return
	}

	// 设置响应头
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	// 直接写文件内容，Gin 会自动关闭连接，不会长时间占用文件
	ctx.Data(http.StatusOK, contentType, data)
}
