// Package service 提供用户相关的业务逻辑服务
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/lin-snow/ech0/internal/transaction"

	authModel "github.com/lin-snow/ech0/internal/model/auth"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	settingModel "github.com/lin-snow/ech0/internal/model/setting"
	model "github.com/lin-snow/ech0/internal/model/user"
	repository "github.com/lin-snow/ech0/internal/repository/user"
	settingService "github.com/lin-snow/ech0/internal/service/setting"
	cryptoUtil "github.com/lin-snow/ech0/internal/util/crypto"
	jwtUtil "github.com/lin-snow/ech0/internal/util/jwt"
)

// UserService 用户服务结构体，提供用户相关的业务逻辑处理
type UserService struct {
	txManager      transaction.TransactionManager         // 事务管理器
	userRepository repository.UserRepositoryInterface     // 用户数据层接口
	settingService settingService.SettingServiceInterface // 系统设置数据层接口
}

// NewUserService 创建并返回新的用户服务实例
//
// 参数:
//   - userRepository: 用户数据层接口实现
//   - settingService: 系统设置数据层接口实现
//
// 返回:
//   - UserServiceInterface: 用户服务接口实现
func NewUserService(
	tm transaction.TransactionManager,
	userRepository repository.UserRepositoryInterface,
	settingService settingService.SettingServiceInterface,
) UserServiceInterface {
	return &UserService{
		txManager:      tm,
		userRepository: userRepository,
		settingService: settingService,
	}
}

// Login 用户登录验证
// 验证用户名和密码，成功后生成JWT token
//
// 参数:
//   - loginDto: 登录数据传输对象，包含用户名和密码
//
// 返回:
//   - string: 生成的JWT token
//   - error: 登录过程中的错误信息
func (userService *UserService) Login(loginDto *authModel.LoginDto) (string, error) {
	// 合法性校验
	if loginDto.Username == "" || loginDto.Password == "" {
		return "", errors.New(commonModel.USERNAME_OR_PASSWORD_NOT_BE_EMPTY)
	}

	// 将密码进行 MD5 加密
	loginDto.Password = cryptoUtil.MD5Encrypt(loginDto.Password)

	// 检查用户是否存在
	user, err := userService.userRepository.GetUserByUsername(loginDto.Username)
	if err != nil {
		return "", errors.New(commonModel.USER_NOTFOUND)
	}

	// 进行密码验证,查看外界传入的密码是否与数据库一致
	if user.Password != loginDto.Password {
		return "", errors.New(commonModel.PASSWORD_INCORRECT)
	}

	// 生成 Token
	token, err := jwtUtil.GenerateToken(jwtUtil.CreateClaims(user))
	if err != nil {
		return "", err
	}

	return token, nil
}

// Register 用户注册
// 注册新用户，包括用户数量限制检查、注册权限检查等
// 第一个注册的用户自动设置为系统管理员
//
// 参数:
//   - registerDto: 注册数据传输对象，包含用户名和密码
//
// 返回:
//   - error: 注册过程中的错误信息
func (userService *UserService) Register(registerDto *authModel.RegisterDto) error {
	return userService.txManager.Run(func(ctx context.Context) error {
		// 检查用户数量是否超过限制
		users, err := userService.userRepository.GetAllUsers()
		if err != nil {
			return err
		}
		if len(users) > authModel.MAX_USER_COUNT {
			return errors.New(commonModel.USER_COUNT_EXCEED_LIMIT)
		}

		// 将密码进行 MD5 加密
		registerDto.Password = cryptoUtil.MD5Encrypt(registerDto.Password)

		newUser := model.User{
			Username: registerDto.Username,
			Password: registerDto.Password,
			IsAdmin:  false,
		}

		// 检查用户是否已经存在
		user, err := userService.userRepository.GetUserByUsername(newUser.Username)
		if err == nil && user.ID != model.USER_NOT_EXISTS_ID {
			return errors.New(commonModel.USERNAME_HAS_EXISTS)
		}

		// 检查是否该系统第一次注册用户
		if len(users) == 0 {
			// 第一个注册的用户为系统管理员
			newUser.IsAdmin = true
		}

		// 检查是否开放注册
		var setting settingModel.SystemSetting
		if err := userService.settingService.GetSetting(&setting); err != nil {
			return err
		}
		if len(users) != 0 && !setting.AllowRegister {
			return errors.New(commonModel.USER_REGISTER_NOT_ALLOW)
		}

		if err := userService.userRepository.CreateUser(ctx, &newUser); err != nil {
			return err
		}

		return nil
	})
}

// UpdateUser 更新用户信息
// 只有管理员可以更新用户信息，支持更新用户名、密码和头像
//
// 参数:
//   - userid: 执行更新操作的用户ID（必须为管理员）
//   - userdto: 用户信息数据传输对象，包含要更新的用户信息
//
// 返回:
//   - error: 更新过程中的错误信息
func (userService *UserService) UpdateUser(userid uint, userdto model.UserInfoDto) error {
	return userService.txManager.Run(func(ctx context.Context) error {
		// 检查执行操作的用户是否为管理员
		user, err := userService.userRepository.GetUserByID(int(userid))
		if err != nil {
			return err
		}
		if !user.IsAdmin {
			return errors.New(commonModel.NO_PERMISSION_DENIED)
		}

		// 检查是否需要更新用户名
		if userdto.Username != "" && userdto.Username != user.Username {
			// 检查用户名是否已存在
			existingUser, _ := userService.userRepository.GetUserByUsername(userdto.Username)
			if existingUser.ID != model.USER_NOT_EXISTS_ID {
				return errors.New(commonModel.USERNAME_ALREADY_EXISTS)
			}
			user.Username = userdto.Username
		}

		// 检查是否需要更新密码
		if userdto.Password != "" && cryptoUtil.MD5Encrypt(userdto.Password) != user.Password {
			// 检查密码是否为空
			if userdto.Password == "" {
				return errors.New(commonModel.USERNAME_OR_PASSWORD_NOT_BE_EMPTY)
			}
			// 更新密码
			user.Password = cryptoUtil.MD5Encrypt(userdto.Password)
		}

		// 检查是否需要更新头像
		if userdto.Avatar != "" && userdto.Avatar != user.Avatar {
			// 更新头像
			user.Avatar = userdto.Avatar
		}
		// 更新用户信息
		if err := userService.userRepository.UpdateUser(ctx, &user); err != nil {
			return err
		}

		return nil
	})
}

// UpdateUserAdmin 更新用户的管理员权限
// 只有系统管理员、管理员可以修改其他用户的管理员权限，不能修改自己和系统管理员的权限
//
// 参数:
//   - userid: 执行操作的用户ID（必须为管理员）
//   - id: 要修改权限的用户ID
//
// 返回:
//   - error: 更新过程中的错误信息
func (userService *UserService) UpdateUserAdmin(userid uint, id uint) error {
	return userService.txManager.Run(func(ctx context.Context) error {
		// 检查执行操作的用户是否为管理员
		user, err := userService.userRepository.GetUserByID(int(userid))
		if err != nil {
			return err
		}
		if !user.IsAdmin {
			return errors.New(commonModel.NO_PERMISSION_DENIED)
		}

		// 检查要修改权限的用户是否存在
		user, err = userService.userRepository.GetUserByID(int(id))
		if err != nil {
			return err
		}

		// 检查系统管理员信息
		sysadmin, err := userService.GetSysAdmin()
		if err != nil {
			return err
		}

		// 检查是否尝试修改自己或系统管理员的权限
		if userid == user.ID || id == sysadmin.ID {
			return errors.New(commonModel.INVALID_PARAMS_BODY)
		}

		user.IsAdmin = !user.IsAdmin

		// 更新用户信息
		if err := userService.userRepository.UpdateUser(ctx, &user); err != nil {
			return err
		}

		return nil
	})
}

// GetAllUsers 获取所有用户列表
// 返回除系统管理员外的所有用户，并移除密码信息
//
// 返回:
//   - []model.User: 用户列表（不包含密码信息）
//   - error: 获取过程中的错误信息
func (userService *UserService) GetAllUsers() ([]model.User, error) {
	allures, err := userService.userRepository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	sysadmin, err := userService.GetSysAdmin()
	if err != nil {
		return nil, err
	}

	// 处理用户信息(去掉管理员用户)
	for i := range allures {
		if allures[i].ID == sysadmin.ID {
			allures = append(allures[:i], allures[i+1:]...)
			break
		}
	}

	// 处理用户信息(去掉密码)
	for i := range allures {
		allures[i].Password = ""
	}

	return allures, nil
}

// GetSysAdmin 获取系统管理员信息
//
// 返回:
//   - model.User: 系统管理员用户信息
//   - error: 获取过程中的错误信息
func (userService *UserService) GetSysAdmin() (model.User, error) {
	sysadmin, err := userService.userRepository.GetSysAdmin()
	if err != nil {
		return model.User{}, err
	}

	return sysadmin, nil
}

// DeleteUser 删除用户
// 只有管理员可以删除用户，不能删除自己和系统管理员
//
// 参数:
//   - userid: 执行删除操作的用户ID（必须为管理员）
//   - id: 要删除的用户ID
//
// 返回:
//   - error: 删除过程中的错误信息
func (userService *UserService) DeleteUser(userid, id uint) error {
	return userService.txManager.Run(func(ctx context.Context) error {
		// 检查执行操作的用户是否为管理员
		user, err := userService.userRepository.GetUserByID(int(userid))
		if err != nil {
			return err
		}
		if !user.IsAdmin {
			return errors.New(commonModel.NO_PERMISSION_DENIED)
		}

		// 检查要删除的用户是否存在
		user, err = userService.userRepository.GetUserByID(int(id))
		if err != nil {
			return err
		}

		sysadmin, err := userService.GetSysAdmin()
		if err != nil {
			return err
		}

		if userid == user.ID || id == sysadmin.ID {
			return errors.New(commonModel.INVALID_PARAMS_BODY)
		}

		if err := userService.userRepository.DeleteUser(ctx, id); err != nil {
			return err
		}

		return nil
	})

}

// GetUserByID 根据用户ID获取用户信息
//
// 参数:
//   - userId: 用户ID
//
// 返回:
//   - model.User: 用户信息
//   - error: 获取过程中的错误信息
func (userService *UserService) GetUserByID(userId int) (model.User, error) {
	return userService.userRepository.GetUserByID(userId)
}

// BindGitHub 绑定 GitHub 账号
//
// 参数:
//   - userID: 当前用户 ID
//
// 返回:
//   - error: 绑定过程中的错误信息
func (userService *UserService) BindGitHub(userID uint, redirect_URI string) (string, error) {
	// 检查当前用户是否存在
	user, err := userService.userRepository.GetUserByID(int(userID))
	if err != nil {
		return "", err
	}

	// 检查用户是否为管理员
	if !user.IsAdmin {
		return "", errors.New(commonModel.NO_PERMISSION_BINDING)
	}

	var setting settingModel.OAuth2Setting
	if err := userService.settingService.GetOAuth2Setting(0, &setting, true); err != nil {
		return "", err
	}

	if !setting.Enable {
		return "", errors.New(commonModel.OAUTH2_NOT_ENABLED)
	}

	if setting.ClientID == "" || setting.RedirectURI == "" || setting.AuthURL == "" || setting.TokenURL == "" || setting.UserInfoURL == "" || setting.ClientSecret == "" {
		return "", errors.New(commonModel.OAUTH2_NOT_CONFIGURED)
	}

	// 生成附带用户 ID 的 state 参数
	state, err := jwtUtil.GenerateOAuthState(string(authModel.OAuth2ActionBind), userID, redirect_URI, string(commonModel.OAuth2GITHUB))
	if err != nil {
		return "", err
	}

	// 拼接 scope 参数
	scope := ""
	if len(setting.Scopes) > 0 {
		scope = strings.Join(setting.Scopes, " ")
	}

	// 拼接 OAuth2 登录 URL
	bingURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		setting.AuthURL,
		url.QueryEscape(setting.ClientID),
		url.QueryEscape(setting.RedirectURI),
		url.QueryEscape(scope),
		url.QueryEscape(state),
	)

	return bingURL, nil
}

// GetGitHubLoginURL 获取 GitHub 登录 URL
//
// 返回:
//   - string: GitHub 登录 URL
//   - error: 获取过程中的错误信息
func (userService *UserService) GetGitHubLoginURL(redirect_URI string) (string, error) {
	var setting settingModel.OAuth2Setting
	if err := userService.settingService.GetOAuth2Setting(0, &setting, true); err != nil {
		return "", err
	}

	if !setting.Enable {
		return "", errors.New(commonModel.OAUTH2_NOT_ENABLED)
	}

	if setting.ClientID == "" || setting.RedirectURI == "" || setting.AuthURL == "" || setting.TokenURL == "" || setting.UserInfoURL == "" || setting.ClientSecret == "" {
		return "", errors.New(commonModel.OAUTH2_NOT_CONFIGURED)
	}

	// 生成随机的 state 参数，防止 CSRF 攻击
	state, err := jwtUtil.GenerateOAuthState(string(authModel.OAuth2ActionLogin), authModel.NO_USER_LOGINED, redirect_URI, string(commonModel.OAuth2GITHUB))
	if err != nil {
		return "", err
	}

	scope := ""
	if len(setting.Scopes) > 0 {
		scope = strings.Join(setting.Scopes, " ")
	}

	loginURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		setting.AuthURL,
		url.QueryEscape(setting.ClientID),
		url.QueryEscape(setting.RedirectURI),
		url.QueryEscape(scope),
		url.QueryEscape(state),
	)

	return loginURL, nil
}

// HandleGitHubCallback 处理 GitHub OAuth2 回调
//
// 参数:
//   - code: GitHub 回调返回的授权码
//   - state: GitHub 回调返回的状态参数
//
// 返回:
//   - string: 重定向的前端 URL，包含登录结果信息
func (userService *UserService) HandleGitHubCallback(code string, state string) string {
	// 获取 OAuth2 设置
	var setting settingModel.OAuth2Setting
	if err := userService.settingService.GetOAuth2Setting(0, &setting, true); err != nil {
		return ""
	}

	if !setting.Enable {
		return ""
	}

	if setting.ClientID == "" || setting.RedirectURI == "" || setting.AuthURL == "" || setting.TokenURL == "" || setting.UserInfoURL == "" || setting.ClientSecret == "" {
		return ""
	}

	// 提取 state 信息
	oauthState, err := jwtUtil.ParseOAuthState(state)
	if err != nil {
		return ""
	}

	// 2. 用 code 换取 access_token
	tokenResp, err := exchangeCodeForToken(&setting, code)
	if err != nil {
		fmt.Println("Error exchanging code for token:", err)
		return ""
	}

	// 3. 获取 GitHub 用户信息
	githubUser, err := fetchGitHubUserInfo(&setting, tokenResp.AccessToken)
	if err != nil {
		fmt.Println("Error fetching GitHub user info:", err)
		return ""
	}

	// 处理不同的 OAuth2 操作
	switch oauthState.Action {
	case string(authModel.OAuth2ActionLogin):
		// 处理登录操作
		if oauthState.UserID != authModel.NO_USER_LOGINED {
			// 非法的登录请求，state 中包含用户 ID
			return ""
		}

		// 根据 GitHub ID 查找用户
		user, err := userService.userRepository.GetUserByOAuthID(context.Background(), string(commonModel.OAuth2GITHUB), fmt.Sprint(githubUser.ID))
		if err != nil {
			fmt.Println("Error fetching user by OAuth ID:", err)
			return ""
		}

		// 根据用户信息生成 JWT token
		token, err := jwtUtil.GenerateToken(jwtUtil.CreateClaims(user))
		if err != nil {
			fmt.Println("Error generating token:", err)
			return ""
		}

		// 构造重定向 URL，包含 token 信息
		redirectURL, err := url.Parse(oauthState.Redirect)
		if err != nil {
			return ""
		}
		query := redirectURL.Query()
		query.Set("token", token)
		redirectURL.RawQuery = query.Encode()

		// 返回重定向 URL
		return redirectURL.String()

	case string(authModel.OAuth2ActionBind):
		// 处理绑定操作
		if oauthState.UserID == authModel.NO_USER_LOGINED {
			// 用户未登录，无法绑定
			return ""
		}

		// 绑定 GitHub 账号
		userService.txManager.Run(func(ctx context.Context) error {
			return userService.userRepository.BindOAuth(ctx, oauthState.UserID, oauthState.Provider, fmt.Sprint(githubUser.ID))
		})

		// 返回绑定成功的前端 URL
		return oauthState.Redirect + "?bind=success"

	default:
		// 未知操作
		return ""
	}
}

// 用 code 换取 access_token
func exchangeCodeForToken(setting *settingModel.OAuth2Setting, code string) (*authModel.GitHubTokenResponse, error) {
	data := map[string]string{
		"client_id":     setting.ClientID,
		"client_secret": setting.ClientSecret,
		"code":          code,
		"redirect_uri":  setting.RedirectURI,
	}
	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", setting.TokenURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, errors.New("GitHub token 响应错误: " + string(body))
	}

	var tokenResp authModel.GitHubTokenResponse
	_ = json.Unmarshal(body, &tokenResp)
	return &tokenResp, nil
}

// 获取 GitHub 用户信息
func fetchGitHubUserInfo(setting *settingModel.OAuth2Setting, accessToken string) (*authModel.GitHubUser, error) {
	req, _ := http.NewRequest("GET", setting.UserInfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, errors.New("GitHub 用户信息请求失败: " + string(body))
	}

	var user authModel.GitHubUser
	_ = json.Unmarshal(body, &user)
	return &user, nil
}

// GetOAuthInfo 获取 OAuth2 信息
func (userService *UserService) GetOAuthInfo(userId uint) (model.OAuthInfoDto, error) {
	var oauthInfo model.OAuthInfoDto

	// 检查当前用户是否存在
	user, err := userService.userRepository.GetUserByID(int(userId))
	if err != nil {
		return oauthInfo, err
	}

	// 检查用户是否为管理员
	if !user.IsAdmin {
		return oauthInfo, errors.New(commonModel.NO_PERMISSION_BINDING)
	}

	oauthInfoBinding, err := userService.userRepository.GetOAuthInfo(userId)
	if err != nil {
		return oauthInfo, err
	}

	oauthInfo = model.OAuthInfoDto{
		Provider: oauthInfoBinding.Provider,
		UserID:   oauthInfoBinding.UserID,
		OAuthID:  oauthInfoBinding.OAuthID,
	}

	return oauthInfo, nil
}
