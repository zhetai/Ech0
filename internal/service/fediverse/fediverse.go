package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	commonModel "github.com/lin-snow/ech0/internal/model/common"
	echoModel "github.com/lin-snow/ech0/internal/model/echo"
	model "github.com/lin-snow/ech0/internal/model/fediverse"
	echoRepository "github.com/lin-snow/ech0/internal/repository/echo"
	repository "github.com/lin-snow/ech0/internal/repository/fediverse"
	userRepository "github.com/lin-snow/ech0/internal/repository/user"
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
	echoRepository     echoRepository.EchoRepositoryInterface
}

func NewFediverseService(
	txManager transaction.TransactionManager,
	fediverseRepository repository.FediverseRepositoryInterface,
	userRepository userRepository.UserRepositoryInterface,
	settingService settingService.SettingServiceInterface,
	echoRepository echoRepository.EchoRepositoryInterface,
) FediverseServiceInterface {
	return &FediverseService{
		txManager:           txManager,
		fediverseRepository: fediverseRepository,
		userRepository:      userRepository,
		settingService:      settingService,
		echoRepository:      echoRepository,
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
		if err := fediverseService.handleFollowActivity(&user, activity); err != nil {
			return err
		}

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
	serverURL, err := normalizeServerURL(setting.ServerURL)
	if err != nil {
		return model.OutboxPage{}, err
	}

	// 查 Echos
	echosByPage, total := fediverseService.echoRepository.GetEchosByPage(page, pageSize, "", false)
	if err != nil {
		return model.OutboxPage{}, err
	}

	// 转 Avtivity
	var activities []model.Activity
	for i := range echosByPage {
		activities = append(activities, fediverseService.ConvertEchoToActivity(&echosByPage[i], &actor, serverURL))
	}

	// 拼装 OutboxPage
	outboxPage := model.OutboxPage{
		Context:      "https://www.w3.org/ns/activitystreams",
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
	if (page * pageSize) < int(total) {
		outboxPage.Next = fmt.Sprintf("%s/users/%s/outbox?page=%d", serverURL, username, page+1)
	}

	return outboxPage, nil
}

// ConvertEchoToActivity 将 Echo 转换为 ActivityPub Activity
func (fediverseService *FediverseService) ConvertEchoToActivity(echo *echoModel.Echo, actor *model.Actor, serverURL string) model.Activity {
	obj := fediverseService.ConvertEchoToObject(echo, actor, serverURL)

	activityID := fmt.Sprintf("%s/activities/%d", serverURL, echo.ID)

	activity := model.Activity{
		Context:    "https://www.w3.org/ns/activitystreams",
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
		Context:      "https://www.w3.org/ns/activitystreams",
		ObjectID:     fmt.Sprintf("%s/objects/%d", serverURL, echo.ID),
		Type:         "Note",
		Content:      echo.Content,
		AttributedTo: actor.ID,
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

	actor, _, err := fediverseService.BuildActor(&user)
	if err != nil {
		return model.FollowersResponse{}, err
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
		ID:           actor.Followers,
		Type:         "OrderedCollection",
		TotalItems:   len(followerURLs),
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
	echo, err := fediverseService.echoRepository.GetEchosById(id)
	if err != nil || echo.Private {
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
	serverURL, err := normalizeServerURL(setting.ServerURL)
	if err != nil {
		return model.Object{}, err
	}

	// 转 Object
	return fediverseService.ConvertEchoToObject(echo, &actor, serverURL), nil
}

func (fediverseService *FediverseService) fetchRemoteActorInbox(actorURL string) (string, error) {
	if actorURL == "" {
		return "", errors.New("remote actor url is empty")
	}

	body, err := httpUtil.SendRequest(actorURL, http.MethodGet, httpUtil.Header{
		Header:  "Accept",
		Content: "application/activity+json",
	})
	if err != nil {
		return "", err
	}

	var resp struct {
		Inbox     string `json:"inbox"`
		Endpoints struct {
			SharedInbox string `json:"sharedInbox"`
		} `json:"endpoints"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	if resp.Inbox != "" {
		return resp.Inbox, nil
	}
	if resp.Endpoints.SharedInbox != "" {
		return resp.Endpoints.SharedInbox, nil
	}

	return "", errors.New("remote actor inbox not found")
}
