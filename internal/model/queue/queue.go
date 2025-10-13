package model

import "time"

type DeadLetter struct {
	ID         int64     `db:"id"`
	Type       string    `db:"type"`        // 业务类型，如 "webhook" / "push" / "email"
	Payload    []byte    `db:"payload"`     // 原始任务数据（序列化 JSON）
	ErrorMsg   string    `db:"error_msg"`   // 失败原因（错误信息）
	RetryCount int       `db:"retry_count"` // 重试次数
	NextRetry  time.Time `db:"next_retry"`  // 下次重试时间（指数退避）
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
