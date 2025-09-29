package service

import (
	"errors"

	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/fediverse"
	settingModel "github.com/lin-snow/ech0/internal/model/setting"
	repository "github.com/lin-snow/ech0/internal/repository/fediverse"
	userRepository "github.com/lin-snow/ech0/internal/repository/user"
	settingService "github.com/lin-snow/ech0/internal/service/setting"
)

type FediverseService struct {
	fediverseRepository repository.FediverseRepositoryInterface
	userRepository      userRepository.UserRepositoryInterface
	settingService      settingService.SettingServiceInterface
}

func NewFediverseService(fediverseRepository repository.FediverseRepositoryInterface, 
	userRepository userRepository.UserRepositoryInterface,
	settingService settingService.SettingServiceInterface,
	) FediverseServiceInterface {
	return &FediverseService{
		fediverseRepository: fediverseRepository,
		userRepository:      userRepository,
		settingService:      settingService,
	}
}

// GetActorByUsername 通过用户名获取 Actor 信息
func (fediverseService *FediverseService) GetActorByUsername(username string) (model.Actor, error) {
	// 查询用户
	user, err := fediverseService.userRepository.GetUserByUsername(username)
	if err != nil {
		return model.Actor{}, errors.New(commonModel.USER_NOTFOUND)
	}

	// 构建 Actor 对象
	// 从设置服务获取服务器域名
	var setting settingModel.SystemSetting
	if err := fediverseService.settingService.GetSetting(&setting); err != nil {
		return model.Actor{}, errors.New(commonModel.GET_ACTOR_ERROR)
	}
	serverURL := setting.ServerURL
	if serverURL == "" {
		return model.Actor{}, errors.New(commonModel.ACTIVEPUB_NOT_ENABLED)
	}

	actor := model.Actor{
		Context: 		 "https://www.w3.org/ns/activitystreams", // 固定值
		ID: serverURL + "/users/" + user.Username,  // 实例地址拼接 域名 + /users/ + username
		Type:               "Person", // 固定值
		Name: user.DisplayName, // 显示名称
		PreferredUsername: user.Username, // 用户名
		Summary: "这是" + user.DisplayName + "的 ActivityPub 个人资料。", // 简短介绍
		Inbox: serverURL + "/users/" + user.Username + "/inbox", // 收件箱地址
		Outbox: serverURL + "/users/" + user.Username + "/outbox", // 发件箱地址
		PublicKey: model.PublicKey{
			ID: 		 serverURL + "/users/" + user.Username + "#main-key", // 公钥ID
			Owner: 	 serverURL + "/users/" + user.Username, // 所有者
			PublicKeyPem: user.PublicKeyPEM, // 从用户数据获取公钥
		},
	}

	return actor, nil
}