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

// @title 图片托管系统 API
// @version 1.0
// @description 图片托管系统的 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API 支持
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 请输入 'Bearer {token}' 格式的认证信息

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
	config.LoadConfig() // 确保这行在使用 JWT 之前执行
	testUserHasPermissions()
	addr := fmt.Sprintf(":%d", config.AppConfigInstance.App.Port)

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
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	//router.POST("/upload", controllers.Uploads)
	//router.Use(middleware.AuthMiddleware())
	//router.Use(middleware.PermissionMiddleware())
	services.InitValidator()
	//router.Use(middleware.RequestID())

	router = routes.SetupRouter()
	log.Info("Starting server on ", addr)
	if router.Run(addr) != nil {
		log.Fatal("Failed to start server")
	}

}
