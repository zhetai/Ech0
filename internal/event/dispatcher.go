package event

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap/zapcore"

	webhookModel "github.com/lin-snow/ech0/internal/model/webhook"
	webhookRepository "github.com/lin-snow/ech0/internal/repository/webhook"
	logUtil "github.com/lin-snow/ech0/internal/util/log"
)

type WebhookDispatcher struct {
	bus    IEventBus
	client *http.Client
	repo   webhookRepository.WebhookRepositoryInterface
	wg     sync.WaitGroup
}

func NewWebhookDispatcher(ebp func() IEventBus, repo webhookRepository.WebhookRepositoryInterface) *WebhookDispatcher {
	return &WebhookDispatcher{
		bus:  ebp(),
		repo: repo,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Handle 由事件总线调用，负责调度事件到每个活跃的 webhook
func (wd *WebhookDispatcher) Handle(ctx context.Context, e *Event) error {
	log.Println("Webhook Handle:", e)

	// 获取所有开启的webhook
	webhooks, err := wd.repo.ListActiveWebhooks()
	if err != nil {
		return err
	}
	for _, wh := range webhooks {
		// 可以根据 event type 做过滤
		// if !wh.ShouldHandle(e.Type) {
		// 	continue
		// }

		wd.wg.Add(1)
		go func(wh *webhookModel.Webhook) {
			defer wd.wg.Done()
			wd.Dispatch(ctx, wh, e)
		}(&wh)
	}

	return nil
}

// Dispatch 负责将事件发送到指定的 webhook
func (wd *WebhookDispatcher) Dispatch(ctx context.Context, wh *webhookModel.Webhook, e *Event) {
	log.Println("Webhook :", wh, "Event", e)

	// 构建 HTTP 请求
	req, err := wd.buildRequest(wh, e)
	if err != nil {
		// 记录日志或处理错误
		logUtil.GetLogger().Error(err.Error())
		return
	}

	// 发送 HTTP 请求
	resp, err := wd.client.Do(req)
	if err != nil {
		// 记录日志或处理错误
		logUtil.GetLogger().Error(err.Error())
		return
	}
	defer resp.Body.Close()

	// 处理响应
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// 成功处理
		logUtil.GetLogger().
			Info("Webhook %s dispatched event %s successfully", zapcore.Field{Key: "webhook", String: wh.URL}, zapcore.Field{Key: "event", String: string(e.Type)})
	} else {
		// 记录失败日志
		logUtil.GetLogger().Error("Webhook %s dispatch event %s failed with status %d", zapcore.Field{Key: "webhook", String: wh.URL}, zapcore.Field{Key: "event", String: string(e.Type)}, zapcore.Field{Key: "status", Integer: int64(resp.StatusCode)})
	}
}

// buildRequest 构建 HTTP 请求(POST)
func (wd *WebhookDispatcher) buildRequest(wh *webhookModel.Webhook, e *Event) (*http.Request, error) {
	// 构造 HTTP 请求头
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	headers.Set("X-Ech0-Event", string(e.Type))
	headers.Set("User-Agent", "Ech0-Webhook-Client")

	// 构造 HTTP 请求体
	body, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	bodyReader := io.NopCloser(bytes.NewReader(body))

	// 构造 HTTP 请求
	req, err := http.NewRequest("POST", wh.URL, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header = headers

	// 设置 GetBody 以支持重试
	req.GetBody = func() (bodyReader io.ReadCloser, err error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}

	// 返回请求对象
	return req, nil
}

// Wait 等待所有事件处理完成
func (wd *WebhookDispatcher) Wait() {
	wd.wg.Wait()
}
