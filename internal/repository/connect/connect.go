package repository

import (
	model "github.com/lin-snow/ech0/internal/model/connect"
	"gorm.io/gorm"
)

type ConnectRepository struct {
	db *gorm.DB
}

func NewConnectRepository(db *gorm.DB) ConnectRepositoryInterface {
	return &ConnectRepository{
		db: db,
	}
}

// GetAllConnects 获取所有连接
func (connectRepository *ConnectRepository) GetAllConnects() ([]model.Connected, error) {
	var connects []model.Connected
	// 查询数据库
	if err := connectRepository.db.Find(&connects).Error; err != nil {
		return nil, err
	}
	// 如果没有找到，返回空切片
	if len(connects) == 0 {
		return []model.Connected{}, nil
	}
	// 返回查询到的 connects
	return connects, nil
}

// CreateConnect 创建一个新的连接
func (connectRepository *ConnectRepository) CreateConnect(connect *model.Connected) error {
	if err := connectRepository.db.Create(connect).Error; err != nil {
		return err
	}
	return nil
}

// DeleteConnect 删除连接
func (connectRepository *ConnectRepository) DeleteConnect(id uint) error {
	// 根据 ID 删除 Connect
	if err := connectRepository.db.Delete(&model.Connected{}, id).Error; err != nil {
		return err
	}

	return nil
}
