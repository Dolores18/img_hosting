// logger/logger.go
package logger

import (
	"log"
	"os"
	"path/filepath"

	"img_hosting/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *logrus.Logger

// Init 初始化日志记录器
func Init() {
	cfg := config.GetConfig()
	logger = logrus.New()

	// 设置日志格式为JSON
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 确保日志目录存在
	if err := os.MkdirAll(cfg.Log.Path, 0755); err != nil {
		log.Fatalf("创建日志目录失败: %v", err)
	}

	// 配置日志轮转
	logWriter := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Log.Path, cfg.Log.Filename),
		MaxSize:    cfg.Log.MaxSize,    // MB
		MaxBackups: cfg.Log.MaxBackups, // 文件个数
		MaxAge:     cfg.Log.MaxAge,     // 天数
		Compress:   cfg.Log.Compress,   // 是否压缩
	}

	// 设置输出
	logger.SetOutput(logWriter)

	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 添加一些默认字段
	logger.AddHook(&contextHook{})
}

// contextHook 用于添加额外的上下文信息
type contextHook struct{}

func (hook *contextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *contextHook) Fire(entry *logrus.Entry) error {
	if entry.Data == nil {
		entry.Data = make(logrus.Fields)
	}
	entry.Data["app"] = "img_hosting"
	entry.Data["environment"] = os.Getenv("APP_ENV")
	return nil
}

// GetLogger 获取日志记录器实例
func GetLogger() *logrus.Logger {
	if logger == nil {
		log.Fatal("Logger not initialized")
	}
	return logger
}
