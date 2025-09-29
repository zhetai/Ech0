// package task declaration to use task related functionalities
package task

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	commonService "github.com/lin-snow/ech0/internal/service/common"
	logUtil "github.com/lin-snow/ech0/internal/util/log"
	"go.uber.org/zap/zapcore"
)

type Tasker struct {
	scheduler     gocron.Scheduler
	commonService commonService.CommonServiceInterface
}

func NewTasker(commonService commonService.CommonServiceInterface) *Tasker {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		logUtil.GetLogger().Error("Failed to create scheduler", zapcore.Field{
			Key:  "error",
			String: err.Error(),
		})
	}

	return &Tasker{
		scheduler:     scheduler,
		commonService: commonService,
	}
}

func (t *Tasker) Start() {
	t.CleanupTempFilesTask() // 启动清理临时文件任务
	t.scheduler.Start()
}

func (t *Tasker) Stop() {
	t.scheduler.Shutdown()
}

// CleanupTempFilesTask 清理过期的临时文件任务
func (t *Tasker) CleanupTempFilesTask() {
	// 每三天执行一次
	_, err := t.scheduler.NewJob(
		gocron.DurationJob(72 * time.Hour),
		gocron.NewTask(
			func() {
				if err := t.commonService.CleanupTempFiles(); err != nil {
					logUtil.GetLogger().Error("Failed to clean up temporary files", zapcore.Field{
						Key:  "error",
						String: err.Error(),
					})
				}
			},
		),
	)
	if err != nil {
		logUtil.GetLogger().Error("Failed to schedule CleanupTempFilesTask", zapcore.Field{
			Key:  "error",
			String: err.Error(),
		})
	}
}