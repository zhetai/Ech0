package util

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/lin-snow/ech0/internal/config"
	authModel "github.com/lin-snow/ech0/internal/model/auth"
	userModel "github.com/lin-snow/ech0/internal/model/user"
	cryptoUtil "github.com/lin-snow/ech0/internal/util/crypto"
)

// CreateClaims 创建Claims
func CreateClaims(user userModel.User) jwt.Claims {
	leeway := time.Second * 60 // 允许的时间偏差
	claims := authModel.MyClaims{
		Userid:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Config.Auth.Jwt.Issuer,
			Subject:   user.Username,
			Audience:  jwt.ClaimStrings{config.Config.Auth.Jwt.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Config.Auth.Jwt.Expires) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-leeway)),
		},
	}

	return claims
}

// CreateClaims 创建Claims 带过期时间
func CreateClaimsWithExpiry(user userModel.User, expiry int64) jwt.Claims {
	leeway := time.Second * 60 // 允许的时间偏差
	claims := authModel.MyClaims{
		Userid:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Config.Auth.Jwt.Issuer,
			Subject:   user.Username,
			Audience:  jwt.ClaimStrings{config.Config.Auth.Jwt.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-leeway)),
		},
	}

	return claims
}

// GenerateToken 生成JWT Token
func GenerateToken(claim jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(config.JWT_SECRET)
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*authModel.MyClaims, error) {
	claims := &authModel.MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_SECRET, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*authModel.MyClaims); ok {
		return claims, nil
	}

	log.Println("unknown claims type, cannot proceed")
	return nil, errors.New("unknown claims type, cannot proceed")
}

// GenerateOAuthState 生成 OAuth2 state token
func GenerateOAuthState(action string, userID uint, redirect, provider string) (string, error) {
	now := time.Now()
	expiration := now.Add(10 * time.Minute)

	claims := jwt.MapClaims{
		"action":   action,
		"user_id":  userID,
		"nonce":    cryptoUtil.GenerateRandomString(16),
		"redirect": redirect,
		"exp":      expiration.Unix(),
		"iat":      now.Unix(),
		"provider": provider,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWT_SECRET)
}

// ParseOAuthState 解析并验证 OAuth2 state token
func ParseOAuthState(stateStr string) (*authModel.OAuthState, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(stateStr, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_SECRET, nil
	})
	if err != nil {
		return nil, err
	}

	return &authModel.OAuthState{
		Action:   claims["action"].(string),
		UserID:   uint(claims["user_id"].(float64)),
		Nonce:    claims["nonce"].(string),
		Redirect: claims["redirect"].(string),
		Exp:      int64(claims["exp"].(float64)),
		Provider: claims["provider"].(string),
	}, nil
}
