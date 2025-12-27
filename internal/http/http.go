package http

import (
	"fmt"
	"huo_jian_qiang/internal/logger"
	"huo_jian_qiang/internal/warning"
	"io"
	"net/http"
	"strings"
	"time"
)

func PostRequest(basedUrl, body string, timeout int) (string, int, error) {
	if !strings.HasPrefix(basedUrl, "http://") && !strings.HasPrefix(basedUrl, "https://") {
		return "", 0, fmt.Errorf("URL格式错误：必须以http://或https://开头")
	}

	client := createHTTPClient(timeout)

	request, err := http.NewRequest("POST", basedUrl, strings.NewReader(body))

	if err != nil {
		return "", 0, fmt.Errorf("创建请求失败: %v", err)
	}

	request.Header.Set("Accept", "*/*")
	request.Header.Set("User-Agent", "Mozilla/5.0")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	logger.Infof("发送POST请求: URL=%s, 内容长度=%d, Content-Type=%s", basedUrl, len(body), request.Header.Get("Content-Type"))

	resp, err := client.Do(request)

	if err != nil {
		if err == io.EOF {
			return "", 0, fmt.Errorf("服务器已关闭")
		}

		if strings.Contains(err.Error(), "timeout") {
			return "", 0, fmt.Errorf("请求超时")
		}

		if strings.Contains(err.Error(), "unsupported protocol scheme") {
			return "", 0, fmt.Errorf("不支持的协议方案或链接不存在")
		}

		if strings.Contains(err.Error(), "connection reset by peer") {
			return "", 0, fmt.Errorf("服务器已关闭")
		}

		return "", 0, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("读取响应失败: %v", err)
	}

	switch resp.StatusCode {
	case 504:
		logger.Warnf("网关超时 | 状态码: %d", resp.StatusCode)
	case 408:
		logger.Warnf("请求超时 | 响应体: %s", string(bodyBytes))
	case 400:
		logger.Warnf("请求错误 | 状态码: %d", resp.StatusCode)
	case 403:
		logger.Warnf("请求被限制 | 状态码: %d", resp.StatusCode)
	case 429:
		logger.Warnf("收到[429]状态码 | 请求被限制, 建议增加延迟后重试")
	case 200:
	case 201:
		logger.Infof("请求成功！")
		logger.Infof("状态码：%d", resp.StatusCode)
		logger.Infof("响应长度：%d 字节", len(bodyBytes))

		if len(body) > 0 {
			logger.Infof("    响应体预览：")

			previewLength := 300

			if len(body) < previewLength {
				previewLength = len(body)
			}

			preview := body[:previewLength]

			preview = strings.ReplaceAll(preview, "\n", "↲ ")
			preview = strings.ReplaceAll(preview, "\r", "")

			logger.Infof("   ┌──────────────────────────────────────")
			logger.Infof("   │ %s", preview)

			if len(body) > previewLength {
				logger.Infof("   │ ... (还有 %d 个字符)", len(body)-previewLength)
			}

			logger.Infof("   └──────────────────────────────────────")
		}
	default:
		logger.Infof("   └──────────────────────────────────────")
	}

	responseBody := string(bodyBytes)
	return responseBody, resp.StatusCode, nil
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

	ipaddrv4, ipaddrv6, err := warning.GetIP(baseURL)

	if err != nil {
		return "", fmt.Errorf("解析IP地址失败: %v", err)
	}

	logger.Infof("地址 %s: IPv4=%s, IPv6=%s", baseURL, ipaddrv4, ipaddrv6)

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
