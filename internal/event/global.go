package event

import (
	"sync/atomic"
)

// 使用 atomic.Value 存储全局 EventBus，线程安全
var globalEventBus atomic.Value // 存储 *EventBus

// GetEventBus 获取全局 EventBus 实例
func GetEventBus() IEventBus {
	v := globalEventBus.Load()
	if v == nil {
		panic("EventBus not initialized. Call InitEventBus() first.")
	}
	return v.(IEventBus)
}

// SetEventBus 设置全局 EventBus 实例（仅框架启动时调用）
func SetEventBus(bus IEventBus) {
	globalEventBus.Store(bus)
}

// InitEventBus 初始化并注册全局 EventBus
func InitEventBus() {
	bus := NewEventBus()
	SetEventBus(bus)
}

// Provider 用于 wire 构建
// func EventBusProvider() func() IEventBus {
// 	return GetEventBus
// }
