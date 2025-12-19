package main

import (
	"huo_jian_qiang/internal/logger"
)

func main() {
	logger.InitDefault("火尖枪", logger.DEBUG)
	logger.Infof("火箭悄然点火，工具已苏醒！")
}
