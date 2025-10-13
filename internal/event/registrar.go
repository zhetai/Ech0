package event

// EventHandlers 事件处理器集合
type EventHandlers struct {
	wbd *WebhookDispatcher // webhook 事件处理器
}

// NewEventHandlers 创建一个新的事件处理器集合
func NewEventHandlers(wbd *WebhookDispatcher) *EventHandlers {
	return &EventHandlers{wbd: wbd}
}

// EventRegistrar 事件注册器
type EventRegistrar struct {
	eb IEventBus // 事件总线
	eh *EventHandlers // 事件处理器集合
}

// NewEventRegistry 创建一个新的事件注册表
func NewEventRegistry(eb IEventBus, eh *EventHandlers) *EventRegistrar {
	return &EventRegistrar{eb: eb, eh: eh}
}

// Register 注册事件处理函数
func (er *EventRegistrar) Register(eventType EventType) error {
	// 系统级 Dispatcher

	// 业务级 Dispatcher
    er.eb.SubscribeAll(er.eh.wbd.Handle) // 订阅所有事件，交给 WebhookDispatcher 处理
	
	return nil
}
