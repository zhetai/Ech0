package models

import (
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &Message{})
	if err != nil {
		return err
	}
	return nil
}
