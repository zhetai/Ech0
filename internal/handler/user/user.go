package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	res "github.com/lin-snow/ech0/internal/handler/response"
	authModel "github.com/lin-snow/ech0/internal/model/auth"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/user"
	service "github.com/lin-snow/ech0/internal/service/user"
)

type UserHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Login 用户登陆
func (userHandler *UserHandler) Login() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 从请求体获取用户名和密码
		var loginDto authModel.LoginDto
		if err := ctx.ShouldBindJSON(&loginDto); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		// 调用 Service 层处理登陆
		token, err := userHandler.userService.Login(&loginDto)
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		// 返回成功响应， 包含 JWT Token
		return res.Response{
			Data: token,
			Msg:  commonModel.LOGIN_SUCCESS,
		}
	})

}

// Register 用户注册
func (userHandler *UserHandler) Register() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		var registerDto authModel.RegisterDto
		if err := ctx.ShouldBindJSON(&registerDto); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		// 调用 Service 层处理注册
		if err := userHandler.userService.Register(&registerDto); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.REGISTER_SUCCESS,
		}
	})

}

// UpdateUser 更新用户信息
func (userHandler *UserHandler) UpdateUser() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 解析用户请求体中的参数
		var userdto model.UserInfoDto
		if err := ctx.ShouldBindJSON(&userdto); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)
		if err := userHandler.userService.UpdateUser(userid, userdto); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.UPDATE_USER_SUCCESS,
		}
	})

}

// UpdateUserAdmin 更新用户权限
func (userHandler *UserHandler) UpdateUserAdmin() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		if err := userHandler.userService.UpdateUserAdmin(userid); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.UPDATE_USER_SUCCESS,
		}
	})

}

// GetAllUsers 获取所有用户
func (userHandler *UserHandler) GetAllUsers() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		allusers, err := userHandler.userService.GetAllUsers()
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: allusers,
			Msg:  commonModel.GET_USER_SUCCESS,
		}
	})

}

// DeleteUser 删除用户
func (userHandler *UserHandler) DeleteUser() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return res.Response{
				Msg: commonModel.INVALID_PARAMS,
				Err: err,
			}
		}

		if err := userHandler.userService.DeleteUser(userid, uint(id)); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.DELETE_USER_SUCCESS,
		}
	})

}

// GetUserInfo 获取当前用户信息
func (userHandler *UserHandler) GetUserInfo() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		// 调用 Service 层获取用户信息
		user, err := userHandler.userService.GetUserByID(int(userid))
		user.Password = "" // 不返回密码
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		// 返回成功响应
		return res.Response{
			Data: user,
			Msg:  commonModel.GET_USER_INFO_SUCCESS,
		}
	})

}
