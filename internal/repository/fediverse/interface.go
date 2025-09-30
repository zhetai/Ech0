package repository

import (
	"context"

	model "github.com/lin-snow/ech0/internal/model/fediverse"
)

type FediverseRepositoryInterface interface {
	// 通过用户ID获取粉丝列表
	GetFollowers(userID uint) ([]model.Follower, error)

	// 通过用户ID获取关注列表
	GetFollowing(userID uint) ([]model.Follow, error)

	// 存储新的粉丝
	SaveFollower(ctx context.Context, follower *model.Follower) error
}
