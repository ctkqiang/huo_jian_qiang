package main

import (
	"fmt"
	"huo_jian_qiang/cmd"
	"huo_jian_qiang/internal/logger"
	"os"
)

func main() {
	logger.InitDefault("火尖枪", logger.INFO)
	logger.Infof("火箭悄然点火，工具已苏醒！")

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

	logger.Infof("配置加载成功:")
	logger.Infof("  用户文件: %s", cfg.UsersFile)
	logger.Infof("  密码文件: %s", cfg.PasswordsFile)
	logger.Infof("  请求体: %s", cfg.RequestBody)

	if cfg.Delay > 0 {
		logger.Infof("  请求间隔: %d秒", cfg.Delay)
	}

	if cfg.Threads > 0 {
		logger.Infof("  线程数: %d", cfg.Threads)
	}

	startProcessing(cfg)
}

func printUsage() {
	fmt.Println("🔥 火尖枪 - 高性能请求工具")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  go run main.go -u <用户文件> -p <密码文件> -a <请求体> [选项]")
	fmt.Println()
	fmt.Println("必填参数:")
	fmt.Println("  -u string   包含用户名的文件")
	fmt.Println("  -p string   包含密码的文件")
	fmt.Println("  -a string   附加用户输入（请求体模板）")
	fmt.Println()
	fmt.Println("可选参数:")
	fmt.Println("  -d int      请求间隔（秒）")
	fmt.Println("  -t int      线程数")
	fmt.Println()
}

// startProcessing 开始处理逻辑
func startProcessing(cfg *cmd.Config) {
	logger.Infof("开始处理...")
	logger.Infof("处理完成！")
}
