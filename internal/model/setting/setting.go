package model

// SystemSetting 定义系统设置实体
type SystemSetting struct {
	SiteTitle     string `json:"site_title"`     // 站点标题
	ServerName    string `json:"server_name"`    // 服务器名称
	ServerURL     string `json:"server_url"`     // 服务器地址
	AllowRegister bool   `json:"allow_register"` // 是否允许注册'
	ICPNumber     string `json:"ICP_number"`     // 备案号
	MetingAPI     string `json:"meting_api"`     // Meting API 地址
	CustomCSS     string `json:"custom_css"`     // 自定义 CSS
	CustomJS      string `json:"custom_js"`      // 自定义 JS
}

// CommentSetting 定义评论设置实体
type CommentSetting struct {
	EnableComment bool   `json:"enable_comment"` // 是否启用评论
	Provider      string `json:"provider"`       // 评论提供者
	CommentAPI    string `json:"comment_api"`    // 评论 API 地址
}

// S3Setting 定义 S3 存储设置实体
type S3Setting struct {
	Enable     bool   `json:"enable"`      // 是否启用 S3 存储
	Provider   string `json:"provider"`    // S3 服务提供商，例如 "aws", "aliyun", "minio", "other"
	Endpoint   string `json:"endpoint"`    // S3 端点
	AccessKey  string `json:"access_key"`  // 访问密钥 ID
	SecretKey  string `json:"secret_key"`  // 秘密访问密钥
	BucketName string `json:"bucket_name"` // 存储桶名称
	Region     string `json:"region"`      // 区域
	UseSSL     bool   `json:"use_ssl"`     // 是否使用 SSL
	CDNURL     string `json:"cdn_url"`     // CDN 加速域名（可选，没有就走 Endpoint）
	PathPrefix string `json:"path_prefix"` // 存储路径前缀，例如 "uploads/"，方便隔离目录
	PublicRead bool   `json:"public_read"` // 上传时是否默认设置对象为 public-read
}
