package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"huo_jian_qiang/internal/constant"
	"strings"
	"time"
)

// AttemptConnect 尝试使用给定的用户名和密码连接到MySQL数据库
// 返回: (是否成功, 错误信息)
func AttemptConnect(user string, password string, host string) (bool, error) {
	return AttemptConnectWithConfig(user, password, constant.Config{
		Host:    host,
		Port:    3306,
		Timeout: 5 * time.Second,
	})
}

// AttemptConnectWithConfig 使用配置尝试连接MySQL
func AttemptConnectWithConfig(user string, password string, config constant.Config) (bool, error) {
	var result int

	// 设置默认值
	if config.Port == 0 {
		config.Port = 3306
	}
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}

	// 构建DSN
	dsn := buildDSN(user, password, config)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return false, fmt.Errorf("连接数据库失败: %w", err)
	}
	defer db.Close()

	// 设置连接参数
	db.SetConnMaxLifetime(config.Timeout)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// 测试连接
	err = db.PingContext(ctx)
	if err != nil {
		if isAuthError(err) {
			return false, nil // 认证失败，返回false但不返回error
		}
		return false, fmt.Errorf("测试连接数据库失败: %w", err)
	}

	// 执行简单查询验证连接
	query := "SELECT 1"
	if config.Database != "" {
		// 尝试切换到指定数据库
		query = fmt.Sprintf("USE %s; SELECT 1", config.Database)
	}

	err = db.QueryRowContext(ctx, query).Scan(&result)
	if err != nil {
		// 如果是指定数据库且数据库不存在，这可能是一个错误
		if config.Database != "" && strings.Contains(err.Error(), "Unknown database") {
			return true, nil // 连接成功但数据库不存在
		}
		return false, fmt.Errorf("连接成功后查询失败: %w", err)
	}

	if result != 1 {
		return false, errors.New("查询结果异常")
	}

	return true, nil
}

// AttemptConnectBatch 批量尝试连接，返回第一个成功的连接
func AttemptConnectBatch(host string, credentials []struct {
	User     string
	Password string
}) (ConnectionResult, error) {
	return AttemptConnectBatchWithConfig(constant.Config{
		Host:    host,
		Port:    3306,
		Timeout: 5 * time.Second,
	}, credentials)
}

// AttemptConnectBatchWithConfig 批量尝试连接（带配置）
func AttemptConnectBatchWithConfig(config constant.Config, credentials []struct {
	User     string
	Password string
}) (ConnectionResult, error) {
	for _, cred := range credentials {
		success, err := AttemptConnectWithConfig(cred.User, cred.Password, config)
		result := ConnectionResult{
			Success:  success,
			User:     cred.User,
			Password: cred.Password,
			Host:     config.Host,
			Error:    err,
		}

		if success {
			return result, nil
		}

		// 如果是认证错误，继续尝试下一个
		if err == nil {
			result.AuthError = true
		} else {
			// 其他错误（如网络问题）可能意味着需要停止
			return result, fmt.Errorf("连接尝试失败: %w", err)
		}
	}

	return ConnectionResult{
		Success: false,
		Host:    config.Host,
	}, errors.New("所有凭证都失败")
}

// buildDSN 构建MySQL连接字符串
func buildDSN(user, password string, config constant.Config) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, password, config.Host, config.Port)

	// 添加参数
	params := []string{
		"timeout=" + config.Timeout.String(),
		"readTimeout=" + config.Timeout.String(),
		"writeTimeout=" + config.Timeout.String(),
		"parseTime=true",
		"charset=utf8mb4",
	}

	if config.Database != "" {
		dsn = fmt.Sprintf("%s%s", dsn, config.Database)
	}

	return fmt.Sprintf("%s?%s", dsn, strings.Join(params, "&"))
}

// isAuthError 检查错误是否为认证错误
func isAuthError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	// MySQL认证相关错误关键词
	authErrors := []string{
		"access denied",
		"invalid credentials",
		"authentication",
		"password",
		"no such user",
		"1045", // MySQL错误码：访问被拒绝
		"1044", // MySQL错误码：拒绝访问数据库
		"1042", // MySQL错误码：错误的主机名
		"1043", // MySQL错误码：握手错误
	}

	for _, authErr := range authErrors {
		if strings.Contains(errStr, authErr) {
			return true
		}
	}
	return false
}

// IsConnectionError 检查是否为连接错误（网络问题等）
func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	connectionErrors := []string{
		"connection refused",
		"connection timeout",
		"no route to host",
		"network is unreachable",
		"i/o timeout",
		"dial tcp",
		"2003", // MySQL错误码：无法连接到MySQL服务器
		"2005", // MySQL错误码：未知的MySQL服务器主机
	}

	for _, connErr := range connectionErrors {
		if strings.Contains(errStr, connErr) {
			return true
		}
	}
	return false
}
