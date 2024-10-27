package routes

import (
	"img_hosting/controllers"
	"img_hosting/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine) {
	// 公开路由
	router.POST("/register", controllers.RegisterUser)
	router.POST("/signin", controllers.Sinngin)

	// 修改路由组配置 - 去掉斜杠前缀
	authorized := router.Group("") // 改为空字符串
	authorized.Use(middleware.AuthMiddleware())
	authorized.Use(middleware.PermissionMiddleware())
	{
		authorized.POST("/searchimg", controllers.SearchImage)
		authorized.POST("/searchAllimg", controllers.GetAllimage)
		authorized.POST("/imgupload", controllers.Uploads)
		authorized.POST("/searchbytag", controllers.SearchImageByTags)
		authorized.POST("/createtag", controllers.CreateTag)
		authorized.POST("/getalltag", controllers.GetAllTag)
		authorized.POST("/addimagetag", controllers.AddImageTag)
	}
}
