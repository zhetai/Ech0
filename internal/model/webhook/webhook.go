package model

import "time"

// Webhook 定义 Webhook 设置实体
type Webhook struct {
	ID          uint      `gorm:"primaryKey"   json:"id"`           // Webhook ID
	Name        string    `                    json:"name"`         // Webhook 名称
	URL         string    `                    json:"url"`          // Webhook URL
	Secret      string    `                    json:"secret"`       // 签名密钥，用于请求验证（HMAC等）
	IsActive    bool      `gorm:"default:true" json:"is_active"`    // 启用/禁用状态
	LastStatus  string    `                    json:"last_status"`  // 最近调用状态（如 success, failed）
	LastTrigger time.Time `                    json:"last_trigger"` // 最近触发时间
	CreatedAt   time.Time `                    json:"created_at"`   // 创建时间
	UpdatedAt   time.Time `                    json:"updated_at"`   // 更新时间
}
