package units

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Crawler performs web crawling to extract content from a URL
func Crawler(url string) (string, error) {
	fmt.Printf("正在使用 URL 进行自定义爬取:%s\n", url)

	reqBody := map[string]string{
		"url": url,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %v", err)
	}

	resp, err := http.Post("https://crawl.search1api.com", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("API请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败, 状态码: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return "", fmt.Errorf("收到的响应不是有效的JSON格式")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("JSON解码失败: %v", err)
	}

	fmt.Println("自定义爬取服务调用完成")

	// Convert the result back to JSON string
	responseData, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("响应JSON编码失败: %v", err)
	}

	return string(responseData), nil
}
