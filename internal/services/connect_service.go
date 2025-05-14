package services

import (
	"errors"
	"fmt"

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
