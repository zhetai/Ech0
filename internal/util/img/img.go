package util

import (
	"github.com/disintegration/imaging"
)

// GetImageSize 使用 imaging 获取图片宽高
func GetImageSize(path string) (int, int, error) {
	img, err := imaging.Open(path, imaging.AutoOrientation(true))
	if err != nil {
		return 0, 0, err
	}
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	return width, height, nil
}
