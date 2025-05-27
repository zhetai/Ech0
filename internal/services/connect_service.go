package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/repository"
	"github.com/lin-snow/ech0/pkg"
)

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
		connect.Logo = fmt.Sprintf("%s/api/%s", trimmedServerURL, logoPath)
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

func GetConnectsInfo() ([]models.Connect, error) {
	// 获取所有连接地址
	connects, err := repository.GetAllConnects()
	if err != nil {
		return nil, err
	}

	if len(connects) == 0 {
		return []models.Connect{}, nil
	}

	var connectList []models.Connect
	connectList = make([]models.Connect, 0, len(connects))

	var wg sync.WaitGroup
	connectChan := make(chan models.Connect, len(connects))

	seenURLs := make(map[string]struct{})
	var seenMutex sync.Mutex

	// 重试配置
	const maxRetries = 3
	const retryDelay = 500 * time.Millisecond

	for _, conn := range connects {
		wg.Add(1)
		go func(conn models.Connected) {
			defer wg.Done()

			url := pkg.TrimURL(conn.ConnectURL) + "/api/connect"

			var lastErr error

			// 重试循环
			for attempt := 0; attempt < maxRetries; attempt++ {
				resp, err := pkg.SendRequest(url, "GET", struct {
					Header  string
					Content string
				}{
					Header:  "Ech0_URL",
					Content: conn.ConnectURL,
				})

				if err != nil {
					lastErr = err
					// 如果不是最后一次重试，等待后继续
					if attempt < maxRetries-1 {
						time.Sleep(retryDelay * time.Duration(attempt+1))
						continue
					}
					// 最后一次重试也失败，放弃该连接
					fmt.Printf("[连接信息获取失败] 地址: %s，阶段: 发送请求，错误: %v\n", conn.ConnectURL, lastErr)
					return
				}

				var connectInfo dto.Result[models.Connect]
				if err := json.Unmarshal(resp, &connectInfo); err != nil {
					lastErr = err
					// JSON 解析失败，如果不是最后一次重试，等待后继续
					if attempt < maxRetries-1 {
						time.Sleep(retryDelay * time.Duration(attempt+1))
						continue
					}
					// 最后一次重试也失败，放弃该连接
					fmt.Printf("[连接信息获取失败] 地址: %s，阶段: 解析响应，错误: %v\n", conn.ConnectURL, lastErr)
					return
				}

				// 验证响应数据
				if connectInfo.Code != 1 || connectInfo.Data.ServerURL == "" {
					lastErr = fmt.Errorf("无效响应: code=%d, serverURL=%s", connectInfo.Code, connectInfo.Data.ServerURL)
					// 数据无效，如果不是最后一次重试，等待后继续
					if attempt < maxRetries-1 {
						time.Sleep(retryDelay * time.Duration(attempt+1))
						continue
					}
					// 最后一次重试也失败，放弃该连接
					fmt.Printf("[连接信息获取失败] 地址: %s，阶段: 校验数据，错误: %v\n", conn.ConnectURL, lastErr)
					return
				}

				// 成功获取有效数据，检查重复并发送
				seenMutex.Lock()
				if _, exists := seenURLs[connectInfo.Data.ServerURL]; exists {
					seenMutex.Unlock()
					return // 重复数据，直接返回
				}
				seenURLs[connectInfo.Data.ServerURL] = struct{}{}
				seenMutex.Unlock()

				connectChan <- connectInfo.Data
				return // 成功处理，退出重试循环
			}
		}(conn)
	}

	go func() {
		wg.Wait()
		close(connectChan)
	}()

	for connect := range connectChan {
		if connect.ServerURL == "" {
			continue
		}
		connectList = append(connectList, connect)
	}

	return connectList, nil
}
