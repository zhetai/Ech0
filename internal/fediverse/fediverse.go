package fediverse

import (
	echoRepository "github.com/lin-snow/ech0/internal/repository/echo"
	repository "github.com/lin-snow/ech0/internal/repository/fediverse"
	userRepository "github.com/lin-snow/ech0/internal/repository/user"
	settingService "github.com/lin-snow/ech0/internal/service/setting"
)

// FediverseCore 联邦宇宙核心
type FediverseCore struct {
	repo           repository.FediverseRepositoryInterface
	settingService settingService.SettingServiceInterface
	userRepository userRepository.UserRepositoryInterface
	echoRepository echoRepository.EchoRepositoryInterface
}

func NewFediverseCore(
	repo repository.FediverseRepositoryInterface,
	settingService settingService.SettingServiceInterface,
	userRepository userRepository.UserRepositoryInterface,
	echoRepository echoRepository.EchoRepositoryInterface,
) *FediverseCore {
	return &FediverseCore{
		repo:           repo,
		settingService: settingService,
		userRepository: userRepository,
		echoRepository: echoRepository,
	}
}
