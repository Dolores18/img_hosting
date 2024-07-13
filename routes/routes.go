package routes

import (
	"github.com/gin-gonic/gin"
	"img_hosting/controllers"
	"img_hosting/middleware"
)

// RegisterRoutes 注册所有路由
// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine) {
	// 公开路由
	router.POST("/register", controllers.RegisterUser)
	router.POST("/signin", controllers.Sigin) // 注意：这里我假设 "sigin" 是有意为之的拼写

	// 需要认证的路由组
	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/imgupload", controllers.Uploads)
		// 在这里添加其他需要认证的路由
	}
}
