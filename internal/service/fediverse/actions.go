package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/fediverse"
	httpUtil "github.com/lin-snow/ech0/internal/util/http"
)

func (fediverseService *FediverseService) SearchActorByActorID(actorID string) (map[string]any, error) {
	actorID = strings.TrimSpace(actorID)
	if actorID == "" {
		return nil, errors.New(commonModel.FEDIVERSE_INVALID_INPUT)
	}

	resolvedActorURL, err := resolveActorURL(actorID)
	if err != nil {
		return nil, err
	}

	body, err := httpUtil.SendRequest(resolvedActorURL, http.MethodGet, httpUtil.Header{
		Header:  "Accept",
		Content: "application/activity+json",
	}, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", commonModel.GET_ACTOR_ERROR, err)
	}

	var actor map[string]any
	if err := json.Unmarshal(body, &actor); err != nil {
		return nil, fmt.Errorf("%s: %w", commonModel.GET_ACTOR_ERROR, err)
	}
	if len(actor) == 0 {
		return nil, errors.New(commonModel.GET_ACTOR_ERROR)
	}

	return actor, nil
}

func resolveActorURL(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New(commonModel.FEDIVERSE_INVALID_INPUT)
	}

	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return trimmed, nil
	}

	resource := trimmed
	if after, ok :=strings.CutPrefix(resource, "acct:"); ok  {
		resource = after
	}
	resource = strings.TrimPrefix(resource, "@")

	if !strings.Contains(resource, "@") {
		return "", errors.New(commonModel.GET_ACTOR_ERROR)
	}

	parts := strings.SplitN(resource, "@", 2)
	username := strings.TrimSpace(parts[0])
	domain := strings.TrimSpace(parts[1])
	if username == "" || domain == "" {
		return "", errors.New(commonModel.GET_ACTOR_ERROR)
	}

	webfingerURL := fmt.Sprintf("https://%s/.well-known/webfinger?resource=%s", domain, url.QueryEscape("acct:"+username+"@"+domain))
	body, err := httpUtil.SendRequest(webfingerURL, http.MethodGet, httpUtil.Header{
		Header:  "Accept",
		Content: "application/jrd+json, application/json",
	}, 5*time.Second)
	if err != nil {
		return "", fmt.Errorf("%s: %w", commonModel.GET_ACTOR_ERROR, err)
	}

	var resp struct {
		Links []struct {
			Rel  string `json:"rel"`
			Type string `json:"type"`
			Href string `json:"href"`
		} `json:"links"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("%s: %w", commonModel.GET_ACTOR_ERROR, err)
	}

	for _, link := range resp.Links {
		if link.Rel == "self" && link.Href != "" {
			if link.Type == "application/activity+json" || link.Type == "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" || link.Type == "" {
				return link.Href, nil
			}
		}
	}

	return "", errors.New(commonModel.GET_ACTOR_ERROR)
}

func (fediverseService *FediverseService) FollowActor(userID uint, req model.FollowActionRequest) (map[string]string, error) {
	target := strings.TrimSpace(req.TargetActor)
	if target == "" {
		return nil, errors.New(commonModel.FEDIVERSE_INVALID_INPUT)
	}

	user, err := fediverseService.commonService.CommonGetUserByUserId(userID)
	if err != nil {
		return nil, err
	}

	actor, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return nil, err
	}

	serverURL, err := normalizeServerURL(setting.ServerURL)
	if err != nil {
		return nil, err
	}

	published := time.Now().UTC()
	activityID := fmt.Sprintf("%s/activities/%s/follow/%d", serverURL, actor.PreferredUsername, published.UnixNano())

	payload, err := buildFollowActivityPayload(&actor, target, activityID, published)
	if err != nil {
		return nil, err
	}

	inboxURL, err := fediverseService.fetchRemoteActorInbox(target)
	if err != nil {
		return nil, err
	}

	if err := httpUtil.PostActivity(payload, inboxURL, actor.ID); err != nil {
		return nil, err
	}

	if err := fediverseService.txManager.Run(func(ctx context.Context) error {
		follow := &model.Follow{
			UserID:     user.ID,
			ActorID:    actor.ID,
			ObjectID:   target,
			ActivityID: activityID,
			Status:     "pending",
		}
		return fediverseService.fediverseRepository.SaveOrUpdateFollow(ctx, follow)
	}); err != nil {
		return nil, err
	}

	return map[string]string{
		"activityId": activityID,
	}, nil
}

func (fediverseService *FediverseService) UnfollowActor(userID uint, req model.FollowActionRequest) (map[string]string, error) {
	target := strings.TrimSpace(req.TargetActor)
	if target == "" {
		return nil, errors.New(commonModel.FEDIVERSE_INVALID_INPUT)
	}

	user, err := fediverseService.commonService.CommonGetUserByUserId(userID)
	if err != nil {
		return nil, err
	}

	actor, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return nil, err
	}

	serverURL, err := normalizeServerURL(setting.ServerURL)
	if err != nil {
		return nil, err
	}

	follow, err := fediverseService.fediverseRepository.GetFollowByUserAndObject(context.Background(), user.ID, target)
	if err != nil {
		return nil, err
	}
	if follow == nil || follow.ActivityID == "" {
		return nil, errors.New(commonModel.FOLLOW_RELATION_MISSING)
	}

	published := time.Now().UTC()
	undoID := fmt.Sprintf("%s/activities/%s/unfollow/%d", serverURL, actor.PreferredUsername, published.UnixNano())

	payload, err := buildUndoFollowActivityPayload(&actor, target, undoID, follow.ActivityID, published)
	if err != nil {
		return nil, err
	}

	inboxURL, err := fediverseService.fetchRemoteActorInbox(target)
	if err != nil {
		return nil, err
	}

	if err := httpUtil.PostActivity(payload, inboxURL, actor.ID); err != nil {
		return nil, err
	}

	if err := fediverseService.txManager.Run(func(ctx context.Context) error {
		return fediverseService.fediverseRepository.DeleteFollow(ctx, follow.ID)
	}); err != nil {
		return nil, err
	}

	return map[string]string{
		"activityId":       undoID,
		"followActivityId": follow.ActivityID,
	}, nil
}

func (fediverseService *FediverseService) LikeObject(userID uint, req model.LikeActionRequest) (map[string]string, error) {
	targetActor := strings.TrimSpace(req.TargetActor)
	object := strings.TrimSpace(req.Object)
	if targetActor == "" || object == "" {
		return nil, errors.New(commonModel.FEDIVERSE_INVALID_INPUT)
	}

	user, err := fediverseService.commonService.CommonGetUserByUserId(userID)
	if err != nil {
		return nil, err
	}

	actor, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return nil, err
	}

	serverURL, err := normalizeServerURL(setting.ServerURL)
	if err != nil {
		return nil, err
	}

	likeID := generateDeterministicActivityID(serverURL, actor.PreferredUsername, "like", object)
	published := time.Now().UTC()

	payload, err := buildLikeActivityPayload(&actor, targetActor, object, likeID, published)
	if err != nil {
		return nil, err
	}

	inboxURL, err := fediverseService.fetchRemoteActorInbox(targetActor)
	if err != nil {
		return nil, err
	}

	if err := httpUtil.PostActivity(payload, inboxURL, actor.ID); err != nil {
		return nil, err
	}

	return map[string]string{
		"activityId": likeID,
	}, nil
}

func (fediverseService *FediverseService) UndoLikeObject(userID uint, req model.LikeActionRequest) (map[string]string, error) {
	targetActor := strings.TrimSpace(req.TargetActor)
	object := strings.TrimSpace(req.Object)
	if targetActor == "" || object == "" {
		return nil, errors.New(commonModel.FEDIVERSE_INVALID_INPUT)
	}

	user, err := fediverseService.commonService.CommonGetUserByUserId(userID)
	if err != nil {
		return nil, err
	}

	actor, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return nil, err
	}

	serverURL, err := normalizeServerURL(setting.ServerURL)
	if err != nil {
		return nil, err
	}

	likeID := generateDeterministicActivityID(serverURL, actor.PreferredUsername, "like", object)
	published := time.Now().UTC()
	undoID := fmt.Sprintf("%s/activities/%s/undo-like/%d", serverURL, actor.PreferredUsername, published.UnixNano())

	payload, err := buildUndoLikeActivityPayload(&actor, targetActor, object, likeID, undoID, published)
	if err != nil {
		return nil, err
	}

	inboxURL, err := fediverseService.fetchRemoteActorInbox(targetActor)
	if err != nil {
		return nil, err
	}

	if err := httpUtil.PostActivity(payload, inboxURL, actor.ID); err != nil {
		return nil, err
	}

	return map[string]string{
		"activityId":     undoID,
		"likeActivityId": likeID,
	}, nil
}

func buildFollowActivityPayload(actor *model.Actor, targetActor string, activityID string, published time.Time) ([]byte, error) {
	if actor == nil {
		return nil, errors.New("actor is nil")
	}
	if activityID == "" {
		return nil, errors.New("activity id is empty")
	}
	if targetActor == "" {
		return nil, errors.New("target actor is empty")
	}

	payload := map[string]any{
		"@context":  []any{"https://www.w3.org/ns/activitystreams"},
		"id":        activityID,
		"type":      model.ActivityTypeFollow,
		"actor":     actor.ID,
		"object":    targetActor,
		"to":        []string{targetActor},
		"published": published.Format(time.RFC3339),
	}

	return json.Marshal(payload)
}

func buildUndoFollowActivityPayload(actor *model.Actor, targetActor string, undoID string, followActivityID string, published time.Time) ([]byte, error) {
	if actor == nil {
		return nil, errors.New("actor is nil")
	}
	if undoID == "" || followActivityID == "" {
		return nil, errors.New("activity id is empty")
	}
	if targetActor == "" {
		return nil, errors.New("target actor is empty")
	}

	payload := map[string]any{
		"@context": []any{"https://www.w3.org/ns/activitystreams"},
		"id":       undoID,
		"type":     model.ActivityTypeUndo,
		"actor":    actor.ID,
		"object": map[string]any{
			"id":     followActivityID,
			"type":   model.ActivityTypeFollow,
			"actor":  actor.ID,
			"object": targetActor,
		},
		"to":        []string{targetActor},
		"published": published.Format(time.RFC3339),
	}

	return json.Marshal(payload)
}

func buildLikeActivityPayload(actor *model.Actor, targetActor string, object string, activityID string, published time.Time) ([]byte, error) {
	if actor == nil {
		return nil, errors.New("actor is nil")
	}
	if activityID == "" {
		return nil, errors.New("activity id is empty")
	}
	if targetActor == "" || object == "" {
		return nil, errors.New("target actor or object is empty")
	}

	payload := map[string]any{
		"@context":  []any{"https://www.w3.org/ns/activitystreams"},
		"id":        activityID,
		"type":      model.ActivityTypeLike,
		"actor":     actor.ID,
		"object":    object,
		"to":        []string{targetActor},
		"published": published.Format(time.RFC3339),
	}

	return json.Marshal(payload)
}

func buildUndoLikeActivityPayload(actor *model.Actor, targetActor string, object string, likeActivityID string, undoID string, published time.Time) ([]byte, error) {
	if actor == nil {
		return nil, errors.New("actor is nil")
	}
	if likeActivityID == "" || undoID == "" {
		return nil, errors.New("activity id is empty")
	}
	if targetActor == "" || object == "" {
		return nil, errors.New("target actor or object is empty")
	}

	payload := map[string]any{
		"@context": []any{"https://www.w3.org/ns/activitystreams"},
		"id":       undoID,
		"type":     model.ActivityTypeUndo,
		"actor":    actor.ID,
		"object": map[string]any{
			"id":     likeActivityID,
			"type":   model.ActivityTypeLike,
			"actor":  actor.ID,
			"object": object,
		},
		"to":        []string{targetActor},
		"published": published.Format(time.RFC3339),
	}

	return json.Marshal(payload)
}

func generateDeterministicActivityID(serverURL, username, prefix, key string) string {
	hash := sha256.Sum256([]byte(strings.ToLower(key)))
	short := hex.EncodeToString(hash[:16])
	return fmt.Sprintf("%s/activities/%s/%s/%s", serverURL, username, prefix, short)
}
