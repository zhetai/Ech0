package handler

import (
	"github.com/gin-gonic/gin"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	errorUtil "github.com/lin-snow/ech0/internal/util/err"
	"net/http"
)

// Response 代表 handler 层的执行结果封装
type Response struct {
	Code int
	Data any
	Msg  string
	Err  error
}

// Execute 包装器，自动根据 Response 返回统一格式的 HTTP 响应 (仅处理返回类型为JSON的handler)
func Execute(fn func(ctx *gin.Context) Response) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res := fn(ctx)
		if res.Err != nil {
			ctx.JSON(http.StatusOK, commonModel.Fail[string](
				errorUtil.HandleError(&commonModel.ServerError{
					Msg: res.Msg,
					Err: res.Err,
				}),
			))
			return
		}

		// 支持自定义 code
		if res.Code != 0 {
			ctx.JSON(http.StatusOK, commonModel.OKWithCode(res.Data, res.Code, res.Msg))
		} else {
			ctx.JSON(http.StatusOK, commonModel.OK(res.Data, res.Msg))
		}
	}
}
