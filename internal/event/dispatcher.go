package event

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/lin-snow/ech0/internal/async"
	webhookModel "github.com/lin-snow/ech0/internal/model/webhook"
	webhookRepository "github.com/lin-snow/ech0/internal/repository/webhook"
	logUtil "github.com/lin-snow/ech0/internal/util/log"
)

type WebhookDispatcher struct {
	bus    IEventBus                                    // 事件总线
	client *http.Client                                 // HTTP 客户端
	repo   webhookRepository.WebhookRepositoryInterface // Webhook 仓储层
	pool   *async.WorkerPool                            // 任务池pool
}

func NewWebhookDispatcher(ebp func() IEventBus, repo webhookRepository.WebhookRepositoryInterface) *WebhookDispatcher {
	return &WebhookDispatcher{
		bus:  ebp(),
		repo: repo,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		pool: async.NewWorkerPool(6, 6), // 假设最大并发数为 6，任务队列大小为 6
	}
}

// Handle 由事件总线调用，负责调度事件到每个活跃的 webhook
func (wd *WebhookDispatcher) Handle(ctx context.Context, e *Event) error {
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
		wh := wh // 捕获变量
		// 提交任务到池中异步处理
		wd.pool.Submit(func() error {
			wd.Dispatch(ctx, &wh, e)
			return nil
		})
	}

	return nil
}

// Dispatch 负责将事件发送到指定的 webhook
func (wd *WebhookDispatcher) Dispatch(ctx context.Context, wh *webhookModel.Webhook, e *Event) {
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
	} else {
		// 记录失败日志
		logUtil.GetLogger().Error("Webhook Handle Failed: ", zap.String("name", wh.Name), zap.String("url", wh.URL))
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
	wd.pool.Wait()
}
