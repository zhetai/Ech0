package pkg

import (
	"io"
	"net/http"
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

// 请求发送函数
func SendRequest(url, method string, cutsomHeader struct {
	Header  string
	Content string
}) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// 添加自定义请求头
	if cutsomHeader.Header != "" && cutsomHeader.Content != "" {
		req.Header.Set(cutsomHeader.Header, cutsomHeader.Content)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
