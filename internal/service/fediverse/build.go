package service

import (
	"errors"
	"fmt"

	"github.com/lin-snow/ech0/internal/config"
	authModel "github.com/lin-snow/ech0/internal/model/auth"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/fediverse"
	settingModel "github.com/lin-snow/ech0/internal/model/setting"
	userModel "github.com/lin-snow/ech0/internal/model/user"
	httpUtil "github.com/lin-snow/ech0/internal/util/http"
)

// BuildOutbox 构建 Outbox 元信息
func (fediverseService *FediverseService) BuildOutbox(username string) (model.OutboxResponse, error) {
	// 查询用户，确保用户存在
	user, err := fediverseService.userRepository.GetUserByUsername(username)
	if err != nil {
		return model.OutboxResponse{}, errors.New(commonModel.USER_NOTFOUND)
	}

	// 获取 Actor和 setting
	_, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return model.OutboxResponse{}, err
	}
	serverURL := httpUtil.ExtractDomain(httpUtil.TrimURL(setting.ServerURL))

	// 查 Echos
	echosByPage, err := fediverseService.echoService.GetEchosByPage(authModel.NO_USER_LOGINED, commonModel.PageQueryDto{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		return model.OutboxResponse{}, err
	}

	// 拼装 OutboxResponse
	outbox := model.OutboxResponse{
		Context:    "https://www.w3.org/ns/activitystreams",
		ID:         fmt.Sprintf("%s/users/%s/outbox", serverURL, username),
		Type:       "OrderedCollection",
		TotalItems: int(echosByPage.Total),
		First:      fmt.Sprintf("%s/users/%s/outbox?page=1", serverURL, username),
		Last:       "", // 这里暂时设为空，实际应用中应计算最后一页的链接
	}

	return outbox, nil
}

// BuildActor 构建 Actor 对象
func (fediverseService *FediverseService) BuildActor(user *userModel.User) (model.Actor, *settingModel.SystemSetting, error) {
	// 从设置服务获取服务器域名
	var setting settingModel.SystemSetting
	if err := fediverseService.settingService.GetSetting(&setting); err != nil {
		return model.Actor{}, nil, err
	}
	serverURL := httpUtil.TrimURL(setting.ServerURL)
	if serverURL == "" {
		return model.Actor{}, nil, errors.New(commonModel.ACTIVEPUB_NOT_ENABLED)
	}
	// 构建头像信息 (域名 + /api + 头像路径)
	if user.Avatar == "" {
		user.Avatar = "/Ech0.png" // 默认头像路径
	} else {
		user.Avatar = "/api" + user.Avatar
	}
	avatarURL := serverURL + user.Avatar
	avatarMIME := httpUtil.GetMIMETypeFromFilenameOrURL(avatarURL)

	// 构建 Actor 对象
	return model.Actor{
		Context: []any{
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
		},
		ID:                serverURL + "/users/" + user.Username, // 实例地址拼接 域名 + /users/ + username
		Type:              "Person",                              // 固定值
		Name:              setting.ServerName,                         // 显示名称
		PreferredUsername: user.Username, 					   // 用户名
		Summary:           "你好呀!👋 我是来自Ech0的" + user.Username,     // 简介
		Icon: model.Preview{
			Type: 	"Image",
			MediaType: avatarMIME,
			URL: avatarURL,
		},
		Image: model.Preview{
			Type: "Image",
			MediaType: "image/png",
			URL: serverURL + "/banner.png", // 封面图片，固定为 /banner.png
		},
		Followers: 	 serverURL + "/users/" + user.Username + "/followers", // 粉丝列表地址
		Following: 	 serverURL + "/users/" + user.Username + "/following", // 关注列表地址
		Inbox:             serverURL + "/users/" + user.Username + "/inbox", // 收件箱地址
		Outbox:            serverURL + "/users/" + user.Username + "/outbox", // 发件箱地址
		PublicKey: model.PublicKey{
			ID:           serverURL + "/users/" + user.Username + "#main-key",
			Owner:        serverURL + "/users/" + user.Username,
			PublicKeyPem: string(config.RSA_PUBLIC_KEY),
		},
	}, &setting, nil
}
