package model

// UserInfoDto 定义用户信息数据传输对象
type UserInfoDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
	Avatar   string `json:"avatar"`
}
