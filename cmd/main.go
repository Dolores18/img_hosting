package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"img_hosting/config"
	"img_hosting/middleware"
	"img_hosting/models"
	"img_hosting/pkg/logger"
	"img_hosting/routes"
	"img_hosting/services"
	"time"
)

func main() {
	// 初始化数据库
	models.GetDB()

	//日志初始化
	logger.Init()
	log := logger.GetLogger()
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	services.InitValidator()
	//router.Use(middleware.RequestID())
	// 设置静态文件目录
	router.Static("/statics", "./statics")
	router.LoadHTMLGlob("statics/html/*")

	// 将根目录的请求重定向到 index.html
	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/statics/html/index.html")
	})
	// CORS middleware configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"}, // 允许的域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	config.LoadConfig() //加载配置
	addr := fmt.Sprintf(":%d", config.AppConfigInstance.App.Port)

	routes.RegisterRoutes(router)
	log.Info("Starting server on ", addr)
	router.Run(addr)
}
