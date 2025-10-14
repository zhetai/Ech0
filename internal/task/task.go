// package task declaration to use task related functionalities
package task

import (
	"context"
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/lin-snow/ech0/internal/event"
	queueRepository "github.com/lin-snow/ech0/internal/repository/queue"
	commonService "github.com/lin-snow/ech0/internal/service/common"
	logUtil "github.com/lin-snow/ech0/internal/util/log"
)

type Tasker struct {
	scheduler     gocron.Scheduler
	commonService commonService.CommonServiceInterface
	eventBus      event.IEventBus
	queueRepo     queueRepository.QueueRepositoryInterface
}

func NewTasker(
	commonService commonService.CommonServiceInterface,
	eventBusProvider func() event.IEventBus,
	queueRepo queueRepository.QueueRepositoryInterface,
) *Tasker {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		logUtil.GetLogger().Error("Failed to create scheduler", zapcore.Field{
			Key:    "error",
			String: err.Error(),
		})
	}

	return &Tasker{
		scheduler:     scheduler,
		commonService: commonService,
		eventBus:      eventBusProvider(),
		queueRepo:     queueRepo,
	}
}

func (t *Tasker) Start() {
	t.CleanupTempFilesTask()  // 启动清理临时文件任务
	t.DeadLetterConsumeTask() // 启动死信任务消费任务
	t.scheduler.Start()
}

func (t *Tasker) Stop() {
	t.scheduler.Shutdown()
}

// CleanupTempFilesTask 清理过期的临时文件任务
func (t *Tasker) CleanupTempFilesTask() {
	// 每三天执行一次
	_, err := t.scheduler.NewJob(
		gocron.DurationJob(72*time.Hour),
		gocron.NewTask(
			func() {
				if err := t.commonService.CleanupTempFiles(); err != nil {
					logUtil.GetLogger().Error("Failed to clean up temporary files", zap.String("error", err.Error()))
				}
			},
		),
	)
	if err != nil {
		logUtil.GetLogger().Error("Failed to schedule CleanupTempFilesTask", zap.String("error", err.Error()))
	}
}

// DeadLetterConsumeTask 死信任务消费任务
func (t *Tasker) DeadLetterConsumeTask() {
	// 每天12点执行一次, 测试时为每30秒执行一次
	_, err := t.scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(12, 0, 0))),
		// gocron.DurationJob(30*time.Second), // 测试时为每30秒执行一次
		gocron.NewTask(
			func() {
				// 取出死信队列中的任务，逐个重试
				deadLetters, err := t.queueRepo.ListDeadLetters(10)
				if err != nil {
					logUtil.GetLogger().Error("Failed To Get DeadLetters!", zap.String("error", err.Error()))
				}

				// 遍历死信任务，重新发送事件
				for _, dl := range deadLetters {
					// 发布事件到事件总线，触发重试
					t.eventBus.Publish(context.Background(),
						event.NewEvent(
							event.EventTypeDeadLetterRetried,
							event.EventPayload{
								event.EventPayloadDeadLetter: dl,
							},
						),
					)
				}
			},
		),
	)
	if err != nil {
		logUtil.GetLogger().Error("Failed to schedule WebhookRetryTask", zap.String("error", err.Error()))
	}
}
