package service

import (
	"errors"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	"github.com/lin-snow/ech0/internal/model/todo"
	"github.com/lin-snow/ech0/internal/repository/todo"
	commonService "github.com/lin-snow/ech0/internal/service/common"
)

type TodoService struct {
	todoRepository repository.TodoRepositoryInterface
	commonService  commonService.CommonServiceInterface
}

func NewTodoService(todoRepository repository.TodoRepositoryInterface, commonService commonService.CommonServiceInterface) TodoServiceInterface {
	return &TodoService{
		todoRepository: todoRepository,
		commonService:  commonService,
	}
}

func (todoService *TodoService) GetTodoList(userid uint) ([]model.Todo, error) {
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

func (todoService *TodoService) AddTodo(userid uint, todo *model.Todo) error {
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

	if err := todoService.todoRepository.CreateTodo(todo); err != nil {
		return err
	}
	return nil
}

func (todoService *TodoService) UpdateTodo(userid uint, id int64) error {
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

func (todoService *TodoService) DeleteTodo(userid uint, id int64) error {
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
