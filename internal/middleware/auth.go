package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/pkg"
)

// JWT Auth Middleware
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
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
