package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lin-snow/ech0/internal/config"
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
	"github.com/lin-snow/ech0/internal/transaction"
	fileUtil "github.com/lin-snow/ech0/internal/util/file"
	httpUtil "github.com/lin-snow/ech0/internal/util/http"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type FediverseService struct {
	txManager           transaction.TransactionManager
	fediverseRepository repository.FediverseRepositoryInterface
	userRepository      userRepository.UserRepositoryInterface
	settingService      settingService.SettingServiceInterface
	echoService         echoService.EchoServiceInterface
}

func NewFediverseService(
	txManager transaction.TransactionManager,
	fediverseRepository repository.FediverseRepositoryInterface,
	userRepository userRepository.UserRepositoryInterface,
	settingService settingService.SettingServiceInterface,
	echoService echoService.EchoServiceInterface,
) FediverseServiceInterface {
	return &FediverseService{
		txManager:           txManager,
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
		Subject: "acct:" + user.Username + "@" + httpUtil.ExtractDomain(httpUtil.TrimURL(setting.ServerURL)),
		Aliases: []string{
			actor.ID,
			"acct:" + user.Username + "@" + httpUtil.ExtractDomain(httpUtil.TrimURL(setting.ServerURL)),
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
	user, err := fediverseService.userRepository.GetUserByUsername(username)
	if err != nil {
		return errors.New(commonModel.USER_NOTFOUND)
	}

	// 处理不同类型的 Activity
	switch activity.Type {
	case model.ActivityTypeFollow:
		// 处理关注请求
		fediverseService.txManager.Run(func(ctx context.Context) error {
			return fediverseService.fediverseRepository.SaveFollower(ctx, &model.Follower{
				UserID:  user.ID,
				ActorID: activity.ActorID,
			})
		})

	default:
		return errors.New("Unsupported activity type: " + cases.Title(language.English).String(activity.Type))
	}

	return nil
}

// HandleOutbox 处理 Outbox 消息
func (fediverseService *FediverseService) HandleOutboxPage(ctx context.Context, username string, page, pageSize int) (model.OutboxPage, error) {
	// 查询用户，确保用户存在
	user, err := fediverseService.userRepository.GetUserByUsername(username)
	if err != nil {
		return model.OutboxPage{}, errors.New(commonModel.USER_NOTFOUND)
	}

	// 获取 Actor和 setting
	actor, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return model.OutboxPage{}, err
	}
	serverURL := httpUtil.ExtractDomain(httpUtil.TrimURL(setting.ServerURL))

	// 查 Echos
	echosByPage, err := fediverseService.echoService.GetEchosByPage(authModel.NO_USER_LOGINED, commonModel.PageQueryDto{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return model.OutboxPage{}, err
	}

	// 转 Avtivity
	var activities []model.Activity
	for i := range echosByPage.Items {
		activities = append(activities, fediverseService.ConvertEchoToActivity(&echosByPage.Items[i], &actor, serverURL))
	}

	// 拼装 OutboxPage
	outboxPage := model.OutboxPage{
		ID:           fmt.Sprintf("%s/users/%s/outbox?page=%d", serverURL, username, page),
		Type:         "OrderedCollectionPage",
		PartOf:       fmt.Sprintf("%s/users/%s/outbox", serverURL, username),
		Next:         "",
		Prev:         "",
		OrderedItems: activities,
	}

	// 计算 Next && Prev
	if page > 1 {
		outboxPage.Prev = fmt.Sprintf("%s/users/%s/outbox?page=%d", serverURL, username, page-1)
	}
	if (page * pageSize) < int(echosByPage.Total) {
		outboxPage.Next = fmt.Sprintf("%s/users/%s/outbox?page=%d", serverURL, username, page+1)
	}

	return outboxPage, nil
}

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
	avatarURL := serverURL + "/api" + user.Avatar
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

// ConvertEchoToActivity 将 Echo 转换为 ActivityPub Activity
func (fediverseService *FediverseService) ConvertEchoToActivity(echo *echoModel.Echo, actor *model.Actor, serverURL string) model.Activity {
	obj := fediverseService.ConvertEchoToObject(echo, actor, serverURL)

	activityID := fmt.Sprintf("%s/activities/%d", serverURL, echo.ID)

	activity := model.Activity{
		ActivityID: activityID,
		Type:       model.ActivityTypeCreate,
		ActorID:    actor.ID,
		ActorURL:   actor.ID,
		ObjectID:   obj.ObjectID,
		ObjectType: obj.Type,
		Published:  echo.CreatedAt,
		To:         obj.To,
		Cc:         []string{},
		Summary:    "",
		Delivered:  false,
		CreatedAt:  time.Now(),
	}

	activityJSON, _ := json.Marshal(activity)
	activity.ActivityJSON = string(activityJSON)
	return activity
}

// ConvertEchoToObject 将 Echo 转换为 ActivityPub Object
func (fediverseService *FediverseService) ConvertEchoToObject(echo *echoModel.Echo, actor *model.Actor, serverURL string) model.Object {
	var attachments []model.Attachment
	for i := range echo.Images {
		attachments = append(attachments, model.Attachment{
			Type:      "Image",
			MediaType: httpUtil.GetMIMETypeFromFilenameOrURL(echo.Images[i].ImageURL),
			URL:       fileUtil.GetImageURL(echo.Images[i], serverURL),
		})
	}

	return model.Object{
		ObjectID:     fmt.Sprintf("%s/objects/%d", serverURL, echo.ID),
		Type:         "Note",
		Content:      echo.Content,
		AttributedTo: fmt.Sprintf("%s/users/%s", serverURL, echo.Username),
		Published:    echo.CreatedAt,
		To: []string{
			"https://www.w3.org/ns/activitystreams#Public",
		},
		Attachments: attachments,
	}
}

// GetFollowers 获取粉丝列表
func (fediverseService *FediverseService) GetFollowers(username string) (model.FollowersResponse, error) {
	user, err := fediverseService.userRepository.GetUserByUsername(username)
	if err != nil {
		return model.FollowersResponse{}, errors.New(commonModel.USER_NOTFOUND)
	}

	followers, err := fediverseService.fediverseRepository.GetFollowers(user.ID)
	if err != nil {
		return model.FollowersResponse{}, err
	}

	// 构建Actor URL 列表
	var followerURLs []string
	for _, follower := range followers {
		followerURLs = append(followerURLs, follower.ActorID)
	}

	return model.FollowersResponse{
		Context:      "https://www.w3.org/ns/activitystreams",
		ID:           "",
		Type:         "OrderedCollection",
		TotalItems:   len(followerURLs),
		First:        "",
		OrderedItems: followerURLs,
	}, nil
}

// GetFollowing 获取关注列表
func (fediverseService *FediverseService) GetFollowing(username string) (model.FollowingResponse, error) {
	user, err := fediverseService.userRepository.GetUserByUsername(username)
	if err != nil {
		return model.FollowingResponse{}, errors.New(commonModel.USER_NOTFOUND)
	}

	following, err := fediverseService.fediverseRepository.GetFollowing(user.ID)
	if err != nil {
		return model.FollowingResponse{}, err
	}

	// 构建Actor URL 列表
	var followingURLs []string
	for _, follow := range following {
		followingURLs = append(followingURLs, follow.ObjectID)
	}

	return model.FollowingResponse{
		Context:      "https://www.w3.org/ns/activitystreams",
		ID:           "",
		Type:         "OrderedCollection",
		TotalItems:   len(followingURLs),
		First:        "",
		OrderedItems: followingURLs,
	}, nil
}

// GetObjectByID 通过 ID 获取内容对象
func (fediverseService *FediverseService) GetObjectByID(id uint) (model.Object, error) {
	// 获取 Echo
	echo, err := fediverseService.echoService.GetEchoById(authModel.NO_USER_LOGINED, id)
	if err != nil {
		return model.Object{}, err
	}

	// 获取 Actor 和 setting
	user, err := fediverseService.userRepository.GetUserByUsername(echo.Username)
	if err != nil {
		return model.Object{}, err
	}
	actor, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return model.Object{}, err
	}
	serverURL := httpUtil.ExtractDomain(httpUtil.TrimURL(setting.ServerURL))

	// 转 Object
	return fediverseService.ConvertEchoToObject(echo, &actor, serverURL), nil
}
