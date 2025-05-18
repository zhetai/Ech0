package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
	Database struct {
		Type string `yaml:"type"`
		Path string `yaml:"path"`
	} `yaml:"database"`
	Setting struct {
		SiteTitle     string `yaml:"sitetitle"`
		Servername    string `yaml:"servername"`
		Serverurl     string `yaml:"serverurl"`
		AllowRegister bool   `yaml:"allowregister"`
		Icpnumber     string `yaml:"icpnumber"`
	}
	Upload struct {
		MaxSize      int      `yaml:"maxsize"`
		AllowedTypes []string `yaml:"allowedtypes"`
		ImagePath    string   `yaml:"imagepath"`
		AudioPath    string   `yaml:"audiopath"`
	} `yaml:"upload"`
	Auth struct {
		Jwt struct {
			Expires  int    `yaml:"expires"`
			Issuer   string `yaml:"issuer"`
			Audience string `yaml:"audience"`
		} `yaml:"jwt"`
	} `yaml:"auth"`
}

var Config AppConfig
var JWT_SECRET []byte

func LoadConfig() error {
	viper.SetConfigFile("config/config.yaml")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Failed to load config: %s", err)
		return err
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Printf("Failed to parse config: %s", err)
		return err
	}

	JWT_SECRET = GetJWTSecret()

	return nil
}

// 获取JWT密钥
func GetJWTSecret() []byte {
	// 从环境变量中获取JWT密钥
	secret := os.Getenv("JWT_SECRET")
	if secret == "" { // 如果没有设置环境变量，则使用UUID生成默认密钥
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			log.Fatal("failed to generate random JWT secret:", err)
		}
		secret = hex.EncodeToString(b)
	}

	// log.Println("JWT secret:", secret)

	return []byte(secret)
}
