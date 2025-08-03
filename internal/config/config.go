package config

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"log"
	"os"

	model "github.com/lin-snow/ech0/internal/model/common"
	"github.com/spf13/viper"
)

// Config 全局配置变量
var Config AppConfig

// JWT_SECRET 用于JWT签名的密钥
var JWT_SECRET []byte

// AppConfig 应用程序配置结构体
type AppConfig struct {
	Server struct {
		Port string `yaml:"port"` // 服务器端口
		Host string `yaml:"host"` // 服务器主机地址
		Mode string `yaml:"mode"` // 运行模式，可能的值为 "debug" 或 "release"
	} `yaml:"server"`
	Database struct {
		Type string `yaml:"type"` // 数据库类型
		Path string `yaml:"path"` // 数据库文件路径
	} `yaml:"database"`
	Auth struct {
		Jwt struct {
			Expires  int    `yaml:"expires"`  // JWT的过期时间，单位为秒
			Issuer   string `yaml:"issuer"`   // JWT的发行者
			Audience string `yaml:"audience"` // JWT的受众
		} `yaml:"jwt"`
	} `yaml:"auth"`
	Upload struct {
		ImageMaxSize int      `yaml:"imagemaxsize"` // 图片文件的最大上传大小，单位为字节
		AudioMaxSize int      `yaml:"audiomaxsize"` // 音频文件的最大上传大小，单位为字节
		AllowedTypes []string `yaml:"allowedtypes"` // 允许上传的文件类型
		ImagePath    string   `yaml:"imagepath"`    // 图片文件存储路径
		AudioPath    string   `yaml:"audiopath"`    // 音频文件存储路径
	} `yaml:"upload"`
	Setting struct {
		SiteTitle     string `yaml:"sitetitle"`     // 网站标题
		Servername    string `yaml:"servername"`    // 服务器名称
		Serverurl     string `yaml:"serverurl"`     // 服务器 URL
		AllowRegister bool   `yaml:"allowregister"` // 是否允许注册
		Icpnumber     string `yaml:"icpnumber"`     // ICP 备案号
		MetingAPI     string `yaml:"metingapi"`     // Meting API 地址
		CustomCSS     string `yaml:"customcss"`     // 自定义 CSS 样式
		CustomJS      string `yaml:"customjs"`      // 自定义 JS 脚本
	} `yaml:"setting"`
	Comment struct {
		EnableComment bool   `yaml:"enablecomment"` // 是否启用评论
		Provider      string `yaml:"provider"`      // 评论提供者
		CommentAPI    string `yaml:"commentapi"`    // 评论 API 地址
	} `yaml:"comment"`
	SSH struct {
		Port string `yaml:"port"` // SSH 端口
		Host string `yaml:"host"` // SSH 主机地址
		Key  string `yaml:"key"`  // SSH 私钥路径
	} `yaml:"ssh"`
}

//go:embed config.yaml
var configData []byte

// LoadAppConfig 加载应用程序配置
func LoadAppConfig() {
	// viper.SetConfigFile("config/config.yaml")
	viper.SetConfigType("yaml")
	// 使用嵌入的配置数据而不是从文件系统读取
	err := viper.ReadConfig(bytes.NewReader(configData))
	if err != nil {
		panic(model.READ_CONFIG_PANIC + ":" + err.Error())
	}

	// 将配置文件内容反序列化到结构体 Config 中
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(model.READ_CONFIG_PANIC + ":" + err.Error())
	}

	JWT_SECRET = GetJWTSecret()
}

// GetJWTSecret 加载JWT密钥
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

	return []byte(secret)
}
