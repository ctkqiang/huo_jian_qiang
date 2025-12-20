package download

import (
	"bytes"
	"fmt"
	"huo_jian_qiang/internal/download"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// TestPasswordDownloader 测试 PasswordDownloader 函数
func TestPasswordDownloader(t *testing.T) {
	tests := []struct {
		name    string
		dataDir string
		wantErr bool
	}{
		{
			name:    "有效的目录路径",
			dataDir: t.TempDir(),
			wantErr: false,
		},
		{
			name:    "空目录路径",
			dataDir: "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			downloader := download.PasswordDownloader(tt.dataDir)

			// 验证返回的 Downloader 不为 nil
			if downloader == nil {
				t.Error("PasswordDownloader 返回了 nil")
			}

			// 验证 DataDir 设置正确
			if downloader.DataDir != tt.dataDir {
				t.Errorf("期望的 DataDir = %v, 实际 = %v", tt.dataDir, downloader.DataDir)
			}

			// 验证目录是否被创建
			if tt.dataDir != "" {
				if _, err := os.Stat(tt.dataDir); os.IsNotExist(err) {
					t.Errorf("目录未被创建: %v", err)
				}
			}
		})
	}
}

// TestDownloader_DownloadFile 测试 DownloadFile 方法的正常流程和错误情况
func TestDownloader_DownloadFile(t *testing.T) {
	// 创建测试 HTTP 服务器
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success.txt":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test file content"))
		case "/error404.txt":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
		case "/slow.txt":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("slow content"))
		case "/binary.txt":
			w.WriteHeader(http.StatusOK)
			// 写入一些二进制数据
			w.Write([]byte{0x00, 0x01, 0x02, 0x03})
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	tests := []struct {
		name        string
		url         string
		wantErr     bool
		checkResult func(t *testing.T, filePath string, err error)
	}{
		{
			name:    "正常下载 .txt 文件",
			url:     server.URL + "/success.txt",
			wantErr: false,
			checkResult: func(t *testing.T, filePath string, err error) {
				if err != nil {
					t.Errorf("预期无错误，但得到: %v", err)
				}

				// 验证文件存在
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("文件未创建: %v", err)
				}

				// 验证文件内容
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("读取文件失败: %v", err)
				}

				expected := "test file content"
				if string(content) != expected {
					t.Errorf("文件内容不匹配，期望: %s, 实际: %s", expected, string(content))
				}

				// 验证文件名有 .txt 后缀
				if filepath.Ext(filePath) != ".txt" {
					t.Errorf("文件名缺少 .txt 后缀: %s", filePath)
				}
			},
		},
		{
			name:    "下载非 .txt 文件，自动添加后缀",
			url:     server.URL + "/data",
			wantErr: false,
			checkResult: func(t *testing.T, filePath string, err error) {
				if err != nil {
					t.Errorf("预期无错误，但得到: %v", err)
				}

				// 验证文件名自动添加了 .txt 后缀
				if filepath.Ext(filePath) != ".txt" {
					t.Errorf("文件名未自动添加 .txt 后缀: %s", filePath)
				}
			},
		},
		{
			name:    "HTTP 404 错误",
			url:     server.URL + "/error404.txt",
			wantErr: true,
			checkResult: func(t *testing.T, filePath string, err error) {
				if err == nil {
					t.Error("预期有错误，但得到 nil")
				}
			},
		},
		{
			name:    "无效的 URL",
			url:     "://invalid-url",
			wantErr: true,
			checkResult: func(t *testing.T, filePath string, err error) {
				if err == nil {
					t.Error("预期有错误（无效URL），但得到 nil")
				}
			},
		},
		{
			name:    "网络请求失败",
			url:     "http://localhost:99999/nonexistent",
			wantErr: true,
			checkResult: func(t *testing.T, filePath string, err error) {
				if err == nil {
					t.Error("预期有错误（网络失败），但得到 nil")
				}
			},
		},
		{
			name:    "下载二进制文件",
			url:     server.URL + "/binary.txt",
			wantErr: false,
			checkResult: func(t *testing.T, filePath string, err error) {
				if err != nil {
					t.Errorf("预期无错误，但得到: %v", err)
				}

				// 验证文件存在
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("文件未创建: %v", err)
				}

				// 验证文件内容
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("读取文件失败: %v", err)
				}

				expected := []byte{0x00, 0x01, 0x02, 0x03}
				if !bytes.Equal(content, expected) {
					t.Errorf("二进制文件内容不匹配")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建临时目录作为测试环境
			tempDir := t.TempDir()
			downloader := download.PasswordDownloader(tempDir)

			// 执行下载
			filePath, err := downloader.DownloadFile(tt.url)

			// 检查错误是否符合预期
			if tt.wantErr && err == nil {
				t.Error("预期有错误，但得到 nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("预期无错误，但得到: %v", err)
			}

			// 如果提供了检查函数，执行额外的验证
			if tt.checkResult != nil {
				tt.checkResult(t, filePath, err)
			}

			// 清理测试文件
			if filePath != "" {
				os.Remove(filePath)
			}
		})
	}
}

// TestDownloader_DownloadFile_NetworkErrors 测试网络错误情况
func TestDownloader_DownloadFile_NetworkErrors(t *testing.T) {
	// 测试服务器关闭的情况
	t.Run("服务器关闭", func(t *testing.T) {
		// 创建一个服务器然后立即关闭
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("content"))
		}))
		server.Close()

		tempDir := t.TempDir()
		downloader := download.PasswordDownloader(tempDir)

		_, err := downloader.DownloadFile(server.URL + "/test.txt")
		if err == nil {
			t.Error("预期有错误（服务器关闭），但得到 nil")
		}
	})

	// 测试连接超时
	t.Run("连接超时", func(t *testing.T) {
		// 创建一个永远不会响应的服务器
		server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 这里永远不会返回
		}))

		tempDir := t.TempDir()
		downloader := download.PasswordDownloader(tempDir)

		// 启动服务器但立即关闭连接
		server.Start()
		server.CloseClientConnections()
		defer server.Close()

		_, err := downloader.DownloadFile(server.URL + "/test.txt")
		if err == nil {
			t.Error("预期有错误（连接超时），但得到 nil")
		}
	})
}

// TestDownloader_DownloadFile_Validation 测试文件校验失败情况
func TestDownloader_DownloadFile_Validation(t *testing.T) {
	t.Run("空文件内容", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			// 不写入任何内容
		}))
		defer server.Close()

		tempDir := t.TempDir()
		downloader := download.PasswordDownloader(tempDir)

		filePath, err := downloader.DownloadFile(server.URL + "/empty.txt")
		if err != nil {
			t.Errorf("下载空文件不应该返回错误: %v", err)
		}

		// 验证文件存在但为空
		info, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("文件不存在: %v", err)
		}

		if info.Size() != 0 {
			t.Errorf("文件应该为空，但大小为: %d", info.Size())
		}
	})

	t.Run("大文件下载", func(t *testing.T) {
		// 创建一个大文件内容（1MB）
		largeContent := make([]byte, 1024*1024)
		for i := range largeContent {
			largeContent[i] = byte(i % 256)
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(largeContent)
		}))
		defer server.Close()

		tempDir := t.TempDir()
		downloader := download.PasswordDownloader(tempDir)

		filePath, err := downloader.DownloadFile(server.URL + "/large.txt")
		if err != nil {
			t.Errorf("下载大文件失败: %v", err)
		}

		// 验证文件大小
		info, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("文件不存在: %v", err)
		}

		if info.Size() != int64(len(largeContent)) {
			t.Errorf("文件大小不匹配，期望: %d, 实际: %d", len(largeContent), info.Size())
		}
	})
}

// TestDownloader_ConcurrentDownload 测试并发下载场景
func TestDownloader_ConcurrentDownload(t *testing.T) {
	// 创建测试服务器
	requestCount := 0
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		requestID := requestCount
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("content for request %d", requestID)))
	}))
	defer server.Close()

	tempDir := t.TempDir()
	downloader := download.PasswordDownloader(tempDir)

	// 并发下载次数
	concurrentCount := 10
	var wg sync.WaitGroup
	errors := make(chan error, concurrentCount)
	filePaths := make(chan string, concurrentCount)

	for i := 0; i < concurrentCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			filePath, err := downloader.DownloadFile(fmt.Sprintf("%s/file%d.txt", server.URL, id))
			if err != nil {
				errors <- err
				return
			}

			filePaths <- filePath
		}(i)
	}

	wg.Wait()
	close(errors)
	close(filePaths)

	// 检查是否有错误
	for err := range errors {
		t.Errorf("并发下载出现错误: %v", err)
	}

	// 验证所有文件都已创建
	fileCount := 0
	for filePath := range filePaths {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("文件未创建: %s", filePath)
		} else {
			fileCount++
		}
	}

	if fileCount != concurrentCount {
		t.Errorf("期望下载 %d 个文件，实际下载 %d 个", concurrentCount, fileCount)
	}

	// 验证请求计数
	if requestCount != concurrentCount {
		t.Errorf("期望 %d 个请求，实际 %d 个请求", concurrentCount, requestCount)
	}
}

// TestCreateUsersTxt 测试 CreateUsersTxt 函数
func TestCreateUsersTxt(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		wantErr bool
	}{
		{
			name:    "有效目录",
			dir:     t.TempDir(),
			wantErr: false,
		},
		{
			name:    "当前目录",
			dir:     ".",
			wantErr: false,
		},
		{
			name:    "空目录",
			dir:     "",
			wantErr: false,
		},
		{
			name:    "嵌套目录",
			dir:     filepath.Join(t.TempDir(), "nested", "deep"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := download.CreateUsersTxt(tt.dir)

			if tt.wantErr && err == nil {
				t.Error("预期有错误，但得到 nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("预期无错误，但得到: %v", err)
			}

			// 验证文件是否被创建
			if err == nil {
				expectedPath := filepath.Join(tt.dir, "data", "users.txt")
				if tt.dir == "" {
					expectedPath = filepath.Join("data", "users.txt")
				}

				// 检查文件是否存在
				if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
					t.Errorf("用户文件未创建: %s", expectedPath)
				}

				// 验证文件内容
				content, err := os.ReadFile(expectedPath)
				if err != nil {
					t.Errorf("读取用户文件失败: %v", err)
				}

				// 验证文件不为空
				if len(content) == 0 {
					t.Error("用户文件内容为空")
				}

				// 验证包含预期的用户名
				if !bytes.Contains(content, []byte("admin")) {
					t.Error("用户文件不包含 'admin'")
				}

				if !bytes.Contains(content, []byte("root")) {
					t.Error("用户文件不包含 'root'")
				}

				// 统计行数（应该是100个用户）
				lines := bytes.Count(content, []byte("\n"))
				expectedLines := 100
				if lines != expectedLines {
					t.Errorf("用户数量不匹配，期望: %d, 实际: %d", expectedLines, lines)
				}
			}
		})
	}
}

// TestCreateUsersTxt_ErrorCases 测试 CreateUsersTxt 的错误场景
func TestCreateUsersTxt_ErrorCases(t *testing.T) {
	t.Run("目录权限错误", func(t *testing.T) {
		// 尝试创建到只读目录的路径
		// 注意：这个测试可能在某些环境下失败或跳过
		if os.Getenv("CI") != "true" {
			t.Skip("跳过权限测试（非 CI 环境）")
		}

		err := download.CreateUsersTxt("/proc/self/readonly")
		if err == nil {
			t.Error("预期有错误（权限问题），但得到 nil")
		}
	})

	t.Run("文件写入错误", func(t *testing.T) {
		tempDir := t.TempDir()
		dataDir := filepath.Join(tempDir, "data")

		// 创建 data 目录，但不给写入权限
		os.MkdirAll(dataDir, 0555) // 只读权限

		// 尝试创建 users.txt 文件
		filePath := filepath.Join(dataDir, "users.txt")
		_, err := os.Create(filePath)
		if err == nil {
			// 如果文件创建成功，删除它并恢复权限
			os.Remove(filePath)
			os.Chmod(dataDir, 0755)
			t.Skip("测试环境允许创建文件，跳过此测试")
		}
	})
}

// TestIntegration 集成测试
func TestIntegration(t *testing.T) {
	t.Run("完整流程测试", func(t *testing.T) {
		tempDir := t.TempDir()

		// 1. 创建下载器
		downloader := download.PasswordDownloader(tempDir)
		if downloader == nil {
			t.Fatal("PasswordDownloader 返回 nil")
		}

		// 2. 创建测试 HTTP 服务器
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("integration test content"))
		}))
		defer server.Close()

		// 3. 下载文件
		downloadedFile, err := downloader.DownloadFile(server.URL + "/integration.txt")
		if err != nil {
			t.Fatalf("下载失败: %v", err)
		}

		// 4. 验证下载的文件
		if _, err := os.Stat(downloadedFile); os.IsNotExist(err) {
			t.Fatalf("下载的文件不存在: %s", downloadedFile)
		}

		// 5. 创建用户文件
		err = download.CreateUsersTxt(tempDir)
		if err != nil {
			t.Fatalf("创建用户文件失败: %v", err)
		}

		// 6. 验证用户文件
		usersFile := filepath.Join(tempDir, "data", "users.txt")
		if _, err := os.Stat(usersFile); os.IsNotExist(err) {
			t.Fatalf("用户文件不存在: %s", usersFile)
		}

		t.Logf("集成测试完成。下载文件: %s, 用户文件: %s", downloadedFile, usersFile)
	})
}

// 测试初始化函数
func TestMain(m *testing.M) {
	// 测试前的设置
	fmt.Println("开始下载器单元测试...")

	// 运行测试
	code := m.Run()

	// 测试后的清理
	fmt.Println("下载器单元测试完成")

	os.Exit(code)
}
