package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

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

	// 如果没有找到，返回空切片
	if len(connects) == 0 {
		return []models.Connect{}, nil
	}

	var connectList []models.Connect
	// 预分配切片容量，减少动态扩容
	connectList = make([]models.Connect, 0, len(connects))

	var wg sync.WaitGroup
	connectChan := make(chan models.Connect, len(connects))

	// 使用map记录已添加的连接，防止重复
	// 使用结构体存储连接和优先级信息
	type ConnectEntry struct {
		connect  models.Connect
		priority int // 优先级：1=非特殊connect，2=特殊connect
	}

	seenURLs := make(map[string]ConnectEntry)
	var seenMutex sync.Mutex // 保护map的并发访问

	// 遍历连接地址，获取每个连接的状态
	for _, conn := range connects {
		wg.Add(1)
		go func(conn models.Connected) {
			defer wg.Done()

			url := pkg.TrimURL(conn.ConnectURL) + "/api/connect"
			resp, err := pkg.SendRequest(url, "GET", struct {
				Header  string
				Content string
			}{
				Header:  "Ech0_URL",
				Content: conn.ConnectURL,
			})
			if err != nil {
				// 处理请求错误
				return
			}

			var connectInfo dto.Result[models.Connect]
			if err := json.Unmarshal(resp, &connectInfo); err != nil {
				// 解析失败，抛弃该实例的数据
				// 检查是否为特殊connect,即为[]models.Connect
				var specialConnects []models.Connect
				if err := json.Unmarshal(resp, &specialConnects); err != nil {
					// 解析失败，抛弃该实例的数据
					return
				}

				// 处理特殊connect数组 - 优先级较低
				for _, specialConnect := range specialConnects {
					// 确保ServerURL不为空
					if specialConnect.ServerURL == "" {
						continue
					}

					// 检查重复并应用优先级规则
					seenMutex.Lock()
					entry, exists := seenURLs[specialConnect.ServerURL]
					if exists && entry.priority <= 2 {
						// 已存在且优先级相同或更高，跳过
						seenMutex.Unlock()
						continue
					}
					// 添加或更新，特殊connect优先级为2
					seenURLs[specialConnect.ServerURL] = ConnectEntry{
						connect:  specialConnect,
						priority: 2,
					}
					seenMutex.Unlock()

					// 发送到通道 - 后续会根据map筛选
					connectChan <- specialConnect
				}
				return
			}

			// 非特殊connect处理 - 优先级较高
			if connectInfo.Code != 1 || connectInfo.Data.ServerURL == "" {
				return
			}

			// 检查重复并应用优先级规则
			seenMutex.Lock()
			entry, exists := seenURLs[connectInfo.Data.ServerURL]
			if exists && entry.priority <= 1 {
				// 已存在且优先级相同或更高，跳过
				seenMutex.Unlock()
				return
			}
			// 添加或更新，非特殊connect优先级为1（更高）
			seenURLs[connectInfo.Data.ServerURL] = ConnectEntry{
				connect:  connectInfo.Data,
				priority: 1,
			}
			seenMutex.Unlock()

			// 发送到通道
			connectChan <- connectInfo.Data
		}(conn)
	}

	go func() {
		wg.Wait()
		close(connectChan)
	}()

	// 创建最终的去重结果集
	finalResults := make(map[string]models.Connect)

	// 从通道收集所有connect
	for connect := range connectChan {
		// 再次检查ServerURL是否为空
		if connect.ServerURL == "" {
			continue
		}

		// 根据优先级处理
		seenMutex.Lock()
		entry, exists := seenURLs[connect.ServerURL]
		// 只添加优先级最高的项到最终结果
		if !exists || entry.priority == 1 && entry.connect.ServerURL == connect.ServerURL {
			finalResults[connect.ServerURL] = connect
		}
		seenMutex.Unlock()
	}

	// 转换map为切片
	for _, connect := range finalResults {
		connectList = append(connectList, connect)
	}

	return connectList, nil
}
