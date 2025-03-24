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
	// 获取系统管理员信息
	sysuser, err := repository.GetSysAdmin()
	if err != nil {
		return models.Status{}, errors.New(models.UserNotFoundMessage)
	}

	// 获取所有用户状态信息
	var users []models.UserStatus
	allusers, err := repository.GetAllUsers()
	if err != nil {
		return models.Status{}, errors.New(models.GetAllUsersFailMessage)
	}
	for _, user := range allusers {
		users = append(users, models.UserStatus{
			UserID:   user.ID,
			UserName: user.Username,
			IsAdmin:  user.IsAdmin,
		})
	}

	status := models.Status{}

	messages, err := repository.GetAllMessages()
	if err != nil {
		return status, errors.New(models.GetAllMessagesFailMessage)
	}

	status.SysAdminID = sysuser.ID
	status.Username = sysuser.Username
	status.Users = users
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

func UpdateUser(user models.User, userdto dto.UserInfoDto) error {
	// 检查更新的用户名是否与之前一样，如果一样则不更新
	if user.Username == userdto.Username {
		return nil
	}

	// 检查用户名是否为空
	if userdto.Username == "" {
		return errors.New(models.UsernameCannotBeEmptyMessage)
	}

	// 更新用户名
	user.Username = userdto.Username
	if err := repository.UpdateUser(&user); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func ChangePassword(user models.User, userdto dto.UserInfoDto) error {
	// 检查待更新的密码是否与之前一样，如果一样则不更新
	if user.Password == pkg.MD5Encrypt(userdto.Password) {
		return errors.New(models.PasswordCannotBeSameAsBeforeMessage)
	}

	// 检查密码是否为空
	if userdto.Password == "" {
		return errors.New(models.PasswordCannotBeEmptyMessage)
	}

	// 更新密码
	user.Password = pkg.MD5Encrypt(userdto.Password)
	if err := repository.UpdateUser(&user); err != nil {
		return errors.New(err.Error())
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
