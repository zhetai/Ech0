package service

import (
	"errors"
	"github.com/lin-snow/ech0/internal/transaction"

	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/todo"
	repository "github.com/lin-snow/ech0/internal/repository/todo"
	commonService "github.com/lin-snow/ech0/internal/service/common"
)

type TodoService struct {
	txManager      transaction.TransactionManager       // 事务管理器
	todoRepository repository.TodoRepositoryInterface   // To do数据层接口
	commonService  commonService.CommonServiceInterface // 公共服务接口
}

func NewTodoService(
	tm transaction.TransactionManager,
	todoRepository repository.TodoRepositoryInterface,
	commonService commonService.CommonServiceInterface,
) TodoServiceInterface {
	return &TodoService{
		txManager:      tm,
		todoRepository: todoRepository,
		commonService:  commonService,
	}
}

// GetTodoList 获取当前用户的Todo列表
func (todoService *TodoService) GetTodoList(userid uint) ([]model.Todo, error) {
	// 检查执行操作的用户是否为管理员
	user, err := todoService.commonService.CommonGetUserByUserId(userid)
	if err != nil {
		return nil, err
	}
	if !user.IsAdmin {
		return nil, errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	todos, err := todoService.todoRepository.GetTodosByUserID(userid)
	if err != nil {
		return nil, err
	}

	// 除去已完成的 To do
	for i := len(todos) - 1; i >= 0; i-- {
		if todos[i].Status == uint(model.Done) {
			todos = append(todos[:i], todos[i+1:]...)
		}
	}
	return todos, nil
}

// AddTodo 创建新的Todo
func (todoService *TodoService) AddTodo(userid uint, todo *model.Todo) error {
	// 检查执行操作的用户是否为管理员
	user, err := todoService.commonService.CommonGetUserByUserId(userid)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	todos, err := todoService.todoRepository.GetTodosByUserID(userid)
	if err != nil {
		return err
	}
	// 除去已完成的 To do
	for i := len(todos) - 1; i >= 0; i-- {
		if todos[i].Status == uint(model.Done) {
			todos = append(todos[:i], todos[i+1:]...)
		}
	}
	if len(todos) >= model.MaxTodoCount {
		return errors.New(commonModel.TODO_EXCEED_LIMIT)
	}

	// 设置TO DO
	todo.UserID = userid
	todo.Username = user.Username
	todo.Status = uint(model.NotDone)

	// 创建 To do
	if err := todoService.todoRepository.CreateTodo(todo); err != nil {
		return err
	}
	return nil
}

// UpdateTodo 更新指定ID的Todo
func (todoService *TodoService) UpdateTodo(userid uint, id int64) error {
	// 检查执行操作的用户是否为管理员
	user, err := todoService.commonService.CommonGetUserByUserId(userid)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 获取 To do
	theTodo, err := todoService.todoRepository.GetTodoByID(id)
	if err != nil {
		return err
	}

	// 检查该 To do 是否属于当前用户
	if theTodo.UserID != userid {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 设置To do的状态
	if theTodo.Status == uint(model.NotDone) {
		theTodo.Status = uint(model.Done)
	} else {
		theTodo.Status = uint(model.NotDone)
	}

	if err := todoService.todoRepository.UpdateTodo(theTodo); err != nil {
		return err
	}

	return nil
}

// DeleteTodo 删除指定ID的Todo
func (todoService *TodoService) DeleteTodo(userid uint, id int64) error {
	// 检查执行操作的用户是否为管理员
	user, err := todoService.commonService.CommonGetUserByUserId(userid)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 获取 To do
	theTodo, err := todoService.todoRepository.GetTodoByID(id)
	if err != nil {
		return err
	}

	// 检查该 To do 是否属于当前用户
	if theTodo.UserID != userid {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	if err := todoService.todoRepository.DeleteTodo(id); err != nil {
		return err
	}

	return nil
}
