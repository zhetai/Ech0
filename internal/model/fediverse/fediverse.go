package model

import "time"

// Activities 活动表
type Activity struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ActorID      uint      `gorm:"index;not null" json:"actor_id"`
	Type         string    `gorm:"size:64;not null" json:"type"` // Create, Follow, Like, Accept...
	ActivityJSON string    `gorm:"type:text;not null" json:"activity_json"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	Delivered    bool      `gorm:"default:false" json:"delivered"`
}

// OutboxItems 发件记录表
type OutboxItem struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityID uint      `gorm:"index;not null" json:"activity_id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	DeliveredTo string   `gorm:"type:text" json:"delivered_to"` // 已送达的远端 inbox 地址
}

// Followers 关注者表
type Follower struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uint      `gorm:"index;not null" json:"user_id"`
	FollowerActorURL string   `gorm:"size:256;not null" json:"follower_actor_url"`
	Accepted         bool      `gorm:"default:false" json:"accepted"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// DeliveryQueue 投递队列表
type DeliveryQueue struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityID      uint      `gorm:"index;not null" json:"activity_id"`
	TargetInboxURL  string    `gorm:"size:512;not null" json:"target_inbox_url"`
	Tries           uint      `gorm:"default:0" json:"tries"`
	LastError       string    `gorm:"type:text" json:"last_error"`
	NextTryAt       time.Time `json:"next_try_at"`
	Status          string    `gorm:"size:32;default:'pending'" json:"status"` // pending, delivered, failed
}

// Objects 内容对象表
type Object struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ObjectID   string    `gorm:"size:512;unique;not null" json:"object_id"` // 全局唯一 URL
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	ObjectJSON string    `gorm:"type:text;not null" json:"object_json"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}