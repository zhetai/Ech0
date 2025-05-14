package dto

type ImageDto struct {
	// 图片的 URL 地址
	URL string `json:"url" binding:"required"`
}
