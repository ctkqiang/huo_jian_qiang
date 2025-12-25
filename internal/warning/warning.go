package warning

import "strings"

func DisplayWarning(domain string) (string, error) {
	if strings.Contains(domain, ".gov.cn") {
		return "警告：根据中国网络安全法规，禁止对.gov.cn域名进行未经授权的扫描或测试。请立即停止操作并确保已获得合法授权。", nil
	}

	return "", nil
}
