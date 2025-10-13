package event

// EventRegistrar 事件注册器
type EventRegistrar struct {
	eb IEventBus // 事件总线
}

// EventHandlers 事件监听者集合
type EventHandlers struct {
	Webhook WebhookDispatcher // Webhook 事件分发器
}

// NewEventRegistry 创建一个新的事件注册表
func NewEventRegistry(eb IEventBus) *EventRegistrar {
	return &EventRegistrar{eb: eb}
}

// Register 注册事件处理函数
func (er *EventRegistrar) Register(eventType EventType, handler EventHandler) error {
	return nil
}
