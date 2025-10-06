package model

import (
	"github.com/golang-jwt/jwt/v5"
)

// MyClaims 是自定义的 JWT 声明结构体
type MyClaims struct {
	Userid   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

const (
	// MAX_USER_COUNT 定义最大用户数量
	MAX_USER_COUNT = 5
	// NO_USER_LOGINED 定义未登录用户的 ID
	NO_USER_LOGINED = uint(0)
)

// OAuth2Setting 定义 OAuth2 配置结构体
type OAuth2Setting struct {
	Enable       bool     `json:"enable"`
	Provider     string   `json:"provider"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURI  string   `json:"redirect_uri"`
	Scopes       []string `json:"scopes"`
	AuthURL      string   `json:"auth_url"`
	TokenURL     string   `json:"token_url"`
	UserInfoURL  string   `json:"user_info_url"`
}
