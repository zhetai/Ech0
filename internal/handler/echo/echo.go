package handler

import (
	"github.com/gin-gonic/gin"
	res "github.com/lin-snow/ech0/internal/handler/response"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/echo"
	service "github.com/lin-snow/ech0/internal/service/echo"
	"strconv"
)

type EchoHandler struct {
	echoService service.EchoServiceInterface
}

func NewEchoHandler(echoService service.EchoServiceInterface) *EchoHandler {
	return &EchoHandler{
		echoService: echoService,
	}
}

func (echoHandler *EchoHandler) PostEcho() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		var newEcho model.Echo
		if err := ctx.ShouldBindJSON(&newEcho); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		userId := ctx.MustGet("userid").(uint)
		if err := echoHandler.echoService.PostEcho(userId, &newEcho); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.POST_ECHO_SUCCESS,
		}
	})

}

func (echoHandler *EchoHandler) GetEchosByPage() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取分页参数
		var pageRequest commonModel.PageQueryDto
		if err := ctx.ShouldBindJSON(&pageRequest); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)
		result, err := echoHandler.echoService.GetEchosByPage(userid, pageRequest)
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: result,
			Msg:  commonModel.GET_ECHOS_BY_PAGE_SUCCESS,
		}
	})

}

func (echoHandler *EchoHandler) DeleteEcho() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		// 从 URL 参数获取留言 ID
		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return res.Response{
				Msg: commonModel.INVALID_PARAMS,
			}
		}

		if err := echoHandler.echoService.DeleteEchoById(userid, uint(id)); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.DELETE_ECHO_SUCCESS,
		}
	})
}

func (echoHandler *EchoHandler) GetTodayEchos() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)
		result, err := echoHandler.echoService.GetTodayEchos(userid)
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: result,
			Msg:  commonModel.GET_TODAY_ECHOS_SUCCESS,
		}
	})
}
