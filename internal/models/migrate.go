package models

import (
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &Message{}, &Image{}, &Todo{}, &KeyValue{}, &Connected{})
	if err != nil {
		return err
	}
	return nil
}

// 从旧的 Message 表中迁移图片数据到新的 Image 表
func MigrateMessageImages(db *gorm.DB) error {
	var messages []Message
	// 找出旧图字段非空的消息
	if err := db.Where("image_url IS NOT NULL AND image_url != ''").Find(&messages).Error; err != nil {
		return err
	}

	for _, msg := range messages {
		var count int64
		// 如果该消息已经有 image 记录，就跳过
		err := db.Model(&Image{}).Where("message_id = ? AND image_url = ?", msg.ID, msg.ImageURL).Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			continue // 已存在，跳过
		}

		img := Image{
			MessageID:   msg.ID,
			ImageURL:    msg.ImageURL,
			ImageSource: msg.ImageSource,
		}
		if err := db.Create(&img).Error; err != nil {
			return err
		}
	}
	return nil
}
