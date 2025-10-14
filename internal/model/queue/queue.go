package model

import (
	"encoding/json"
	"time"
)

const (
	// DeadLetterTypeWebhook webhook 类型的死信任务
	DeadLetterTypeWebhook = "webhook"
	// DeadLetterTypePush push 类型的死信任务
	DeadLetterTypePush = "push"
	// DeadLetterTypeEmail email 类型的死信任务
	DeadLetterTypeEmail = "email"
)

const (
	// DeadLetterStatusPending 待处理状态
	DeadLetterStatusPending = "pending"
	// DeadLetterStatusProcessing 处理中状态
	DeadLetterStatusProcessing = "processing"
	// DeadLetterStatusFailed 失败状态
	DeadLetterStatusFailed = "failed"
	// DeadLetterStatusCompleted 完成状态
	DeadLetterStatusCompleted = "completed"
	// DeadLetterStatusDiscarded 丢弃状态
	DeadLetterStatusDiscarded = "discarded"
)

const (
	DeadLetterMetaKey = "dead_letter" // 存储在 Payload 中的元数据键
)

// ReplayPayload 重放任务的载荷，使用 map 以支持灵活的字段
type ReplayPayload map[string]any

// DeadLetter 死信任务模型
type DeadLetter struct {
	ID         int64     `db:"id"`
	Type       string    `db:"type"`        // 业务类型，如 "webhook" / "push" / "email"
	Payload    []byte    `db:"payload"`     // 原始任务数据（序列化 JSON）
	ErrorMsg   string    `db:"error_msg"`   // 失败原因（错误信息）
	RetryCount int       `db:"retry_count"` // 重试次数
	NextRetry  time.Time `db:"next_retry"`  // 下次重试时间（指数退避）
	Status     string    `db:"status"`      // 任务状态，如 "pending", "processing", "failed", "completed","discarded"
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (dl *DeadLetter) SetType(t string) {
	dl.Type = t
}

func (dl *DeadLetter) SetPayload(replay ReplayPayload) {
	// 将 ReplayPayload 序列化为 JSON 字节数组
	payload, _ := json.Marshal(replay)
	dl.Payload = payload
}

func (dl *DeadLetter) GetPayload() ReplayPayload {
	var replay ReplayPayload
	_ = json.Unmarshal(dl.Payload, &replay)
	return replay
}
