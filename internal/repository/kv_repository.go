package repository

import (
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/pkg"
)

// 处理Key Value表的增删改查

// 增加Key Value
func AddKeyValue[T any](key string, value T) error {
	// 处理数据(序列化value)
	resolvedValue, err := (pkg.JSONMarshal(value))
	if err != nil {
		return err
	}

	// 插入数据
	database.DB.Create(&models.KeyValue{
		Key:   key,
		Value: string(resolvedValue),
	})

	return nil
}

// 删除Key Value
func DeleteKeyValue(key string) error {
	// 删除数据
	if err := database.DB.Where("key = ?", key).Delete(&models.KeyValue{}).Error; err != nil {
		return err
	}
	return nil
}

// 获取Key Value
func GetKeyValue[T any](key string) (T, error) {
	var kv models.KeyValue
	var value T

	// 查询数据
	if err := database.DB.Where("key = ?", key).First(&kv).Error; err != nil {
		return value, err
	}

	// 处理数据(反序列化value)
	if err := pkg.JSONUnmarshal([]byte(kv.Value), &value); err != nil {
		return value, err
	}

	return value, nil
}

// 更新Key Value
func UpdateKeyValue[T any](key string, value T) error {
	// 处理数据(序列化value)
	resolvedValue, err := (pkg.JSONMarshal(value))
	if err != nil {
		return err
	}

	// 更新数据
	if err := database.DB.Model(&models.KeyValue{}).Where("key = ?", key).Update("value", string(resolvedValue)).Error; err != nil {
		return err
	}
	return nil
}
