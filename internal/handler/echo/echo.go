package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	res "github.com/lin-snow/ech0/internal/handler/response"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/echo"
	service "github.com/lin-snow/ech0/internal/service/echo"
)

type EchoHandler struct {
	echoService service.EchoServiceInterface
}

func NewEchoHandler(echoService service.EchoServiceInterface) *EchoHandler {
	return &EchoHandler{
		echoService: echoService,
	}
}

// PostEcho 创建新的Echo
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

// GetEchosByPage 获取Echo列表，支持分页
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

// DeleteEcho 删除Echo
func (echoHandler *EchoHandler) DeleteEcho() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		// 从 URL 参数获取Echo ID
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

// GetTodayEchos 获取今天的Echo列表
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

// UpdateEcho 更新Echo
func (echoHandler *EchoHandler) UpdateEcho() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		var updateEcho model.Echo
		if err := ctx.ShouldBindJSON(&updateEcho); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		userId := ctx.MustGet("userid").(uint)
		if err := echoHandler.echoService.UpdateEcho(userId, &updateEcho); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.UPDATE_ECHO_SUCCESS,
		}
	})
}

// LikeEcho 点赞Echo
func (echoHandler *EchoHandler) LikeEcho() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 从 URL 参数获取Echo ID
		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return res.Response{
				Msg: commonModel.INVALID_PARAMS,
			}
		}

		if err := echoHandler.echoService.LikeEcho(uint(id)); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.LIKE_ECHO_SUCCESS,
		}
	})
}

// GetEchoById 获取指定 ID 的 Echo
func (echoHandler *EchoHandler) GetEchoById() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 从 URL 参数获取Echo ID
		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return res.Response{
				Msg: commonModel.INVALID_PARAMS,
			}
		}

		echo, err := echoHandler.echoService.GetEchoById(uint(id))
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: echo,
			Msg:  commonModel.GET_ECHO_BY_ID_SUCCESS,
		}
	})
}
