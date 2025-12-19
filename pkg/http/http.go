package http

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"huo_jian_qiang/internal/model"
	"huo_jian_qiang/pkg/logger"
)

var log = logger.New("debug")
var userRequest model.UserRequest

func GET(uri string, userRequest model.UserRequest) error {

	basedURL, err := url.Parse(uri)

	param := url.Values{}
	resp, err := http.Get(basedURL.String() + "?" + param.Encode())

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("请求GET %s 失败，状态码: %d", basedURL, resp.StatusCode)

		return fmt.Errorf("GET %s 返回非200状态码: %d |", basedURL, resp.StatusCode)
	}

	log.Info("请求GET %s 成功，状态码: %d | %s", basedURL, body, resp.StatusCode)

	return nil
}
