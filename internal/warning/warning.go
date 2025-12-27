package warning

import (
	"fmt"
	"net"
	"strings"
)

func DisplayWarning(domain string) (string, error) {
	var message string = "[警告：根据中国网络安全法规，禁止对.gov.cn域名进行未经授权的扫描或测试。请立即停止操作并确保已获得合法授权。]"

	if strings.HasSuffix(domain, ".gov.cn") {
		return message, fmt.Errorf("%s", message)
	}

	lowerDomain := strings.ToLower(domain)

	if strings.HasSuffix(lowerDomain, ".gov.cn") {
		return message, fmt.Errorf("%s", message)
	}

	return "", nil
}

func GetIP(url string) (string, string, error) {
	var ipv4 string
	var ipv6 string
	ip, err := net.LookupIP(url)

	if err != nil {
		return "", "", fmt.Errorf("解析域名失败: %v", err)
	}

	for _, _ip := range ip {
		if _ip.To4() != nil {
			ipv4 = _ip.String()
		} else {
			ipv6 = _ip.String()
		}
	}

	return ipv4, ipv6, nil
}
