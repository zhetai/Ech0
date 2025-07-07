package model

// PageQueryDto 用于分页查询的请求数据传输对象
type PageQueryDto struct {
	Page     int    `json:"page" form:"page"`         // 页码，从1开始
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页大小
	Search   string `json:"search" form:"search"`     // 用于搜索的关键字
}

// ImageDto 用于图片相关的请求数据传输对象
type ImageDto struct {
	// 图片的 URL 地址
	URL    string `json:"url" binding:"required"`
	SOURCE string `json:"source" binding:"required"`
}
