package model

// PageQueryDto 用于分页查询的请求数据传输对象
type PageQueryDto struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Search   string `json:"search"`
}

// ImageDto 用于图片相关的请求数据传输对象
type ImageDto struct {
	// 图片的 URL 地址
	URL    string `json:"url" binding:"required"`
	SOURCE string `json:"source" binding:"required"`
}
