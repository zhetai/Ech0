package models

// 成功相关
const (
	SuccessMessage               = "请求成功"
	GetAllMessagesSuccess        = "获取留言成功"
	GetMessageByIDSuccess        = "获取留言成功"
	PostMessageSuccess           = "发布留言成功"
	CreateSuccessMessage         = "创建成功"
	QuerySuccessMessage          = "查询成功"
	DeleteSuccessMessage         = "删除成功"
	ServerLaunchSuccessMessage   = "服务器启动成功"
	GetMessagesByPageSuccess     = "分页获取留言成功"
	RegisterSuccessMessage       = "注册成功"
	LoginSuccessMessage          = "登录成功"
	UpdateUserSuccessMessage     = "更新用户成功"
	ChangePasswordSuccessMessage = "修改密码成功"
	GetSettingsSuccessMessage    = "获取设置成功"
	UpdateSettingsSuccessMessage = "更新设置成功"
	GetHeatMapSuccessMessage     = "获取热力图成功"
	CreateTodoSuccessMessage     = "创建待办事项成功"
	GetTodosSuccessMessage       = "获取待办事项成功"
	UpdateTodoSuccessMessage     = "更新待办事项成功"
	DeleteTodoSuccessMessage     = "删除待办事项成功"
	GetConnectSuccessMessage     = "连接成功"
	AddConnectSuccessMessage     = "添加连接成功"
	GetConnectsSuccessMessage    = "获取连接列表成功"
	DeleteConnectSuccessMessage  = "连接已取消"
	DeleteUserSuccessMessage     = "删除用户成功"
	DeleteImageSuccessMessage    = "删除图片成功"
	PleaseSignUpFirstMessage     = "请先注册用户"
	GetPlayMusicSuccessMessage   = "获取音乐成功"
	DeleteAudioSuccessMessage    = "删除音频成功"
)

// 失败相关
const (
	ErrorMessage                           = "请求失败"
	CreateFailMessage                      = "创建失败"
	GetAllMessagesFailMessage              = "获取留言失败"
	GetMessageByIDFailMessage              = "获取留言失败"
	MessageNotFoundMessage                 = "留言未找到"
	PostMessageFailMessage                 = "发布留言失败"
	DeleteFailMessage                      = "删除失败"
	QueryFailMessage                       = "查询失败"
	InvalidIDMessage                       = "无效的ID"
	NotFoundMessage                        = "资源未找到"
	ValidationErrorMessage                 = "数据验证失败"
	InternalServerErrorMessage             = "内部服务器错误"
	ServerLaunchErrorMessage               = "服务器启动失败"
	LoadConfigErrorMessage                 = "加载配置失败"
	ParseConfigErrorMessage                = "解析配置失败"
	InvalidRequestMessage                  = "无效的请求"
	InvalidRequestBodyMessage              = "无效的请求体"
	ImageUploadErrorMessage                = "图片上传失败"
	NotUploadImageErrorMessage             = "未上传图片"
	NotSupportedImageTypeMessage           = "不支持的图片类型"
	ImageSizeLimitErrorMessage             = "图片大小超过限制"
	InvalidPageParametersMessage           = "无效的分页参数"
	InvalidPageSizeMessage                 = "无效的页码或页大小"
	GetMessagesByPageFailMessage           = "分页获取留言失败"
	TokenNotFoundMessage                   = "未找到令牌,请点击右上角登录"
	TokenInvalidMessage                    = "令牌无效，请点击右上角登录"
	UsernameOrPasswordCannotBeEmptyMessage = "用户名或密码不能为空"
	UsernameCannotBeEmptyMessage           = "用户名不能为空"
	PasswordCannotBeEmptyMessage           = "密码不能为空"
	UsernameAlreadyExistsMessage           = "用户名已存在"
	CreateUserFailMessage                  = "创建用户失败"
	UserNotFoundMessage                    = "用户未找到"
	PasswordIncorrectMessage               = "密码错误"
	GenerateTokenFailMessage               = "生成令牌失败"
	GetAllUsersFailMessage                 = "获取所有用户失败"
	GetStatusFailMessage                   = "获取状态失败"
	GetStatusSuccessMessage                = "获取状态成功"
	CannotBeEmptyMessage                   = "内容不能为空"
	NoPermissionMessage                    = "没有权限,请联系管理员"
	GenerateRSSFailMessage                 = "生成 RSS 失败"
	PasswordCannotBeSameAsBeforeMessage    = "新密码不能与旧密码相同"
	GetSettingsFailMessage                 = "获取设置失败"
	RegisterNotAllowedMessage              = "当前系统不允许注册新用户"
	GetHeatMapFailMessage                  = "获取热力图失败"
	UpdateUserFailMessage                  = "更新用户失败"
	CreateTodoFailMessage                  = "创建待办事项失败"
	GetTodosFailMessage                    = "获取待办事项失败"
	UpdateTodoFailMessage                  = "更新待办事项失败"
	DeleteTodoFailMessage                  = "删除待办事项失败"
	TodoNotFoundMessage                    = "待办事项未找到"
	MaxTodoCountMessage                    = "待办事项数量已达上限"
	GetConnectFailMessage                  = "获取 Connect 信息失败"
	ConnectAlreadyExistsMessage            = "请不要重复添加"
	GetConnectsFailMessage                 = "获取 Connect 列表失败"
	ConnectURLIsEmptyMessage               = "连接地址不能为空"
	UserCountExceedsLimitMessage           = "用户数量超过限制"
	NoSysPermissionMessage                 = "请使用系统管理员权限"
	ImageNotFoundMessage                   = "图片未找到"
	NotUploadFileErrorMessage              = "未上传文件"
	NotSupportedFileTypeErrorMessage       = "不支持的文件类型"
	AudioUploadErrorMessage                = "音频上传失败"
	AudioSizeLimitErrorMessage             = "音频大小超过限制"
	AudioNotFoundMessage                   = "音频未找到"
)

// 数据库相关
const (
	DatabaseTypeMessage       = "数据库类型错误"
	DatabaseConnectionError   = "数据库连接失败"
	DatabaseMigrationError    = "数据库迁移失败"
	DatabaseInitErrorMessage  = "数据库初始化失败"
	DatabaseErrorMessage      = "数据库操作失败"
	DatabaseConnectionSuccess = "数据库连接成功"
	DuplicateEntryMessage     = "重复的条目"
	RecordNotFoundMessage     = "记录未找到"
)

// 键值对相关
const (
	SystemSettingsKey = "system_settings" // 系统设置的键
	ConnectKey        = "connect"         // Connect 信息的键
)

// Ech0 相关
const (
	Extension_MUSIC      = "MUSIC"
	Extension_VIDEO      = "VIDEO"
	Extension_GITHUBPROJ = "GITHUBPROJ"
)

// User 相关
const (
	MaxUserCount = 4 // 最多注册用户数
)

// Todo 相关
const (
	Done         = 1 // 待办事项已完成状态
	NotDone      = 0 // 待办事项状态
	MaxTodoCount = 3 // 最大待办事项数量
)

// File 相关
type FileType string

const (
	ImageType FileType = "image" // 图片类型
	AudioType FileType = "audio" // 音频类型
)

// 验证相关
const (
	InvalidDataMessage  = "无效数据"
	MissingFieldMessage = "缺少必填字段"
)

// 其他

const (
	InitInstallCode = 666
)

const (
	Version = "1.2.7" // 当前版本号
)

const (
	GreetingBanner = `
███████╗     ██████╗    ██╗  ██╗     ██████╗ 
██╔════╝    ██╔════╝    ██║  ██║    ██╔═████╗
█████╗      ██║         ███████║    ██║██╔██║
██╔══╝      ██║         ██╔══██║    ████╔╝██║
███████╗    ╚██████╗    ██║  ██║    ╚██████╔╝
╚══════╝     ╚═════╝    ╚═╝  ╚═╝     ╚═════╝ 
                                             
`
)
