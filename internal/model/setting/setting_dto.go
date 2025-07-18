package model

// SystemSettingDto 定义系统设置数据传输对象
type SystemSettingDto struct {
	SiteTitle     string `json:"site_title"`     // 站点标题
	ServerName    string `json:"server_name"`    // 服务器名称
	ServerURL     string `json:"server_url"`     // 服务器地址
	AllowRegister bool   `json:"allow_register"` // 是否允许注册
	ICPNumber     string `json:"ICP_number"`     // 备案号
	MetingAPI     string `json:"meting_api"`     // Meting API 地址
	CommentAPI    string `json:"comment_api"`    // 评论 API 地址
	CustomCSS     string `json:"custom_css"`     // 自定义 CSS
	CustomJS      string `json:"custom_js"`      // 自定义 JS
}
