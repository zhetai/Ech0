package service

import (
	model "github.com/lin-snow/ech0/internal/model/fediverse"
	settingModel "github.com/lin-snow/ech0/internal/model/setting"
	userModel "github.com/lin-snow/ech0/internal/model/user"
)

type FediverseServiceInterface interface{
	// BuildActor 构建 Actor 对象
	BuildActor(user *userModel.User) (model.Actor, *settingModel.SystemSetting, error)

	// WebFinger 处理 Webfinger 请求
	Webfinger(username string) (model.WebFingerResponse, error)

	// GetActorByUsername 通过用户名获取 Actor 信息
	GetActorByUsername(username string) (model.Actor, error)

	// HandleInbox 处理接收到的 ActivityPub 消息
	HandleInbox(username string, activity *model.Activity) error

	// HandleOutbox 处理 Outbox 消息
	HandleOutbox(username string) (model.OutboxResponse, error)
}