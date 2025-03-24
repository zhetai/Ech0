package repository

import (
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/models"
)

func GetUserByUsername(username string) (models.User, error) {
	user := models.User{}
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func CreateUser(user *models.User) error {
	err := database.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := database.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByID(id uint) (models.User, error) {
	user := models.User{}
	err := database.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetSysAdmin() (models.User, error) {
	// 获取系统管理员（首个注册的用户）
	user := models.User{}
	err := database.DB.Where("is_admin = ?", true).First(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func UpdateUser(user *models.User) error {
	err := database.DB.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}
