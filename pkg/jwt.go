package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lin-snow/ech0/config"
	"github.com/lin-snow/ech0/internal/models"
)

// 创建Claims
func CreateClaims(user models.User) jwt.Claims {
	claims := models.MyCliams{
		Userid:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Config.Auth.Jwt.Issuer,
			Subject:   user.Username,
			Audience:  jwt.ClaimStrings{config.Config.Auth.Jwt.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Config.Auth.Jwt.Expires) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	return claims
}

// 生成JWT Token
func GenerateToken(claim jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(config.Config.Auth.Jwt.Secret))
}

// 解析JWT Token
func ParseToken(tokenString string) (*models.MyCliams, error) {
	claims := &models.MyCliams{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.Auth.Jwt.Secret), nil
	})

	// 验证签名错误
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*models.MyCliams); ok {
		return claims, nil
	} else {
		log.Println("unknown claims type, cannot proceed")
	}

	return nil, nil
}

// MD5 加密
func MD5Encrypt(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}
