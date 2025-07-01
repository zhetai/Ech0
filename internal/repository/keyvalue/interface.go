package keyvalue

type KeyValueRepositoryInterface interface {
	// GetKeyValue 根据键获取值
	GetKeyValue(key string) (interface{}, error)

	// AddKeyValue 添加键值对
	AddKeyValue(key string, value interface{}) error

	// DeleteKeyValue 删除键值对
	DeleteKeyValue(key string) error

	// UpdateKeyValue 更新键值对
	UpdateKeyValue(key string, value interface{}) error
}
