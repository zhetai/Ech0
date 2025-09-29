package service

import (
	model "github.com/lin-snow/ech0/internal/model/fediverse"
)

type FediverseServiceInterface interface{
	// GetActorByUsername 通过用户名获取 Actor 信息
	GetActorByUsername(username string) (model.Actor, error)
}