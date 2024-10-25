package main

import (
	"fmt"
	"img_hosting/config"
	"img_hosting/dao"
	"img_hosting/models"
	"img_hosting/pkg/logger"
	"img_hosting/routes"
	"img_hosting/services"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func testUserHasPermissions() {
	// 连接到 SQLite 数据库
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 获取底层的数据库连接以便之后关闭
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	defer sqlDB.Close()

	// 测试 UserHasPermissions 函数
	userID := uint(2)
	permissionNames := []string{"upload_img"}

	permissions, err := dao.UserHasPermissions(userID, permissionNames)
	if err != nil {
		log.Fatalf("Error checking user permissions: %v", err)
	}

	// 打印结果
	fmt.Printf("Permissions for user %d:\n", userID)
	for perm, has := range permissions {
		fmt.Printf("  %s: %v\n", perm, has)
	}
}

func main() {
	testUserHasPermissions()

	// 初始化数据库
	models.GetDB()

	//日志初始化
	logger.Init()
	log := logger.GetLogger()
	router := gin.Default()

	// CORS middleware configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		// 允许的域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	//router.POST("/upload", controllers.Uploads)
	//router.Use(middleware.AuthMiddleware())
	//router.Use(middleware.PermissionMiddleware())
	services.InitValidator()
	//router.Use(middleware.RequestID())

	config.LoadConfig() //加载配置
	addr := fmt.Sprintf(":%d", config.AppConfigInstance.App.Port)

	routes.RegisterRoutes(router)
	log.Info("Starting server on ", addr)
	if router.Run(addr) != nil {
		log.Fatal("Failed to start server")
	}
}
