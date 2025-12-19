package download

import (
	"huo_jian_qiang/internal/constant"
	"huo_jian_qiang/internal/logger"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Downloader struct {
	DataDir string
}

var log = logger.Default(constant.APP_NAME)

func init() {
	logger.InitDefault(constant.APP_NAME, logger.INFO)
}

func PasswordDownloader(dataDir string) *Downloader {
	os.MkdirAll(dataDir, 0755)

	return &Downloader{
		DataDir: dataDir,
	}
}

func (d *Downloader) DownloadFile(url string) (string, error) {
	targetDir := filepath.Join(d.DataDir, "data")

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Error("创建目录失败: %v", err)
		return "", err
	}

	filename := filepath.Base(url)

	if !strings.HasSuffix(filename, ".txt") {
		filename = filename + ".txt"
	}

	filepath := filepath.Join(targetDir, filename)

	resp, err := http.Get(url)

	if err != nil {
		log.Error("下载文件失败: %v", err)
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("状态异常: %s", resp.Status)
		return "", err
	}

	file, err := os.Create(filepath)

	if err != nil {
		log.Error("创建文件失败: %v", err)
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		log.Error("保存文件失败: %v", err)
		return "", err
	}

	log.Info("已下载: %s -> %s\n", url, filepath)

	return filepath, nil
}
