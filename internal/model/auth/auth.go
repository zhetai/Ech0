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

