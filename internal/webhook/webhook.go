package webhook

import (
	"net/http"
	"sync"

	"github.com/lin-snow/ech0/internal/event"
	webhookRepository "github.com/lin-snow/ech0/internal/repository/webhook"
)

type WebhookDispatcher struct {
	bus    event.IEventBus
	client *http.Client
	repo   webhookRepository.WebhookRepositoryInterface
	wg     sync.WaitGroup
}
