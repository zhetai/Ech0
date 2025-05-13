package repository

import (
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/models"
)

func GetTodosByUserID(userID uint) ([]models.Todo, error) {
	var todos []models.Todo
	// 查询数据库(按创建时间，最新的在前)
	if err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&todos).Error; err != nil {
		return nil, err
	}
	// 如果没有找到，返回空切片
	if len(todos) == 0 {
		return []models.Todo{}, nil
	}
	// 返回查询到的 todos
	return todos, nil
}

func CreateTodo(todo *models.Todo) error {
	if err := database.DB.Create(todo).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTodo(todo *models.Todo) error {
	// 根据 ID 查找 Todo 并更新
	if err := database.DB.Model(&models.Todo{}).Where("id = ?", todo.ID).Updates(todo).Error; err != nil {
		return err
	}

	return nil
}

func DeleteTodo(id uint) error {
	// 根据 ID 删除 Todo
	if err := database.DB.Delete(&models.Todo{}, id).Error; err != nil {
		return err
	}

	return nil
}
