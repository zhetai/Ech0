package model

// PageQueryDto 用于分页查询的请求数据传输对象
type PageQueryDto struct {
	Page     int    `json:"page" from:"page"`
	PageSize int    `json:"pageSize" from:"pageSize"`
	Search   string `json:"search" from:"search"` // 用于搜索的关键字
}

// ImageDto 用于图片相关的请求数据传输对象
type ImageDto struct {
	// 图片的 URL 地址
	URL    string `json:"url" binding:"required"`
	SOURCE string `json:"source" binding:"required"`
}
