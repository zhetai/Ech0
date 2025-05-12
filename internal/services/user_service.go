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
	// 检查是否允许注册
	setting, err := GetSetting()
	if err != nil {
		return errors.New(models.GetSettingsFailMessage)
	}

	// 如果第一次注册用户，则不需要检查是否允许注册，否则检查是否允许注册
	if len(users) != 0 && !setting.AllowRegister {
		return errors.New(models.RegisterNotAllowedMessage)
	}

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

func UpdateUser(user models.User, userdto dto.UserInfoDto) error {
	// 检查是否需要更新用户名
	if userdto.Username != "" && userdto.Username != user.Username {
		// 检查用户名是否已存在
		existingUser, _ := repository.GetUserByUsername(userdto.Username)
		if existingUser.ID != 0 {
			return errors.New(models.UsernameAlreadyExistsMessage)
		}
		user.Username = userdto.Username
	}

	// 检查是否需要更新密码
	if userdto.Password != "" && pkg.MD5Encrypt(userdto.Password) != user.Password {
		// 检查密码是否为空
		if userdto.Password == "" {
			return errors.New(models.PasswordCannotBeEmptyMessage)
		}
		// 更新密码
		user.Password = pkg.MD5Encrypt(userdto.Password)
	}

	// 检查是否需要更新头像
	if userdto.Avatar != "" && userdto.Avatar != user.Avatar {
		// 更新头像
		user.Avatar = userdto.Avatar
	}

	// 更新用户信息
	if err := repository.UpdateUser(&user); err != nil {
		return errors.New(models.UpdateUserFailMessage)
	}

	return nil
}

func UpdateUserAdmin(userID uint) error {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	user.IsAdmin = !user.IsAdmin

	if err := repository.UpdateUser(&user); err != nil {
		return err
	}

	return nil
}
