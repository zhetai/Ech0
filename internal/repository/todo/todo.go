package repository

import (
	model "github.com/lin-snow/ech0/internal/model/todo"
	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepositoryInterface {
	return &TodoRepository{
		db: db,
	}
}

// GetTodosByUserID 根据用户ID获取待办事项
func (todoRepository *TodoRepository) GetTodosByUserID(userid uint) ([]model.Todo, error) {
	var todos []model.Todo
	// 查询数据库(按创建时间，最新的在前)
	if err := todoRepository.db.Where("user_id = ?", userid).Order("created_at DESC").Find(&todos).Error; err != nil {
		return nil, err
	}
	// 如果没有找到，返回空切片
	if len(todos) == 0 {
		return []model.Todo{}, nil
	}
	// 返回查询到的 todos
	return todos, nil
}

// CreateTodo 创建一个新的待办事项
func (todoRepository *TodoRepository) CreateTodo(todo *model.Todo) error {
	if err := todoRepository.db.Create(todo).Error; err != nil {
		return err
	}
	return nil
}

// GetTodoByID 根据ID获取待办事项
func (todoRepository *TodoRepository) GetTodoByID(todoID int64) (*model.Todo, error) {
	var todo model.Todo
	// 根据 ID 查找 To do
	if err := todoRepository.db.Where("id = ?", todoID).First(&todo).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

// UpdateTodo 更新待办事项
func (todoRepository *TodoRepository) UpdateTodo(todo *model.Todo) error {
	// 根据 ID 查找 To do 并更新
	if err := todoRepository.db.Model(&model.Todo{}).Where("id = ?", todo.ID).Updates(todo).Error; err != nil {
		return err
	}

	return nil
}

// DeleteTodo 删除待办事项
func (todoRepository *TodoRepository) DeleteTodo(id int64) error {
	// 根据 ID 删除 To do
	if err := todoRepository.db.Delete(&model.Todo{}, id).Error; err != nil {
		return err
	}

	return nil
}
