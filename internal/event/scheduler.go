package event

import (
	"context"
)

type BackupScheduler struct{}

func NewBackupScheduler() *BackupScheduler {
	return &BackupScheduler{}
}

func (bs *BackupScheduler) Handle(ctx context.Context, e *Event) error {
	// 处理更新备份计划事件
	return nil
}
