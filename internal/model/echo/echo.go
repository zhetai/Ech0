package model

import "time"

// Echo 定义Echo实体
type Echo struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Content       string    `gorm:"type:text;not null" json:"content"`
	Username      string    `gorm:"type:varchar(100)" json:"username,omitempty"`
	Images        []Image   `gorm:"foreignKey:MessageID;constraint:OnDelete:CASCADE" json:"images,omitempty"`
	Private       bool      `gorm:"default:false" json:"private"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	Extension     string    `gorm:"type:text" json:"extension,omitempty"`
	ExtensionType string    `gorm:"type:varchar(100)" json:"extension_type,omitempty"`
	FavCount      int       `gorm:"default:0" json:"fav_count"`
	CreatedAt     time.Time `json:"created_at"`
}

// Message 定义Message实体 (注意⚠️: 该模型为旧版Echo模型,新版已经弃用)
type Message struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Content       string    `gorm:"type:text;not null" json:"content"`
	Username      string    `gorm:"type:varchar(100)" json:"username,omitempty"`
	ImageURL      string    `gorm:"type:text" json:"image_url,omitempty"`
	ImageSource   string    `gorm:"type:varchar(20)" json:"image_source,omitempty"`
	Images        []Image   `gorm:"foreignKey:MessageID" json:"images,omitempty"`
	Private       bool      `gorm:"default:false" json:"private"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	Extension     string    `gorm:"type:text" json:"extension,omitempty"`
	ExtensionType string    `gorm:"type:varchar(100)" json:"extension_type,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// Image 定义Image实体
type Image struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	MessageID   uint   `gorm:"index;not null" json:"message_id"` // 关联的Echo ID(注意⚠️: 该字段名为MessageID, 但实际关联的是Echo表,因为为了兼容旧版Echo用户)
	ImageURL    string `gorm:"type:text" json:"image_url"`
	ImageSource string `gorm:"type:varchar(20)" json:"image_source"`
	ObjectKey   string `gorm:"type:text" json:"object_key,omitempty"` // 对象存储的Key (如果是本地存储则为空)
}

const (
	Extension_MUSIC      = "MUSIC"
	Extension_VIDEO      = "VIDEO"
	Extension_GITHUBPROJ = "GITHUBPROJ"
	Extension_WEBSITE    = "WEBSITE"
	ImageSourceLocal     = "local" // 本地图片
	ImageSourceURL       = "url"   // 直链图片
	ImageSourceS3        = "s3"    // S3 图片
	ImageSourceR2        = "r2"    // R2 图片
)
