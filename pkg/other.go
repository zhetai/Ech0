package pkg

import (
	"os"
	"strings"
)

func TrimURL(url string) string {
	// 去除连接地址前后的空格和斜杠
	url = strings.TrimSpace(url)
	url = strings.TrimPrefix(url, "/")
	url = strings.TrimSuffix(url, "/")
	return url
}

// 删除文件
func DeleteFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}
