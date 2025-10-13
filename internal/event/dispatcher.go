package event

import (
	"net/http"
	"sync"

	webhookRepository "github.com/lin-snow/ech0/internal/repository/webhook"
)

type WebhookDispatcher struct {
	bus    IEventBus
	client *http.Client
	repo   webhookRepository.WebhookRepositoryInterface
	wg     sync.WaitGroup
}
