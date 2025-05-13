package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Username  string    `gorm:"type:varchar(100)" json:"username,omitempty"`
	ImageURL  string    `gorm:"type:text" json:"image_url,omitempty"`
	Private   bool      `gorm:"default:false" json:"private"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Todo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Username  string    `gorm:"type:varchar(100)" json:"username,omitempty"`
	Status    uint      `gorm:"default:0" json:"status"` // 0:未完成 1:已完成
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
	IsAdmin  bool   `json:"is_admin"`
	Avatar   string `gorm:"size:255" json:"avatar"`
}

type SystemSetting struct {
	SiteTitle     string `json:"site_title"`     // 站点标题
	ServerName    string `json:"server_name"`    // 服务器名称
	ServerURL     string `json:"server_url"`     // 服务器地址
	AllowRegister bool   `json:"allow_register"` // 是否允许注册'
	ICPNumber     string `json:"ICP_number"`     // 备案号
}

type UserStatus struct {
	UserID   uint   `json:"user_id"`  // 用户ID
	UserName string `json:"username"` // 用户名
	IsAdmin  bool   `json:"is_admin"` // 是否是管理员
}

type Heapmap struct {
	Date  string `json:"date"`  // 日期
	Count int    `json:"count"` // 留言数量
}

type Status struct {
	SysAdminID    uint         `json:"sys_admin_id"` // 系统管理员ID
	Username      string       `json:"username"`     // 系统管理员用户名
	Users         []UserStatus `json:"users"`        // 所有用户
	Logo          string       `json:"logo"`         // 站点logo
	TotalMessages int          `json:"total_messages"`
}

type Connect struct {
	ServerName  string `json:"server_name"`  // 服务器名称
	ServerURL   string `json:"server_url"`   // 服务器地址
	Logo        string `json:"logo"`         // 站点logo
	Ech0s       int    `json:"ech0s"`        // 留言数量
	SysUsername string `json:"sys_username"` // 系统管理员用户名
}

type MyCliams struct {
	Userid   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// key value表
type KeyValue struct {
	Key   string `json:"key" gorm:"primaryKey"`
	Value string `json:"value"`
}
