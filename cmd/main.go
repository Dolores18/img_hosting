package main

import (
	"fmt"
	_ "fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"img_hosting/config"
	_ "img_hosting/config"
	_ "img_hosting/middleware"
	"img_hosting/pkg/logger"
	_ "img_hosting/pkg/logger"
	"img_hosting/routes"
	_ "img_hosting/routes"
	"img_hosting/services"
)

func main() {

	logger.Init()
	log := logger.GetLogger()
	router := gin.Default()
	services.InitValidator()
	//router.Use(middleware.RequestID())

	config.LoadConfig() //加载配置
	addr := fmt.Sprintf(":%d", config.AppConfigInstance.App.Port)

	routes.RegisterRoutes(router)
	log.Info("Starting server on ", addr)
	router.Run(addr)
}
