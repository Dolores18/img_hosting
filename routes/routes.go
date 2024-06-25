package routes

import (
	"github.com/gin-gonic/gin"
	"img_hosting/controllers"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine) {
	router.POST("/register", controllers.RegisterUser)
	router.POST("/sigin", controllers.Sigin)
}
