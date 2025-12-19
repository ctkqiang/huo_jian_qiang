package http

import (
	"fmt"
	"net/http"

	"huo_jian_qiang/pkg/logger"
)

var log = logger.New("debug")

func GET(url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("请求GET %s 失败，状态码: %d", url, resp.StatusCode)

		return fmt.Errorf("GET %s 返回非200状态码: %d", url, resp.StatusCode)
	}

	log.Info("请求GET %s 成功，状态码: %d", url, resp.StatusCode)

	return nil
}
