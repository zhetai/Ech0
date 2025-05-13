package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/repository"
	"github.com/lin-snow/ech0/pkg"
)

func GenerateRSS(c *gin.Context) (string, error) {
	// 获取所有留言
	showPrivate := false
	messages, err := GetAllMessages(showPrivate)
	if err != nil {
		return "", err
	}

	// 生成 RSS 订阅链接
	schema := "http"
	if c.Request.TLS != nil {
		schema = "https"
	}
	host := c.Request.Host
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

	for _, msg := range messages {
		renderedContent := pkg.MdToHTML([]byte(msg.Content))

		title := msg.Username + " - " + msg.CreatedAt.Format("2006-01-02")

		// 添加图片链接到正文前(scheme://host/api/ImageURL)
		if msg.ImageURL != "" {
			image := fmt.Sprintf("%s://%s/api%s", schema, host, msg.ImageURL)
			renderedContent = append([]byte(fmt.Sprintf("<img src=\"%s\" alt=\"Image\" style=\"max-width:100%%;height:auto;\" />", image)), renderedContent...)
		}

		item := &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s://%s/api/messages/%d", schema, host, msg.ID)},
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

func GetStatus() (models.Status, error) {
	// 获取系统管理员信息
	sysuser, err := repository.GetSysAdmin()
	if err != nil {
		return models.Status{}, errors.New(models.UserNotFoundMessage)
	}

	// 获取所有用户状态信息
	var users []models.UserStatus
	allusers, err := repository.GetAllUsers()
	if err != nil {
		return models.Status{}, errors.New(models.GetAllUsersFailMessage)
	}
	for _, user := range allusers {
		users = append(users, models.UserStatus{
			UserID:   user.ID,
			UserName: user.Username,
			IsAdmin:  user.IsAdmin,
		})
	}

	status := models.Status{}

	messages, err := repository.GetAllMessages(true)
	if err != nil {
		return status, errors.New(models.GetAllMessagesFailMessage)
	}

	status.SysAdminID = sysuser.ID
	status.Username = sysuser.Username
	status.Logo = sysuser.Avatar
	status.Users = users
	status.TotalMessages = len(messages)

	return status, nil
}

func GetHeatMap() ([]models.Heapmap, error) {
	// 获取当前日期
	today := time.Now()

	// 获取一个月前的日期
	oneMonthAgo := today.AddDate(0, -1, 0)

	// 格式化为YYYY-MM-DD
	startDate := oneMonthAgo.Format("2006-01-02") // 一个月前的日期
	endDate := today.Format("2006-01-02")         // 当前日期

	// 数据库查询 （只返回某天count >= 1的item）
	heapmapData, err := repository.GetHeatMap(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 如果不足30天，补齐数据（date为缺的日期，count为0）
	// Create a map for quick lookup of existing heatmap data
	heapmapMap := make(map[string]models.Heapmap)
	for _, item := range heapmapData {
		heapmapMap[item.Date] = item
	}

	var results [30]models.Heapmap
	for i := 0; i < 30; i++ {
		// 计算日期 (from today back to 29 days ago)
		date := today.AddDate(0, 0, -i).Format("2006-01-02")
		resultIndex := 29 - i

		if item, ok := heapmapMap[date]; ok {
			// 找到数据，填充结果
			results[resultIndex] = item
		} else {
			// 未找到数据，填充默认值
			results[resultIndex] = models.Heapmap{
				Date:  date,
				Count: 0,
			}
		}
	}

	return results[:], nil
}

func GetConnect() (models.Connect, error) {
	var connect models.Connect

	// 获取系统设置
	setting, err := GetSetting()
	if err != nil {
		return connect, err
	}

	// 获取系统状态
	status, err := GetStatus()
	if err != nil {
		return connect, err
	}

	// 设置 Connect 信息
	connect.ServerName = setting.ServerName
	connect.ServerURL = setting.ServerURL
	connect.Ech0s = status.TotalMessages
	connect.SysUsername = status.Username

	// 处理 Logo URL，避免出现重复的斜杠
	trimmedServerURL := setting.ServerURL
	if len(trimmedServerURL) > 0 && trimmedServerURL[len(trimmedServerURL)-1] == '/' {
		trimmedServerURL = trimmedServerURL[:len(trimmedServerURL)-1]
	}

	if status.Logo != "" {
		// 如果 Logo URL 以 / 开头，去掉一个 /
		logoPath := status.Logo
		if len(logoPath) > 0 && logoPath[0] == '/' {
			logoPath = logoPath[1:]
		}
		connect.Logo = fmt.Sprintf("%s/%s", trimmedServerURL, logoPath)
	} else {
		connect.Logo = fmt.Sprintf("%s/favicon.svg", trimmedServerURL)
	}

	return connect, nil
}

func AddConnect(connected models.Connected) error {
	// 检查连接地址是否为空
	if connected.ConnectURL == "" {
		return errors.New(models.ConnectURLIsEmptyMessage)
	}

	// 去除连接地址前后的空格和斜杠
	connected.ConnectURL = pkg.TrimURL(connected.ConnectURL)

	// 检查连接地址是否已存在
	connectedList, err := repository.GetAllConnects()
	if err != nil {
		return err
	}

	// 检查连接地址是否已存在
	for _, conn := range connectedList {
		if conn.ConnectURL == connected.ConnectURL {
			return errors.New(models.ConnectAlreadyExistsMessage)
		}
	}

	// 添加连接地址
	if err := repository.CreateConnect(&connected); err != nil {
		return err
	}

	return nil
}

func GetConnects() ([]models.Connected, error) {
	// 获取所有连接地址
	connects, err := repository.GetAllConnects()
	if err != nil {
		return nil, err
	}

	// 如果没有找到，返回空切片
	if len(connects) == 0 {
		return []models.Connected{}, nil
	}

	// 返回查询到的 connects
	return connects, nil
}

func DeleteConnect(id uint) error {
	// 删除连接地址
	if err := repository.DeleteConnect(id); err != nil {
		return err
	}

	return nil
}
