package repository

import (
	"context"

	model "github.com/lin-snow/ech0/internal/model/echo"
)

type EchoRepositoryInterface interface {
	// CreateEcho 创建一个新的 Echo
	CreateEcho(ctx context.Context, echo *model.Echo) error

	// GetEchosByPage 获取分页的 Echo 列表
	GetEchosByPage(page, pageSize int, search string, showPrivate bool) ([]model.Echo, int64)

	// GetEchosById 根据 ID 获取 Echo
	GetEchosById(id uint) (*model.Echo, error)

	// DeleteEchoById 删除 Echo
	DeleteEchoById(ctx context.Context, id uint) error

	// GetTodayEchos 获取今天的 Echo 列表
	GetTodayEchos(showPrivate bool) []model.Echo

	// UpdateEcho 更新 Echo
	UpdateEcho(ctx context.Context, echo *model.Echo) error

	// LikeEcho 点赞 Echo
	LikeEcho(ctx context.Context, id uint) error

	// GetAllTags 获取所有标签
	GetAllTags() ([]model.Tag, error)

	// DeleteTagById 删除标签
	DeleteTagById(ctx context.Context, id uint) error
}
