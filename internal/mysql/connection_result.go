package mysql

type ConnectionResult struct {
	Success   bool   // 是否连接成功
	User      string // 用户名
	Password  string // 密码
	Host      string // 主机地址
	Error     error  // 错误信息（如果有）
	AuthError bool   // 是否为认证错误
}
