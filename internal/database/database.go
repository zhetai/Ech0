package database

import (
	"fmt"
	"log"
	"os"

	"github.com/lin-snow/ech0/config"
	"github.com/lin-snow/ech0/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() error {
	// 读取数据库类型和路径
	dbType := config.Config.Database.Type
	dbPath := config.Config.Database.Path

	// 确保数据库目录存在
	dir := fmt.Sprintf("%s", dbPath[:len(dbPath)-len("/ech0.db")]) // 提取目录部分
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
		return err
	}

	var err error
	// 根据数据库类型选择不同的数据库驱动
	if dbType == "sqlite" {
		DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	} else {
		log.Fatalf(models.DatabaseTypeMessage+": %s", dbType)
		return err
	}

	if err != nil {
		log.Fatalf(models.DatabaseConnectionError+": %v", err)
		return err
	}

	if err = models.MigrateDB(DB); err != nil {
		log.Fatal(models.DatabaseMigrationError+":", err)
	}

	// log.Println(models.DatabaseConnectionSuccess)
	return nil
}
