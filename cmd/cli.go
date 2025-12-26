package cmd

import (
	"flag"
	"fmt"
	"huo_jian_qiang/internal/logger"
	"huo_jian_qiang/internal/warning"
	"strings"
)

type Config struct {
	Url           string
	UsersFile     string
	PasswordsFile string
	RequestBody   string
	Delay         int
	Threads       int
	Method        string
}

func ReadConfig() (*Config, error) {
	cfg := &Config{}
	flag.StringVar(&cfg.Url, "url", "", "目标URL(必填)")
	flag.StringVar(&cfg.UsersFile, "u", "", "包含用户名的文件(必填)")
	flag.StringVar(&cfg.PasswordsFile, "p", "", "包含密码的文件(必填)")
	flag.StringVar(&cfg.RequestBody, "a", "", "附加用户输入(必填)")
	flag.IntVar(&cfg.Delay, "d", 0, "请求间隔(秒)")
	flag.IntVar(&cfg.Threads, "t", 0, "线程数")
	flag.StringVar(&cfg.Method, "m", "GET", "请求方法(GET/POST)")

	flag.Parse()

	if cfg.Url == "" {
		return nil, fmt.Errorf("缺少必填参数：-url (目标URL)")
	}

	if _, err := warning.DisplayWarning(cfg.Url); err != nil {
		return nil, err
	}

	if cfg.UsersFile == "" {
		return nil, fmt.Errorf("缺少必填参数：-u (用户名字段文件)")
	}

	if cfg.PasswordsFile == "" {
		return nil, fmt.Errorf("缺少必填参数：-p(密码字段文件)")
	}

	if cfg.RequestBody == "" {
		return nil, fmt.Errorf("缺少必填参数：-a(附加用户输入)")
	}

	if cfg.Method == "" {
		return nil, fmt.Errorf("缺少必填参数：-m(请求方法) [GET/POST/PUT]")
	}

	cfg.UsersFile = getDefaultFileName(cfg.UsersFile)
	cfg.PasswordsFile = getDefaultFileName(cfg.PasswordsFile)

	logger.Infof("┌─ 链接:         %s", cfg.Url)
	logger.Infof("├─ 用户文件:      %s", cfg.UsersFile)
	logger.Infof("├─ 密码文件:      %s", cfg.PasswordsFile)
	logger.Infof("├─ 请求体:        %s", cfg.RequestBody)
	logger.Infof("├─ 线程数:        %d", cfg.Threads)
	logger.Infof("├─ 请求方法:      %s", cfg.Method)
	logger.Infof("└─ 延迟:          %d 秒", cfg.Delay)

	return cfg, nil
}

func getDefaultFileName(user_input string) string {
	directory := "downloads/data/"
	user_input = strings.TrimSpace(user_input)

	if user_input == "*U*" {
		return directory + "users.txt"
	}

	if user_input == "*P*" {
		return directory + "rockyou.txt"
	}

	return user_input
}
