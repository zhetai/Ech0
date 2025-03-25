package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/pkg"
)

// JWT Auth Middleware
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")

		if auth == "" || auth == "null" || auth == "undefined" || len(auth) < 7 {
			// 如果只是分页获取首页留言，则不需要鉴权
			if strings.HasPrefix(ctx.Request.URL.Path, "/api/messages/page") {
				ctx.Set("userid", uint(0))
				ctx.Next()
				return
			}

			// 查看留言详情也不需要鉴权
			if strings.HasPrefix(ctx.Request.URL.Path, "/api/messages/") {
				ctx.Set("userid", uint(0))
				ctx.Next()
				return
			}

			ctx.JSON(http.StatusOK, dto.Fail[any](models.TokenNotFoundMessage))
			ctx.Abort()
			return
		}

		mc, err := pkg.ParseToken(auth)
		if err != nil {
			ctx.JSON(http.StatusOK, dto.Fail[any](models.TokenInvalidMessage))
			ctx.Abort()
			return
		}

		ctx.Set("userid", mc.Userid)
		ctx.Next()
	}
}
