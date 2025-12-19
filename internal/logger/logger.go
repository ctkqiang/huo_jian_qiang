package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	levelStrings = map[LogLevel]string{
		DEBUG: "调试",
		INFO:  "信息",
		WARN:  "警告",
		ERROR: "错误",
		FATAL: "致命",
	}

	levelColors = map[LogLevel]string{
		DEBUG: "\033[36m", // 青色
		INFO:  "\033[32m", // 绿色
		WARN:  "\033[33m", // 黄色
		ERROR: "\033[31m", // 红色
		FATAL: "\033[35m", // 洋红
	}

	resetColor = "\033[0m"
)

type Logger struct {
	mu         sync.Mutex
	appName    string
	minLevel   LogLevel
	output     io.Writer
	showTime   bool
	callerInfo bool
	colors     bool
	timeFormat string
}

type Config struct {
	AppName    string
	MinLevel   LogLevel
	Output     io.Writer
	ShowTime   bool
	CallerInfo bool
	Colors     bool
	TimeFormat string
}

// New 创建一个新的日志记录器
func New(cfg Config) *Logger {
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = "2006-01-02 15:04:05"
	}

	return &Logger{
		appName:    cfg.AppName,
		minLevel:   cfg.MinLevel,
		output:     cfg.Output,
		showTime:   cfg.ShowTime,
		callerInfo: cfg.CallerInfo,
		colors:     cfg.Colors,
		timeFormat: cfg.TimeFormat,
	}
}

// Default 使用合理的默认值创建日志记录器
func Default(appName string) *Logger {
	return New(Config{
		AppName:    appName,
		MinLevel:   INFO,
		Output:     os.Stdout,
		ShowTime:   true,
		CallerInfo: false,
		Colors:     true,
		TimeFormat: "15:04:05",
	})
}

// Debug 记录调试消息
func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(DEBUG, format, v...)
}

// Info 记录信息消息
func (l *Logger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

// Warn 记录警告消息
func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(WARN, format, v...)
}

// Error 记录错误消息
func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}

// Fatal 记录致命消息并退出
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.log(FATAL, format, v...)
	os.Exit(1)
}

// log 是内部日志记录方法
func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	if level < l.minLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	var builder strings.Builder

	// 添加时间戳
	if l.showTime {
		builder.WriteString(time.Now().Format(l.timeFormat))
		builder.WriteString(" ")
	}

	// 添加应用名称
	if l.appName != "" {
		builder.WriteString(l.appName)
		builder.WriteString(" ")
	}

	// 添加带颜色的日志级别
	levelStr := levelStrings[level]
	if l.colors && levelColors[level] != "" {
		builder.WriteString(l.colorize(levelStr, level))
	} else {
		builder.WriteString("[")
		builder.WriteString(levelStr)
		builder.WriteString("]")
	}
	builder.WriteString(" ")

	// 添加调用者信息（文件:行号）
	if l.callerInfo && (level == DEBUG || level == ERROR || level == FATAL) {
		if pc, file, line, ok := runtime.Caller(2); ok {
			funcName := runtime.FuncForPC(pc).Name()
			// 仅提取文件名，不包含完整路径
			segments := strings.Split(file, "/")
			filename := segments[len(segments)-1]
			builder.WriteString(fmt.Sprintf("(%s:%d %s) ", filename, line, funcName))
		}
	}

	// 添加消息
	message := fmt.Sprintf(format, v...)
	builder.WriteString(message)

	// 添加换行符
	builder.WriteString("\n")

	// 写入输出
	fmt.Fprint(l.output, builder.String())
}

// colorize 为日志级别字符串添加颜色
func (l *Logger) colorize(text string, level LogLevel) string {
	if !l.colors {
		return "[" + text + "]"
	}
	return levelColors[level] + "[" + text + "]" + resetColor
}

// SetLevel 修改最低日志级别
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = level
}

// SetOutput 修改输出写入器
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
}

// DisableColors 关闭颜色输出
func (l *Logger) DisableColors() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.colors = false
}

// EnableColors 开启颜色输出
func (l *Logger) EnableColors() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.colors = true
}

// 一行初始化辅助函数
func Debugf(format string, v ...interface{}) {
	defaultLogger.Debug(format, v...)
}

func Infof(format string, v ...interface{}) {
	defaultLogger.Info(format, v...)
}

func Warnf(format string, v ...interface{}) {
	defaultLogger.Warn(format, v...)
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.Error(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatal(format, v...)
}

// 全局默认日志记录器
var defaultLogger = Default("APP")

// InitDefault 设置全局日志记录器
func InitDefault(appName string, minLevel LogLevel) {
	defaultLogger = New(Config{
		AppName:    appName,
		MinLevel:   minLevel,
		Output:     os.Stdout,
		ShowTime:   true,
		CallerInfo: false,
		Colors:     true,
		TimeFormat: "15:04:05",
	})
}

// StringToLevel 将字符串转换为日志级别
func StringToLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN", "WARNING":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return INFO
	}
}

// LevelToString 将日志级别转换为字符串
func LevelToString(level LogLevel) string {
	if str, ok := levelStrings[level]; ok {
		return str
	}
	return "未知"
}
