package util

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// TrimURL 去除 URL 前后的空格和斜杠
func TrimURL(url string) string {
	if url == "" {
		return ""
	}

	// 去除连接地址前后的空格和斜杠
	url = strings.TrimSpace(url)
	url = strings.TrimPrefix(url, "/")
	url = strings.TrimSuffix(url, "/")
	return url
}

// Header 自定义请求头结构体
type Header struct {
	Header  string
	Content string
}

// SendRequest 发送 HTTP 请求
func SendRequest(url, method string, customHeader Header) ([]byte, error) {
	// 自定义 HTTP 客户端，忽略 TLS 证书验证
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 忽略证书验证
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
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Println("Failed to close response body:", closeErr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	return body, nil
}
