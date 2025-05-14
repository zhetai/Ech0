package repository

import (
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/models"
)

func GetAllConnects() ([]models.Connected, error) {
	var connects []models.Connected
	// 查询数据库
	if err := database.DB.Find(&connects).Error; err != nil {
		return nil, err
	}
	// 如果没有找到，返回空切片
	if len(connects) == 0 {
		return []models.Connected{}, nil
	}
	// 返回查询到的 connects
	return connects, nil
}

func CreateConnect(connect *models.Connected) error {
	if err := database.DB.Create(connect).Error; err != nil {
		return err
	}
	return nil
}

func DeleteConnect(id uint) error {
	// 根据 ID 删除 Connect
	if err := database.DB.Delete(&models.Connected{}, id).Error; err != nil {
		return err
	}

	return nil
}
