package services

import (
	"errors"

	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/repository"
	"github.com/lin-snow/ech0/pkg"
)

func Register(userdto dto.RegisterDto) error {
	if userdto.Username == "" || userdto.Password == "" {
		return errors.New(models.UsernameOrPasswordCannotBeEmptyMessage)
	}

	// 将密码进行md5加密
	userdto.Password = pkg.MD5Encrypt(userdto.Password)

	newuser := models.User{
		Username: userdto.Username,
		Password: userdto.Password,
		IsAdmin:  false,
	}

	// 检查用户名是否已存在
	user, err := repository.GetUserByUsername(userdto.Username)
	if err == nil && user.ID != 0 {
		return errors.New(models.UsernameAlreadyExistsMessage)
	}

	// 检查是否该系统第一次注册用户
	users, err := repository.GetAllUsers()
	if err != nil {
		return errors.New(models.GetAllUsersFailMessage)
	}
	if len(users) == 0 {
		newuser.IsAdmin = true
	}

	// 创建新用户
	if err := repository.CreateUser(&newuser); err != nil {
		return errors.New(models.CreateUserFailMessage)
	}

	return nil
}

func Login(userdto dto.LoginDto) (string, error) {
	if userdto.Username == "" || userdto.Password == "" {
		return "", errors.New(models.UsernameOrPasswordCannotBeEmptyMessage)
	}

	// 将密码进行md5加密
	userdto.Password = pkg.MD5Encrypt(userdto.Password)

	user, err := repository.GetUserByUsername(userdto.Username)
	if err != nil {
		return "", errors.New(models.UserNotFoundMessage)
	}

	if user.Password != userdto.Password {
		return "", errors.New(models.PasswordIncorrectMessage)
	}

	// 生成Token
	token, err := pkg.GenerateToken(pkg.CreateClaims(user))
	if err != nil {
		return "", errors.New(models.GenerateTokenFailMessage)
	}

	return token, nil
}

func GetStatus() (models.Status, error) {
	status := models.Status{}

	user, err := repository.GetAdmin()
	if err != nil {
		return status, errors.New(models.UserNotFoundMessage)
	}

	messages, err := repository.GetAllMessages()
	if err != nil {
		return status, errors.New(models.GetAllMessagesFailMessage)
	}

	status.UserID = user.ID
	status.Username = user.Username
	status.IsAdmin = user.IsAdmin
	status.TotalMessages = len(messages)

	return status, nil
}

func GetUserByID(userID uint) (models.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return models.User{}, errors.New(models.UserNotFoundMessage)
	}
	return user, nil
}

func IsUserAdmin(userID uint) (bool, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return false, errors.New(models.UserNotFoundMessage)
	}
	return user.IsAdmin, nil
}
