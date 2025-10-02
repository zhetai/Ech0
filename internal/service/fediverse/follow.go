package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	model "github.com/lin-snow/ech0/internal/model/fediverse"
	userModel "github.com/lin-snow/ech0/internal/model/user"
	httpUtil "github.com/lin-snow/ech0/internal/util/http"
)

func (fediverseService *FediverseService) handleFollowActivity(user *userModel.User, activity *model.Activity) error {
	fmt.Println("Handling follow activity:", activity)

	followerActor := activity.ActorURL
	if followerActor == "" {
		followerActor = activity.ActorID
	}
	if followerActor == "" {
		return errors.New("follow activity missing actor")
	}

	actor, setting, err := fediverseService.BuildActor(user)
	if err != nil {
		return err
	}
	serverURL := httpUtil.TrimURL(setting.ServerURL)

	acceptPayload, err := fediverseService.buildAcceptActivityPayload(&actor, activity, followerActor, serverURL)
	if err != nil {
		fmt.Printf("Error building accept activity payload: %v\n", err)
		return err
	}

	inboxURL, err := fediverseService.fetchRemoteActorInbox(followerActor)
	if err != nil {
		fmt.Printf("Error fetching follower inbox: %v\n", err)
		return err
	}

	if err := httpUtil.PostActivity(acceptPayload, inboxURL, actor.ID); err != nil {
		fmt.Printf("Error posting accept activity: %v\n", err)
		return err
	}

	// 保存粉丝记录
	fmt.Printf("Saving follower: userID=%d, actor=%s\n", user.ID, followerActor)

	return fediverseService.txManager.Run(func(ctx context.Context) error {
		return fediverseService.fediverseRepository.SaveFollower(ctx, &model.Follower{
			UserID:  user.ID,
			ActorID: followerActor,
		})
	})
}

func (fediverseService *FediverseService) buildAcceptActivityPayload(actor *model.Actor, follow *model.Activity, followerActor, serverURL string) ([]byte, error) {
	if follow.ActivityID == "" {
		return nil, errors.New("follow activity missing id")
	}

	target := follow.ObjectID
	if target == "" {
		target = actor.ID
	}

	now := time.Now().UTC()
	acceptID := fmt.Sprintf("%s/activities/%s/accept/%d", serverURL, actor.PreferredUsername, now.UnixNano())

	payload := map[string]any{
		"@context": []any{"https://www.w3.org/ns/activitystreams"},
		"id":       acceptID,
		"type":     model.ActivityTypeAccept,
		"actor":    actor.ID,
		"object": map[string]any{
			"id":     follow.ActivityID,
			"type":   model.ActivityTypeFollow,
			"actor":  followerActor,
			"object": target,
		},
		"to":        []string{followerActor},
		"published": now.Format(time.RFC3339),
	}

	return json.Marshal(payload)
}
