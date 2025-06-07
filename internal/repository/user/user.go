package repository

import (
	model "github.com/lin-snow/ech0/internal/model/user"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{
		db: db,
	}
}

func (userRepository *UserRepository) GetUserByUsername(username string) (model.User, error) {
	user := model.User{}
	err := userRepository.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (userRepository *UserRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := userRepository.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (userRepository *UserRepository) CreateUser(user *model.User) error {
	err := userRepository.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (userRepository *UserRepository) GetUserByID(id int) (model.User, error) {
	var user model.User
	if err := userRepository.db.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (userRepository *UserRepository) GetSysAdmin() (model.User, error) {
	// 获取系统管理员（首个注册的用户）
	user := model.User{}
	err := userRepository.db.Where("is_admin = ?", true).First(&user).Error
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (userRepository *UserRepository) UpdateUser(user *model.User) error {
	err := userRepository.db.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (userRepository *UserRepository) DeleteUser(id uint) error {
	err := userRepository.db.Delete(&model.User{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
