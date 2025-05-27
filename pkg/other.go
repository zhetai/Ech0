package pkg

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
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
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		// 只有当错误不是"文件不存在"时才返回错误
		return err
	}
	return nil
}

type Header struct {
	Header  string
	Content string
}

// 请求发送函数
func SendRequest(url, method string, customHeader Header) ([]byte, error) {
	// 加载系统根证书池
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("加载系统证书失败: %w", err)
	}

	// 自定义 HTTP 客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: rootCAs,
			},
		},
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// 添加自定义请求头
	if customHeader.Header != "" && customHeader.Content != "" {
		req.Header.Set(customHeader.Header, customHeader.Content)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求发送失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	return body, nil
}
