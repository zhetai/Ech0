package event

import (
	"context"
	"net/http"
	"sync"

	webhookModel "github.com/lin-snow/ech0/internal/model/webhook"
	webhookRepository "github.com/lin-snow/ech0/internal/repository/webhook"
)

type WebhookDispatcher struct {
	bus    IEventBus
	client *http.Client
	repo   webhookRepository.WebhookRepositoryInterface
	wg     sync.WaitGroup
}

func NewWebhookDispatcher(bus IEventBus, repo webhookRepository.WebhookRepositoryInterface) *WebhookDispatcher {
	return &WebhookDispatcher{
		bus:  bus,
		repo: repo,
	}
}

func (wd *WebhookDispatcher) Handle(ctx context.Context, e *Event) error {
	// 获取所有的webhook

	return nil
}

func (wd *WebhookDispatcher) Dispatch(ctx context.Context, wh *webhookModel.Webhook, e *Event) {
}

func (wd *WebhookDispatcher) Wait() {
	wd.wg.Wait()
}
