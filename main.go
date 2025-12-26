package main

import (
	"fmt"

	"huo_jian_qiang/cmd"
	"huo_jian_qiang/internal/constant"
	"huo_jian_qiang/internal/download"
	"huo_jian_qiang/internal/http"
	"huo_jian_qiang/internal/logger"
	"os"
)

func main() {
	logger.InitDefault(constant.APP_NAME, logger.INFO)

	err := download.CreateUsersTxt(".")

	if err != nil {
		logger.Errorf("创建用户文件失败: %v", err)
		os.Exit(1)
	}

	d := download.PasswordDownloader("downloads")

	localPath, err := d.DownloadFile(constant.PASSWORD_LIST)

	if err != nil {
		logger.Errorf("下载密码文件失败: %v", err)
		os.Exit(1)
	}

	logger.Infof("密码文件已下载: %s", localPath)

	if err != nil {
		logger.Errorf("下载出错: %v", err)
		os.Exit(1)
	}

	cfg, err := cmd.ReadConfig()

	if err != nil {
		if err.Error() == "flag: help requested" {
			printUsage()
			os.Exit(0)
		}

		logger.Errorf("配置读取失败: %v", err)
		printUsage()
		os.Exit(1)
	}

	startProcessing(cfg)
}

func printUsage() {
	fmt.Println()
	fmt.Printf("【%s】", constant.APP_NAME)
	fmt.Println()

	fmt.Println("一个快速、可靠的多线程请求工具，专为大规模测试设计")
	fmt.Println("  就像火箭一样精准、快速！")
	fmt.Println()

	fmt.Println("使用方法:")
	fmt.Println("  go run main.go \\")
	fmt.Println("    -u <用户文件> \\")
	fmt.Println("    -p <密码文件> \\")
	fmt.Println("    -a <请求体> \\")
	fmt.Println("    [其他选项]")
	fmt.Println()

	fmt.Println("必填参数 (缺一不可):")
	fmt.Println("  -u  string  包含用户名的文件")
	fmt.Println("              示例: users.txt, emails.csv")
	fmt.Println()
	fmt.Println("  -p  string  包含密码的文件")
	fmt.Println("              示例: passwords.txt, wordlist.txt")
	fmt.Println()
	fmt.Println("  -a  string  请求体模板 (支持 *USER* 和 *PASS* 占位符)")
	fmt.Println(`              示例: '{"user":"*USER*","pass":"*PASS*"}'`)
	fmt.Println(`              示例: 'user=*USER*&password=*PASS*&submit=login'`)

	fmt.Println()

	fmt.Println("可选参数 (锦上添花):")
	fmt.Println("  -d  int     请求间隔 (秒)")
	fmt.Println("              默认: 0 (无延迟)")
	fmt.Println("              建议: 1-5 (避免封禁)")
	fmt.Println()
	fmt.Println("  -t  int     线程数 (并发数量)")
	fmt.Println("              默认: 1 (单线程)")
	fmt.Println("              建议: 10-50 (根据目标调整)")
	fmt.Println()
}

func startProcessing(cfg *cmd.Config) {
	logger.Infof("-> 开始处理...")

	response, status, err := http.PostRequest(cfg.Url, cfg.RequestBody, cfg.Delay)

	if err != nil {
		logger.Errorf("-> 请求出错: %v", err)
		return
	}

	if status != 200 {
		logger.Errorf("请求失败: 状态码 %d", status)
		return
	}

	logger.Infof("-> 响应: %s", response)

	logger.Infof("-> 处理完成！")
}
