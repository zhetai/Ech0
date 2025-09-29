package service

import (
	"context"
	"errors"
	"fmt"

	authModel "github.com/lin-snow/ech0/internal/model/auth"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	echoModel "github.com/lin-snow/ech0/internal/model/echo"
	model "github.com/lin-snow/ech0/internal/model/fediverse"
	settingModel "github.com/lin-snow/ech0/internal/model/setting"
	userModel "github.com/lin-snow/ech0/internal/model/user"
	repository "github.com/lin-snow/ech0/internal/repository/fediverse"
	userRepository "github.com/lin-snow/ech0/internal/repository/user"
	echoService "github.com/lin-snow/ech0/internal/service/echo"
	settingService "github.com/lin-snow/ech0/internal/service/setting"
	httpUtil "github.com/lin-snow/ech0/internal/util/http"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type FediverseService struct {
	fediverseRepository repository.FediverseRepositoryInterface
	userRepository      userRepository.UserRepositoryInterface
	settingService      settingService.SettingServiceInterface
	echoService         echoService.EchoServiceInterface
}

func NewFediverseService(fediverseRepository repository.FediverseRepositoryInterface,
	userRepository userRepository.UserRepositoryInterface,
	settingService settingService.SettingServiceInterface,
	echoService echoService.EchoServiceInterface,
) FediverseServiceInterface {
	return &FediverseService{
		fediverseRepository: fediverseRepository,
		userRepository:      userRepository,
		settingService:      settingService,
		echoService:         echoService,
	}
}

// Webfinger 处理 Webfinger 请求
func (fediverseService *FediverseService) Webfinger(username string) (model.WebFingerResponse, error) {
	// 查询用户
	user, err := fediverseService.userRepository.GetUserByUsername(username)
	if err != nil {
		return model.WebFingerResponse{}, errors.New(commonModel.USER_NOTFOUND)
	}

	// 构建 Actor 对象
	actor, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return model.WebFingerResponse{}, err
	}

	return model.WebFingerResponse{
		Subject: "acct:" + user.Username + "@" + httpUtil.TrimURL(setting.ServerURL),
		Aliases: []string{
			actor.ID,
		},
		Links: []model.Link{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: actor.ID,
			},
		},
	}, nil
}

// GetActorByUsername 通过用户名获取 Actor 信息
func (fediverseService *FediverseService) GetActorByUsername(username string) (model.Actor, error) {
	// 查询用户
	user, err := fediverseService.userRepository.GetUserByUsername(username)
	if err != nil {
		return model.Actor{}, errors.New(commonModel.USER_NOTFOUND)
	}

	// 构建 Actor 对象
	actor, _, err := fediverseService.BuildActor(&user)
	if err != nil {
		return model.Actor{}, err
	}

	return actor, nil
}

// ProcessInbox 处理接收到的 ActivityPub 消息
func (fediverseService *FediverseService) HandleInbox(username string, activity *model.Activity) error {
	// 查询用户，确保用户存在
	// user, err := fediverseService.userRepository.GetUserByUsername(username)
	// if err != nil {
	// 	return errors.New(commonModel.USER_NOTFOUND)
	// }

	// 处理不同类型的 Activity
	switch activity.Type {
	case model.ActivityTypeFollow:
		// 处理关注请求

	case model.ActivityTypeAccept:
		// 处理关注接受

	case model.ActivityTypeCreate:
		// 处理创建内容

	case model.ActivityTypeLike:
		// 处理点赞

	case model.ActivityTypeAnnounce:
		// 处理转发

	case model.ActivityTypeUndo:
		// 处理撤销

	default:
		return errors.New("Unsupported activity type: " + cases.Title(language.English).String(activity.Type))
	}

	return nil
}

// HandleOutbox 处理 Outbox 消息
func (fediverseService *FediverseService) HandleOutbox(ctx context.Context, username string, page, pageSize int) (model.OutboxResponse, error) {
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
	serverURL := httpUtil.TrimURL(setting.ServerURL)

	// 查 Echos
	echosByPage, err := fediverseService.echoService.GetEchosByPage(authModel.NO_USER_LOGINED, commonModel.PageQueryDto{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return model.OutboxResponse{}, err
	}

	// 转 Avtivity
	var activities []model.Activity
	for i := range echosByPage.Items {
		activities = append(activities, fediverseService.ConvertEchoToActivity(&echosByPage.Items[i], nil))
	}

	// 拼装 OutboxPage
	outboxPage := model.OutboxPage{
		ID:           fmt.Sprintf("%s/users/%s/outbox?page=%d", serverURL, username, page),
		Type:         "OrderedCollectionPage",
		PartOf:       fmt.Sprintf("%s/users/%s/outbox", serverURL, username),
		Next:         fmt.Sprintf("%s/users/%s/outbox?page=%d", serverURL, username, page+1),
		Prev:         "",
		OrderedItems: activities,
	}

	return model.OutboxResponse{
		Context:      []any{},
		ID:           outboxPage.ID,
		Type:         "OrderedCollection",
		TotalItems:   int(echosByPage.Total),
		First:        &outboxPage,
		Last:         nil,
		OrderedItems: nil,
	}, nil
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

	return model.Actor{
		Context: []any{
			"https://www.w3.org/ns/activitystreams",
		},
		ID:                serverURL + "/users/" + user.Username, // 实例地址拼接 域名 + /users/ + username
		Type:              "Person",                              // 固定值
		Name:              user.DisplayName,                      // 显示名称
		PreferredUsername: user.Username,                         // 用户名
		Summary:           "这是" + user.DisplayName + "的 ActivityPub 个人资料。",
		Inbox:             serverURL + "/users/" + user.Username + "/inbox",
		Outbox:            serverURL + "/users/" + user.Username + "/outbox",
		PublicKey: model.PublicKey{
			ID:           serverURL + "/users/" + user.Username + "#main-key",
			Owner:        serverURL + "/users/" + user.Username,
			PublicKeyPem: user.PublicKeyPEM,
		},
	}, &setting, nil
}

func (fediverseService *FediverseService) ConvertEchoToActivity(echo *echoModel.Echo, actor *model.Actor) model.Activity {
	return model.Activity{}
}
