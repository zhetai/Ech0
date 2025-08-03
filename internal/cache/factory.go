package cache

import (
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	echoModel "github.com/lin-snow/ech0/internal/model/echo"
	userModel "github.com/lin-snow/ech0/internal/model/user"
)

// CacheFactory 是一个工厂类，用于创建和管理不同类型的缓存实例
type CacheFactory struct {
	userCache ICache[string, *userModel.User]
	echoCache ICache[string, commonModel.PageQueryResult[[]echoModel.Echo]]
}

// NewCacheFactory 创建一个新的 CacheFactory 实例，并初始化所需的缓存
func NewCacheFactory() *CacheFactory {
	userCache, err := NewCache[string, *userModel.User]()
	if err != nil {
		panic(err)
	}

	echoCache, err := NewCache[string, commonModel.PageQueryResult[[]echoModel.Echo]]()
	if err != nil {
		panic(err)
	}

	return &CacheFactory{
		userCache: userCache,
		echoCache: echoCache,
	}
}

// UserCache 返回用户缓存实例
func (f *CacheFactory) UserCache() ICache[string, *userModel.User] {
	return f.userCache
}

// EchoCache 返回 Echo 缓存实例
func (f *CacheFactory) EchoCache() ICache[string, commonModel.PageQueryResult[[]echoModel.Echo]] {
	return f.echoCache
}
