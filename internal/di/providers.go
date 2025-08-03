package di

import (
	"github.com/lin-snow/ech0/internal/cache"
	backupHandler "github.com/lin-snow/ech0/internal/handler/backup"
	commonHandler "github.com/lin-snow/ech0/internal/handler/common"
	connectHandler "github.com/lin-snow/ech0/internal/handler/connect"
	echoHandler "github.com/lin-snow/ech0/internal/handler/echo"
	settingHandler "github.com/lin-snow/ech0/internal/handler/setting"
	todoHandler "github.com/lin-snow/ech0/internal/handler/todo"
	userHandler "github.com/lin-snow/ech0/internal/handler/user"
	webHandler "github.com/lin-snow/ech0/internal/handler/web"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	echoModel "github.com/lin-snow/ech0/internal/model/echo"
	userModel "github.com/lin-snow/ech0/internal/model/user"
	"github.com/lin-snow/ech0/internal/transaction"
)

// Handlers 聚合各个模块的Handler
type Handlers struct {
	WebHandler     *webHandler.WebHandler
	UserHandler    *userHandler.UserHandler
	EchoHandler    *echoHandler.EchoHandler
	CommonHandler  *commonHandler.CommonHandler
	SettingHandler *settingHandler.SettingHandler
	TodoHandler    *todoHandler.TodoHandler
	ConnectHandler *connectHandler.ConnectHandler
	BackupHandler  *backupHandler.BackupHandler
}

// NewHandlers 创建Handlers实例
func NewHandlers(
	webHandler *webHandler.WebHandler,
	userHandler *userHandler.UserHandler,
	echoHandler *echoHandler.EchoHandler,
	commonHandler *commonHandler.CommonHandler,
	settingHandler *settingHandler.SettingHandler,
	todoHandler *todoHandler.TodoHandler,
	connectHandler *connectHandler.ConnectHandler,
	backupHandler *backupHandler.BackupHandler,
) *Handlers {
	return &Handlers{
		WebHandler:     webHandler,
		UserHandler:    userHandler,
		EchoHandler:    echoHandler,
		CommonHandler:  commonHandler,
		SettingHandler: settingHandler,
		TodoHandler:    todoHandler,
		ConnectHandler: connectHandler,
		BackupHandler:  backupHandler,
	}
}

// ProvideUserCache 提供 User 缓存实例给 wire 注入
func ProvideUserCache(factory *cache.CacheFactory) cache.ICache[string, *userModel.User] {
	return factory.UserCache()
}

// ProvideEchoCache 提供 Echo 缓存实例给 wire 注入
func ProvideEchoCache(factory *cache.CacheFactory) cache.ICache[string, commonModel.PageQueryResult[[]echoModel.Echo]] {
	return factory.EchoCache()
}

// ProvideTransactionManager 提供事务管理器实例给 wire 注入
func ProvideTransactionManager(factory *transaction.TransactionManagerFactory) transaction.TransactionManager {
	return factory.TransactionManager()
}
