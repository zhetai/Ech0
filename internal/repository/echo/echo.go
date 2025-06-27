package repository

import (
	"errors"
	"strings"
	"time"

	model "github.com/lin-snow/ech0/internal/model/echo"
	"gorm.io/gorm"
)

type EchoRepository struct {
	db *gorm.DB
}

func NewEchoRepository(db *gorm.DB) EchoRepositoryInterface {
	return &EchoRepository{db: db}
}

func (echoRepository *EchoRepository) CreateEcho(echo *model.Echo) error {
	echo.Content = strings.TrimSpace(echo.Content)

	result := echoRepository.db.Create(echo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (echoRepository *EchoRepository) GetEchosByPage(page, pageSize int, search string, showPrivate bool) ([]model.Echo, int64) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询数据库
	var echos []model.Echo
	var total int64

	query := echoRepository.db.Model(&model.Echo{})

	// 如果 search 不为空，添加模糊查询条件
	if search != "" {
		searchPattern := "%" + search + "%" // 模糊匹配模式
		query = query.Where("content LIKE ?", searchPattern)
	}

	// 如果不是管理员，过滤私密留言
	if !showPrivate {
		query = query.Where("private = ?", false)
	}

	// 获取总数并进行分页查询
	query.Count(&total).
		Preload("Images").
		Limit(pageSize).
		Offset(offset).
		Order("created_at DESC").
		Find(&echos)

	// 返回结果
	return echos, total
}

func (echoRepository *EchoRepository) GetEchosById(id uint) (*model.Echo, error) {
	var echo model.Echo
	result := echoRepository.db.Preload("Images").First(&echo, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 如果未找到记录，则返回 nil
		}
		return nil, result.Error // 其他错误返回
	}

	return &echo, nil
}

func (echoRepository *EchoRepository) DeleteEchoById(id uint) error {
	var echo model.Echo
	// 删除外键images
	echoRepository.db.Where("message_id = ?", id).Delete(&model.Image{})

	result := echoRepository.db.Delete(&echo, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // 如果没有找到记录
	}
	return nil
}

func (echoRepository *EchoRepository) GetTodayEchos(showPrivate bool) []model.Echo {
	// 查询数据库
	var echos []model.Echo

	// 获取当天开始和结束时间
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := echoRepository.db.Model(&model.Echo{})
	// 如果不是管理员，过滤私密留言
	if !showPrivate {
		query = query.Where("private = ?", false)
	}

	// 添加当天的时间过滤
	query = query.Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay)

	// 获取总数并进行分页查询
	query.
		Preload("Images").
		Order("created_at DESC").
		Find(&echos)

	// 返回结果
	return echos
}
