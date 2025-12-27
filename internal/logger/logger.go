package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
	HIGHLIGHT // 高亮级别，优先级高于 INFO
)

var (
	levelStrings = map[LogLevel]string{
		DEBUG:     "调试",
		INFO:      "信息",
		WARN:      "警告",
		ERROR:     "错误",
		FATAL:     "致命",
		HIGHLIGHT: "重点",
	}

	levelColors = map[LogLevel]string{
		DEBUG:     "\033[36m",      // 青色
		INFO:      "\033[32m",      // 绿色
		WARN:      "\033[33m",      // 黄色
		ERROR:     "\033[31m",      // 红色
		FATAL:     "\033[35m",      // 洋红
		HIGHLIGHT: "\033[1;97;45m", // 粗体+亮白字+紫色背景
	}

	resetColor = "\033[0m"
	boldText   = "\033[1m"
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

	// 自动检测终端是否支持颜色
	if cfg.Colors {
		if f, ok := cfg.Output.(*os.File); ok {
			if !term.IsTerminal(int(f.Fd())) {
				cfg.Colors = false
			}
		}
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

// Highlight 记录高亮醒目的消息
func (l *Logger) Highlight(format string, v ...interface{}) {
	l.log(HIGHLIGHT, format, v...)
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

// log 是内部日志记录核心方法
func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	// 高亮消息不受普通级别限制，除非级别设定高于 HIGHLIGHT
	if level < l.minLevel && level != HIGHLIGHT {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	var builder strings.Builder

	// 1. 添加时间戳
	if l.showTime {
		builder.WriteString(time.Now().Format(l.timeFormat))
		builder.WriteString(" ")
	}

	// 2. 添加应用名称
	if l.appName != "" {
		builder.WriteString("【")
		builder.WriteString(l.appName)
		builder.WriteString("】")
		builder.WriteString(": ")
	}

	// 3. 添加带颜色的日志级别
	levelStr := levelStrings[level]
	if l.colors {
		builder.WriteString(levelColors[level])
		builder.WriteString("[")
		builder.WriteString(levelStr)
		builder.WriteString("]")
		builder.WriteString(resetColor)
	} else {
		builder.WriteString("[")
		builder.WriteString(levelStr)
		builder.WriteString("]")
	}
	builder.WriteString(" ")

	// 4. 添加调用者信息
	if l.callerInfo && (level == DEBUG || level == ERROR || level == FATAL) {
		if pc, file, line, ok := runtime.Caller(2); ok {
			funcName := runtime.FuncForPC(pc).Name()
			segments := strings.Split(file, "/")
			filename := segments[len(segments)-1]
			builder.WriteString(fmt.Sprintf("(%s:%d %s) ", filename, line, funcName))
		}
	}

	// 5. 添加消息内容
	message := fmt.Sprintf(format, v...)
	if level == HIGHLIGHT && l.colors {
		// 高亮级别的正文额外加粗
		builder.WriteString(boldText)
		builder.WriteString(message)
		builder.WriteString(resetColor)
	} else {
		builder.WriteString(message)
	}

	builder.WriteString("\n")

	// 写入输出
	fmt.Fprint(l.output, builder.String())
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

// 全局默认实例指针
var defaultLogger = Default("火尖枪")

// 全局辅助函数
func Debugf(format string, v ...interface{})     { defaultLogger.Debug(format, v...) }
func Infof(format string, v ...interface{})      { defaultLogger.Info(format, v...) }
func Warnf(format string, v ...interface{})      { defaultLogger.Warn(format, v...) }
func Errorf(format string, v ...interface{})     { defaultLogger.Error(format, v...) }
func Fatalf(format string, v ...interface{})     { defaultLogger.Fatal(format, v...) }
func Highlightf(format string, v ...interface{}) { defaultLogger.Highlight(format, v...) }

// InitDefault 重新初始化全局日志记录器
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

func StringToLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "HIGHLIGHT":
		return HIGHLIGHT
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
