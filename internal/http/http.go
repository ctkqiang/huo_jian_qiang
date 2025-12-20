package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Requestbody struct {
	username string
	password string
} //

func PostRequest(basedUrl, body string, timeout int) (string, int, error) {
	fullURL, err := buildURL(basedUrl, "")

	if err != nil {
		return "", 0, fmt.Errorf("构建URL失败: %v", err)
	}

	// client := createHTTPClient(timeout)

	request, err := http.NewRequest("POST", fullURL, strings.NewReader(body))

	if err != nil {
		return "", 0, fmt.Errorf("创建请求失败: %v", err)
	}

	request.Header.Set("Accept", "*/*")

	return "", 0, nil
}

func GetRequest(baseURL, paramA string, timeout int) (string, int, error) {
	fullURL, err := buildURL(baseURL, paramA)
	if err != nil {
		return "", 0, fmt.Errorf("构建URL失败: %v", err)
	}

	client := createHTTPClient(timeout)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", 0, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "huo_jian_qiang/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("读取响应失败: %v", err)
	}

	responseBody := string(bodyBytes)

	return responseBody, resp.StatusCode, nil
}

func buildURL(baseURL, paramA string) (string, error) {
	if baseURL == "" {
		return "", fmt.Errorf("基础URL不能为空")
	}

	if paramA == "" {
		return baseURL, nil
	}

	if strings.Contains(baseURL, "?") {
		return fmt.Sprintf("%s&a=%s", baseURL, paramA), nil
	}

	return fmt.Sprintf("%s?a=%s", baseURL, paramA), nil
}

func createHTTPClient(timeout int) *http.Client {
	if timeout <= 0 {
		timeout = 30
	}

	return &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
}
