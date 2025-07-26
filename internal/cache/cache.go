package cache

import (
	"github.com/dgraph-io/ristretto/v2"
	userModel "github.com/lin-snow/ech0/internal/model/user"
)

type ICache[K ristretto.Key, V any] interface {
	Set(key K, value V, cost int64) bool
	Get(key K) (V, error)
	Delete(key K)
	GetOrSet(key K, cost int64, fn func() (V, error)) (V, error)
}

func NewCache[K ristretto.Key, V any]() (ICache[K, V], error) {
	return NewRistrettoCache[K, V](1000000, 1000000, 100)
}

type CacheFactory struct {
	userCache ICache[string, *userModel.User]
}

func NewCacheFactory() *CacheFactory {
	userCache, err := NewCache[string, *userModel.User]()
	if err != nil {
		panic(err)
	}

	return &CacheFactory{
		userCache: userCache,
	}
}

func (f *CacheFactory) UserCache() ICache[string, *userModel.User] {
	return f.userCache
}
