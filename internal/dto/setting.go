package dto

type SystemSettingDto struct {
	ServerName    string `json:"server_name"`    // 服务器名称
	AllowRegister bool   `json:"allow_register"` // 是否允许注册
}
