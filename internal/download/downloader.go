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

func CreateUsersTxt(dir string) error {
	commonUsers := []string{
		"admin", "user1", "testuser", "administrator", "root",
		"guest", "test", "demo", "user", "oracle",
		"postgres", "mysql", "web", "www", "ftp",
		"mail", "email", "smtp", "pop", "imap",
		"service", "system", "support", "manager", "operator",
		"monitor", "logger", "backup", "sync", "agent",
		"client", "server", "node", "master", "slave",
		"primary", "secondary", "replica", "proxy", "gateway",
		"firewall", "router", "switch", "access", "auth",
		"login", "signin", "register", "account", "profile",
		"dev", "developer", "deploy", "release", "ci",
		"cd", "build", "jenkins", "ansible", "docker",
		"kube", "kubernetes", "stack", "app", "api",
		"webapp", "frontend", "backend", "db", "database",
		"cache", "redis", "mongo", "elastic", "search",
		"analytics", "metrics", "telemetry", "trace", "audit",
		"security", "secure", "vault", "key", "token",
		"oauth", "sso", "ldap", "ad", "domain",
		"user2", "user3", "temp", "tmp", "default",
		"pi", "ubuntu", "debian", "centos", "fedora",
		"ansible", "vagrant", "ec2-user", "azureuser", "core",
	}

	dataDir := filepath.Join(dir, "downloads/data")

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(dataDir, "users.txt")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	for _, u := range commonUsers {
		if _, err := file.WriteString(u + "\n"); err != nil {
			return err
		}
	}

	log.Info("已生成用户列表: %s", filePath)
	return nil
}
