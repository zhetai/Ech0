package model

import "time"

// ActivityPubError 定义 ActivityPub 协议的错误响应格式
type ActivityPubError struct {
	Context string `json:"@context"` // ActivityStreams 上下文
	Type    string `json:"type"`     // 类型，固定为 "Error"
	Error   string `json:"error"`    // 错误信息
	Status  int    `json:"status"`   // HTTP 状态码
}

// ActivityType 定义常见的 ActivityPub 活动类型
const (
	ActivityTypeCreate   string = "Create"
	ActivityTypeFollow   string = "Follow"
	ActivityTypeLike     string = "Like"
	ActivityTypeAccept   string = "Accept"
	ActivityTypeAnnounce string = "Announce"
	ActivityTypeUndo     string = "Undo"
)

// WebFingerResponse 是 WebFinger 返回的标准结构
type WebFingerResponse struct {
	Subject string   `json:"subject"`           // 用户标识，例如 acct:alice@domain.com
	Aliases []string `json:"aliases,omitempty"` // 可选：用户的别名 URL
	Links   []Link   `json:"links"`             // 与用户相关的资源链接
}

// Link 是 WebFingerResponse 中的链接对象
type Link struct {
	Rel  string `json:"rel"`            // 链接关系，例如 "self"
	Type string `json:"type,omitempty"` // MIME 类型，例如 application/activity+json
	Href string `json:"href"`           // 链接 URL
}

// OutboxResponse 定义 Outbox 的响应格式
type OutboxResponse struct {
	Context    interface{} `json:"@context"` // 可以是字符串或数组
	ID         string      `json:"id"`
	Type       string      `json:"type"` // "OrderedCollection"
	TotalItems int         `json:"totalItems"`
	First      *OutboxPage `json:"first,omitempty"`
	Last       *OutboxPage `json:"last,omitempty"`
	// 如果不用分页，可以直接用 OrderedItems
	OrderedItems []Activity `json:"orderedItems,omitempty"`
}

// OutboxPage 表示分页形式的 Outbox
type OutboxPage struct {
	ID           string     `json:"id"`
	Type         string     `json:"type"` // "OrderedCollectionPage"
	PartOf       string     `json:"partOf"`
	Next         string     `json:"next,omitempty"`
	Prev         string     `json:"prev,omitempty"`
	OrderedItems []Activity `json:"orderedItems"`
}

// ------------------ 数据库模型 --------------------

type DeliveryStatus string

const (
	DeliveryStatusPending   DeliveryStatus = "pending"
	DeliveryStatusDelivered DeliveryStatus = "delivered"
	DeliveryStatusFailed    DeliveryStatus = "failed"
)

// Activities 活动表 (面向外部的 ActivityPub 活动)
type Activity struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ActorID      uint      `gorm:"index;not null" json:"actor_id"`
	Type         string    `gorm:"size:64;not null" json:"type"` // Create, Follow, Like, Accept...
	ActivityJSON string    `gorm:"type:text;not null" json:"activity_json"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	Delivered    bool      `gorm:"default:false" json:"delivered"`
}

// OutboxItems 发件记录表 (记录哪些活动发给了哪些用户)
type OutboxItem struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityID  uint      `gorm:"index;not null" json:"activity_id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	DeliveredTo string    `gorm:"type:text" json:"delivered_to"` // 已送达的远端 inbox 地址
}

// Followers 关注者表
type Follower struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uint      `gorm:"index;not null" json:"user_id"`
	FollowerActorURL string    `gorm:"size:256;not null" json:"follower_actor_url"`
	Accepted         bool      `gorm:"default:false" json:"accepted"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// DeliveryQueue 投递队列表
type DeliveryQueue struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityID     uint      `gorm:"index;not null" json:"activity_id"`
	TargetInboxURL string    `gorm:"size:512;not null" json:"target_inbox_url"`
	Tries          uint      `gorm:"default:0" json:"tries"`
	LastError      string    `gorm:"type:text" json:"last_error"`
	NextTryAt      time.Time `json:"next_try_at"`
	Status         string    `gorm:"size:32;default:'pending'" json:"status"` // pending, delivered, failed
}

// Objects 内容对象表
type Object struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ObjectID   string    `gorm:"size:512;unique;not null" json:"object_id"` // 全局唯一 URL
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	ObjectJSON string    `gorm:"type:text;not null" json:"object_json"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// ---------------- 常用数据模型 --------------------

// PublicKey 公钥信息
type PublicKey struct {
	ID           string `json:"id"`           // 公钥 ID，通常为 Actor ID + "#main-key"
	Owner        string `json:"owner"`        // 所有者，通常为 Actor ID
	PublicKeyPem string `json:"publicKeyPem"` // PEM 格式的公钥
	Type         string `json:"type"`         // "PublicKey"
}

// Actor ActivityPub Actor 信息
type Actor struct {
	Context           []interface{} `json:"@context"`          // 上下文，可以是字符串或对象的数组
	ID                string        `json:"id"`                // Actor 的唯一标识 URL
	Type              string        `json:"type"`              // Actor 类型，通常为 "Person"
	Name              string        `json:"name"`              // 显示名称
	PreferredUsername string        `json:"preferredUsername"` // 用户名
	Summary           string        `json:"summary"`           // 简短介绍
	Inbox             string        `json:"inbox"`             // 收件箱 URL
	Outbox            string        `json:"outbox"`            // 发件箱 URL
	PublicKey         PublicKey     `json:"publicKey"`         // 公钥信息
}
