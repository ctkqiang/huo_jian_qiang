package constant

import "time"

type Config struct {
	Host     string        // MySQL主机地址
	Port     int           // MySQL端口，默认3306
	Timeout  time.Duration // 连接超时时间，默认5秒
	Database string        // 数据库名，可选
}
