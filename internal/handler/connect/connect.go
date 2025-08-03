package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	res "github.com/lin-snow/ech0/internal/handler/response"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/connect"
	service "github.com/lin-snow/ech0/internal/service/connect"
)

type ConnectHandler struct {
	connectService service.ConnectServiceInterface
}

// NewConnectHandler ConnectHandler 的构造函数
func NewConnectHandler(connectService service.ConnectServiceInterface) *ConnectHandler {
	return &ConnectHandler{
		connectService: connectService,
	}
}

// AddConnect 添加连接
func (connectHandler *ConnectHandler) AddConnect() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		userId := ctx.MustGet("userid").(uint)

		var connected model.Connected
		if err := ctx.ShouldBindJSON(&connected); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
			}
		}

		if err := connectHandler.connectService.AddConnect(userId, connected); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.ADD_CONNECT_SUCCESS,
		}
	})

}

// DeleteConnect 删除连接
func (connectHandler *ConnectHandler) DeleteConnect() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		userId := ctx.MustGet("userid").(uint)

		// 从 URL 参数获取 ID
		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return res.Response{
				Msg: commonModel.INVALID_PARAMS,
			}
		}

		if err := connectHandler.connectService.DeleteConnect(userId, uint(id)); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.DELETE_CONNECT_SUCCESS,
		}
	})
}

// GetConnectsInfo 获取所有添加的连接的信息
func (connectHandler *ConnectHandler) GetConnectsInfo() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 调用 Service 层获取 Connect 信息
		connects, err := connectHandler.connectService.GetConnectsInfo()
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: connects,
			Msg:  commonModel.GET_CONNECT_INFO_SUCCESS,
		}
	})
}

// GetConnect 提供当前实例的连接信息
func (connectHandler *ConnectHandler) GetConnect() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		connect, err := connectHandler.connectService.GetConnect()
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: connect,
			Msg:  commonModel.CONNECT_SUCCESS,
		}
	})
}

// GetConnects 获取当前实例添加的所有连接
func (connectHandler *ConnectHandler) GetConnects() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 调用 Service 层获取 Connect 列表
		connects, err := connectHandler.connectService.GetConnects()
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: connects,
			Msg:  commonModel.GET_CONNECTED_LIST_SUCCESS,
		}
	})
}
