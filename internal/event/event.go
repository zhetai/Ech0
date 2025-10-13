package event

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

// 定义事件类型
const (
	EventTypeUserCreated EventType = "user.created" // 创建用户
	EventTypeUserUpdated EventType = "user.updated" // 更新用户
	EventTypeUserDeleted EventType = "user.deleted" // 删除用户

	EventTypeEchoCreated EventType = "echo.created" // 创建Echo
	EventTypeEchoUpdated EventType = "echo.updated" // 更新Echo
	EventTypeEchoDeleted EventType = "echo.deleted" // 删除Echo

	EventTypeResourceUploaded EventType = "resource.uploaded" // 资源上传

	EventTypeSystemBackup  EventType = "system.backup"  // 系统快照备份
	EventTypeSystemRestore EventType = "system.restore" // 系统快照恢复
	EventTypeSystemExport  EventType = "system.export"  // 系统快照导出
)

// 定义事件Payload的常用字段
const (
	EventPayloadUser = "user"
	EventPayloadEcho = "echo"
	EventPayloadData = "data"
	EventPayloadInfo = "info"
	EventPayloadMeta = "meta"
	EventPayloadTime = "time"
	EventPayloadType = "type"
	EventPayloadID   = "id"
	EventPayloadURL  = "url"
	EventPayloadSize = "size"
	EventPayloadPath = "path"
	EventPayloadFile = "file"
)

// Event 事件结构体
type Event struct {
	ID        string       `json:"id"`        // Unique event ID
	Type      EventType    `json:"type"`      // Event type
	Payload   EventPayload `json:"data"`      // Event payload
	Timestamp time.Time    `json:"timestamp"` // Event timestamp
}

type (
	EventType    string         // 事件类型
	EventPayload map[string]any // 事件负载
)

// NewEvent 创建一个新的事件
func NewEvent(eventType EventType, payload EventPayload) *Event {
	id := ulid.MustNew(ulid.Timestamp(time.Now()), rand.New(rand.NewSource(time.Now().UnixNano()))).
		String()
		// 使用 ULID 生成唯一 ID
	return &Event{
		ID:        id,
		Type:      eventType,
		Payload:   payload,
		Timestamp: time.Now(),
	}
}

// IEventBus 事件总线接口
type IEventBus interface {
	Publish(ctx context.Context, event *Event) error           // 发布事件
	Subscribe(eventType EventType, handler EventHandler) error // 订阅特定事件
	SubscribeAll(handler EventHandler) error                   // 订阅所有事件
}

// EventHandler 事件处理函数类型
type EventHandler func(ctx context.Context, event *Event) error

// EventBus 事件总线
type EventBus struct {
	mu   sync.RWMutex                 // 读写锁保护订阅者列表
	subs map[EventType][]EventHandler // 订阅者列表
	all  []EventHandler               // 全部事件的订阅者
}

// NewEventBus 创建一个新的事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		subs: make(map[EventType][]EventHandler),
	}
}

// Publish 发布事件
func (eb *EventBus) Publish(ctx context.Context, event *Event) error {
	eb.mu.RLock()
	handlers, ok := eb.subs[event.Type]
	eb.mu.RUnlock()
	if !ok {
		handlers = []EventHandler{}
	}

	// 处理所有订阅该事件类型的处理器
	for _, handler := range handlers {
		go func(h EventHandler) {
			if err := h(ctx, event); err != nil {
				// 错误处理
				log.Println("Event Handler Error:", err)
			}
		}(handler)
	}

	eb.mu.RLock()
	allHandlers := eb.all
	eb.mu.RUnlock()

	// 处理所有订阅所有事件的处理器
	for _, handler := range allHandlers {
		go func(h EventHandler) {
			if err := h(ctx, event); err != nil {
				// 错误处理
				log.Println("Event Handler Error:", err)
			}
		}(handler)
	}

	return nil
}

// Subscribe 订阅指定事件
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subs[eventType] = append(eb.subs[eventType], handler)
	return nil
}

// SubscribeAll 订阅所有事件
func (eb *EventBus) SubscribeAll(handler EventHandler) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.all = append(eb.all, handler)
	return nil
}
