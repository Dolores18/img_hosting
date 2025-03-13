package routes

import (
	"fmt"
	"img_hosting/controllers"
	"img_hosting/docs"
	"img_hosting/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 图片托管系统 API
// @version 1.0
// @description 这是一个图片托管系统的API文档
// @host localhost:8080
// @BasePath /
// @schemes http https
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 初始化Swagger文档
	docs.SwaggerInfo.Title = "图片托管系统 API"
	docs.SwaggerInfo.Description = "这是一个图片托管系统的API文档"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// 添加Swagger UI路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("开始注册路由...")

	// 初始化控制器
	userController := controllers.NewUserController()
	tagController := controllers.NewTagController()
	imageController := controllers.NewImageController()
	authController := controllers.NewAuthController()
	privateFileController := &controllers.PrivateFileController{}
	tokenController := &controllers.TokenController{} // 取消注释，启用令牌控制器
	tokenVerifyController := controllers.NewTokenVerifyController()
	permController := controllers.NewPermissionController()

	fmt.Println("控制器初始化完成")

	// 认证相关路由（无需认证）
	authGroup := r.Group("/auth")
	{
		fmt.Println("注册用户认证路由: /auth/login, /auth/register")
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/register", authController.Register)
	}

	// 令牌验证路由
	r.GET("/api/verify-token", tokenVerifyController.VerifyToken)

	// 令牌管理路由
	tokenGroup := r.Group("/api")
	tokenGroup.Use(middleware.AuthMiddleware())
	{
		tokenGroup.POST("/tokens", tokenController.CreateToken)
		tokenGroup.GET("/tokens", tokenController.ListTokens)
		tokenGroup.DELETE("/tokens/:token", tokenController.RevokeToken)
	}

	// 图片相关路由
	imageGroup := r.Group("/images")
	imageGroup.Use(middleware.AuthMiddleware(), middleware.PermissionMiddleware())
	{
		fmt.Println("注册图片上传路由: /images/upload")
		imageGroup.POST("/upload", imageController.UploadImage)
		imageGroup.POST("/batch-upload", imageController.BatchUploadImages)
		imageGroup.GET("", imageController.ListImages)
		imageGroup.GET("/search", imageController.SearchImages)
		imageGroup.GET("/:id", imageController.GetImage)
		imageGroup.DELETE("/:id", imageController.DeleteImage)
		imageGroup.GET("/me/images", imageController.GetUserImages)

		// 添加调试日志
		fmt.Println("注册图片批量上传路由: POST /images/batch-upload")
	}

	// 标签相关路由
	tagGroup := r.Group("/tags")
	tagGroup.Use(middleware.AuthMiddleware())
	{
		tagGroup.GET("", tagController.GetAllTags)
		tagGroup.POST("", tagController.CreateTag)
		tagGroup.POST("/image/:id", tagController.AddTagToImage)
		tagGroup.GET("/:id/images", tagController.GetImagesByTag)
	}

	// 用户相关路由组
	userGroup := r.Group("/users")
	userGroup.Use(middleware.AuthMiddleware(), middleware.PermissionMiddleware())
	{
		userGroup.GET("/profile", userController.GetProfile)
		userGroup.PUT("/profile", userController.UpdateProfile)
		userGroup.GET("", userController.ListUsers)
		userGroup.DELETE("/:id", userController.DeleteUser)
		userGroup.PUT("/:id/status", userController.UpdateStatus)
		userGroup.POST("/:id/roles", userController.ManageRoles)
		userGroup.GET("/:id/roles", userController.GetRoles)
		userGroup.GET("/me/images", imageController.GetUserImages)
	}

	// 私有文件相关路由
	privateFileGroup := r.Group("/private-files")
	privateFileGroup.Use(middleware.AuthMiddleware())
	{
		privateFileGroup.POST("/upload", privateFileController.UploadFile)
		privateFileGroup.POST("/batch-upload", privateFileController.BatchUpload)
		privateFileGroup.GET("", privateFileController.ListFiles)
		privateFileGroup.GET("/:id", privateFileController.GetFile)
		privateFileGroup.DELETE("/:id", privateFileController.DeleteFile)
		privateFileGroup.PUT("/:id", privateFileController.UpdateFile)

		// 添加调试日志
		fmt.Println("注册私有文件更新路由: PUT /private-files/:id")
		fmt.Println("注册私有文件批量上传路由: POST /private-files/batch-upload")
	}

	// 权限管理路由
	permGroup := r.Group("/permissions")
	permGroup.Use(middleware.AuthMiddleware())
	permGroup.Use(middleware.PermissionMiddleware())
	{
		permGroup.GET("/all", permController.GetAllPermissions)
		permGroup.GET("/roles", permController.GetAllRoles)
		permGroup.GET("/roles/:role", permController.GetRolePermissions)
		permGroup.PUT("/roles/:role", permController.UpdateRolePermissions)
		permGroup.POST("/create", permController.CreatePermission)
		permGroup.POST("/roles/create", permController.CreateRole)
		permGroup.POST("/sync", permController.SyncConfigPermissions)
		permGroup.GET("/users/:id/permissions", permController.GetUserPermissions)
		permGroup.GET("/users/current/permissions", permController.GetCurrentUserPermissions)
	}

	fmt.Println("路由注册完成")

	return r
}
