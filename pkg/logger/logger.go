// logger/logger.go
package logger

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// Init 初始化日志记录器
func Init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	// 创建或打开日志文件
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logger.SetOutput(file)
	logger.SetLevel(logrus.InfoLevel)
}

// GetLogger 获取日志记录器实例
func GetLogger() *logrus.Logger {
	if logger == nil {
		log.Fatal("Logger not initialized")
	}
	return logger
}
