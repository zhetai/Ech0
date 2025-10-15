package event

import (
	"context"

	"github.com/lin-snow/ech0/internal/async"
)

// FediverseAgent 处理联邦相关的事件
type FediverseAgent struct {
	pool *async.WorkerPool // 任务池
}

func NewFediverseAgent() *FediverseAgent {
	return &FediverseAgent{
		pool: async.NewWorkerPool(3, 3), // 假设最大并发数为 3，任务队列大小为 3
	}
}

func (fa *FediverseAgent) Handle(ctx context.Context, event *Event) error {
	// 处理事件，与联邦宇宙交互
	switch event.Type {
	case EventTypeEchoCreated:

	default:
		return nil // 忽略其他事件
	}

	return nil
}

func (fa *FediverseAgent) Wait() {
	fa.pool.Wait()
}
