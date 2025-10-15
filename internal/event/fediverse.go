package event

import (
	"context"
	"time"

	"github.com/lin-snow/ech0/internal/async"
	"github.com/lin-snow/ech0/internal/fediverse"
	echoModel "github.com/lin-snow/ech0/internal/model/echo"
	userModel "github.com/lin-snow/ech0/internal/model/user"
	logUtil "github.com/lin-snow/ech0/internal/util/log"
)

// FediverseAgent 处理联邦相关的事件
type FediverseAgent struct {
	pool *async.WorkerPool // 任务池
	core *fediverse.FediverseCore
}

func NewFediverseAgent(core *fediverse.FediverseCore) *FediverseAgent {
	return &FediverseAgent{
		pool: async.NewWorkerPool(3, 3), // 假设最大并发数为 3，任务队列大小为 3
		core: core,
	}
}

func (fa *FediverseAgent) Handle(ctx context.Context, e *Event) error {
	// 处理事件，与联邦宇宙交互
	switch e.Type {
	case EventTypeEchoCreated:
		fa.HandleCreateEchoEvent(ctx, e)

	default:
		return nil // 忽略其他事件
	}

	return nil
}

func (fa *FediverseAgent) Wait() {
	fa.pool.Wait()
}

func (fa *FediverseAgent) HandleCreateEchoEvent(ctx context.Context, e *Event) error {
	// 将 Echo 推送到联邦宇宙
	payload := e.Payload
	echoData, ok := payload[EventPayloadEcho]
	if !ok {
		return nil
	}
	echo, ok := echoData.(echoModel.Echo)
	if !ok {
		return nil
	}
	userData, ok := payload[EventPayloadUser]
	if !ok {
		return nil
	}
	user, ok := userData.(userModel.User)
	if !ok {
		return nil
	}

	fa.pool.Submit(func() error {
		// 重试机制，最多重试3次，初始延迟1秒
		return fa.retryWithBackoff(3, time.Second, func() error {
			if err := fa.core.PushEchoToFediverse(user.ID, echo); err != nil {
				logUtil.GetLogger().Error(err.Error())
				return err
			}
			return nil
		})
	})

	return nil
}

func (fa *FediverseAgent) retryWithBackoff(retries int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < retries; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(delay)
		delay *= 2
	}
	return err
}
