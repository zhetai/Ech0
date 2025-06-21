package handler

import (
	"github.com/gin-gonic/gin"
	res "github.com/lin-snow/ech0/internal/handler/response"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/todo"
	service "github.com/lin-snow/ech0/internal/service/todo"
	"strconv"
)

type TodoHandler struct {
	todoService service.TodoServiceInterface
}

func NewTodoHandler(todoService service.TodoServiceInterface) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

func (todoHandler *TodoHandler) AddTodo() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		var todo model.Todo
		if err := ctx.ShouldBindJSON(&todo); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		if err := todoHandler.todoService.AddTodo(userid, &todo); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.ADD_TODO_SUCCESS,
		}
	})

}

func (todoHandler *TodoHandler) UpdateTodo() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		// 从 URL 参数获取留言 ID
		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return res.Response{
				Msg: commonModel.INVALID_PARAMS_BODY,
				Err: err,
			}
		}

		if err := todoHandler.todoService.UpdateTodo(userid, int64(id)); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.UPDATE_TODO_SUCCESS,
		}
	})

}

func (todoHandler *TodoHandler) DeleteTodo() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		// 从 URL 参数获取留言 ID
		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return res.Response{
				Msg: commonModel.INVALID_PARAMS_BODY,
				Err: err,
			}
		}

		if err := todoHandler.todoService.DeleteTodo(userid, int64(id)); err != nil {
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

func (todoHandler *TodoHandler) GetTodoList() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 获取当前用户 ID
		userid := ctx.MustGet("userid").(uint)

		todos, err := todoHandler.todoService.GetTodoList(userid)
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: todos,
			Msg:  commonModel.GET_TODO_LIST_SUCCESS,
		}
	})

}
