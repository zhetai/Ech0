package repository

import (
	"context"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	echoModel "github.com/lin-snow/ech0/internal/model/echo"
	userModel "github.com/lin-snow/ech0/internal/model/user"
	"github.com/lin-snow/ech0/internal/transaction"
	"gorm.io/gorm"
)

type CommonRepository struct {
	db *gorm.DB
}

func NewCommonRepository(db *gorm.DB) CommonRepositoryInterface {
	return &CommonRepository{
		db: db,
	}
}

// getDB 从上下文中获取事务
func (commonRepository *CommonRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transaction.TxKey).(*gorm.DB); ok {
		return tx
	}
	return commonRepository.db
}

// GetUserByUserId 根据用户ID获取用户信息
func (commonRepository *CommonRepository) GetUserByUserId(userId uint) (userModel.User, error) {
	var user userModel.User
	if err := commonRepository.db.First(&user, userId).Error; err != nil {
		return user, err
	}
	return user, nil
}

// GetSysAdmin 获取系统管理员信息
func (commonRepository *CommonRepository) GetSysAdmin() (userModel.User, error) {
	// 获取系统管理员（首个注册的用户）
	user := userModel.User{}
	err := commonRepository.db.Where("is_admin = ?", true).First(&user).Error
	if err != nil {
		return userModel.User{}, err
	}

	return user, nil
}

// GetAllUsers 获取所有用户信息
func (commonRepository *CommonRepository) GetAllUsers() ([]userModel.User, error) {
	var users []userModel.User
	err := commonRepository.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetAllEchos 获取所有Echo
func (commonRepository *CommonRepository) GetAllEchos(showPrivate bool) ([]echoModel.Echo, error) {
	var echos []echoModel.Echo

	// 是否将私密内容也查询出来
	if showPrivate {
		if err := commonRepository.db.Preload("Images").Order("created_at DESC").Find(&echos).Error; err != nil {
			return nil, err
		}
	} else {
		if err := commonRepository.db.Preload("Images").Where("private = ?", false).Find(&echos).Error; err != nil {
			return nil, err
		}
	}

	return echos, nil
}

// GetHeatMap 获取热力图数据
func (commonRepository *CommonRepository) GetHeatMap(startDate, endDate string) ([]commonModel.Heatmap, error) {
	var results []commonModel.Heatmap

	// 查询数据
	// 执行查询
	err := commonRepository.db.Table("echos").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", startDate, endDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
