package repository

import (
	"context"
	model "github.com/lin-snow/ech0/internal/model/user"
)

type UserRepositoryInterface interface {
	// GetUserByID 根据用户ID获取用户
	GetUserByID(id int) (model.User, error)

	// GetUserByUsername 根据用户名获取用户
	GetUserByUsername(username string) (model.User, error)

	// GetAllUsers 获取所有用户
	GetAllUsers() ([]model.User, error)

	// CreateUser 创建一个新的用户
	CreateUser(ctx context.Context, newUser *model.User) error

	// GetSysAdmin 获取系统管理员
	GetSysAdmin() (model.User, error)

	// UpdateUser 更新用户
	UpdateUser(ctx context.Context, user *model.User) error

	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, id uint) error
}
