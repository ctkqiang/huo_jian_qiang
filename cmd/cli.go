package cmd

import (
	"flag"
	"fmt"
)

type Config struct {
	UsersFile     string
	PasswordsFile string
	RequestBody   string
	Delay         int
	Threads       int
}

func ReadConfig() (*Config, error) {
	cfg := &Config{}
	flag.StringVar(&cfg.UsersFile, "u", "", "包含用户名的文件(必填)")
	flag.StringVar(&cfg.PasswordsFile, "p", "", "包含密码的文件(必填)")
	flag.StringVar(&cfg.RequestBody, "a", "", "附加用户输入(必填)")
	flag.IntVar(&cfg.Delay, "d", 0, "请求间隔(秒)")
	flag.IntVar(&cfg.Threads, "t", 0, "线程数")

	flag.Parse()

	if cfg.UsersFile == "" {
		return nil, fmt.Errorf("缺少必填参数：-u (用户名字段文件)")
	}

	if cfg.PasswordsFile == "" {
		return nil, fmt.Errorf("缺少必填参数：-p(密码字段文件)")
	}

	if cfg.RequestBody == "" {
		return nil, fmt.Errorf("缺少必填参数：-a(附加用户输入)")
	}

	return cfg, nil
}
