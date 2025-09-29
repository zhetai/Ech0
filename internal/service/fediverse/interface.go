package service

import (
	model "github.com/lin-snow/ech0/internal/model/fediverse"
)

type FediverseServiceInterface interface{
	// GetActorByUsername 通过用户名获取 Actor 信息
	GetActorByUsername(username string) (model.Actor, error)

	// ProcessInbox 处理接收到的 ActivityPub 消息
	HandleInbox(username string, activity *model.Activity) error
}