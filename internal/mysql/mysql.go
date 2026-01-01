package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func AttemptConnect(user string, password string, host string) (bool, error) {
	var result int

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/?timeout=5s", user, password, host)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return false, fmt.Errorf("连接数据库失败: %w", err)
	}

	defer db.Close()

	db.SetConnMaxLifetime(5 * time.Second)
	db.SetMaxOpenConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {

		if isAuthError(err) {
			return false, nil
		}

		return false, fmt.Errorf("测试连接数据库失败: %w", err)
	}

	err = db.QueryRowContext(ctx, "SELECT 1").Scan(&result)

	if err != nil {
		return false, fmt.Errorf("连接成功后查询失败: %w", err)
	}

	if result != 1 {
		return false, errors.New("查询结果异常")
	}

	return true, nil
}

func isAuthError(err error) bool {
	errStr := err.Error()
	authErrors := []string{
		"access denied",       // 访问被拒绝
		"invalid credentials", // 凭据无效
		"1045",                // MySQL 访问被拒绝的错误码
		"authentication",      // 认证
		"password",            // 密码
	}

	for _, authErr := range authErrors {
		if containsIgnoreCase(errStr, authErr) {
			return true
		}
	}
	return false
}

func containsIgnoreCase(s, substr string) bool {
	sLower := toLower(s)
	substrLower := toLower(substr)

	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))

	for i := 0; i < len(s); i++ {
		c := s[i]

		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}
