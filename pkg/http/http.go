package http

import (
	"fmt"
	"io"
	"net/http"

	"huo_jian_qiang/internal/model"
	"huo_jian_qiang/pkg/logger"
)

var log = logger.New("debug")
var userRequest model.UserRequest

func GET(url string, userRequest model.UserRequest) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("请求GET %s 失败，状态码: %d", url, resp.StatusCode)

		return fmt.Errorf("GET %s 返回非200状态码: %d |", url, resp.StatusCode)
	}

	log.Info("请求GET %s 成功，状态码: %d | %s", url, body, resp.StatusCode)

	return nil
}
