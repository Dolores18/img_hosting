package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"img_hosting/middleware"
	"img_hosting/pkg/logger"
	"img_hosting/services"
	"net/http"
)

func SearchImage(c *gin.Context) {
	log := logger.GetLogger() //必须实例化
	//获取用户信息
	claims, err := middleware.ParseAndValidateToken(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Info("没有找到claims")

		return
	}
	//请求得到json数据

	var jsonData map[string]interface{}

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析 JSON 数据"})
		log.Info("无法解析数据")
		return
	}

	// 优雅地提取 name 字段
	imgname, exists := jsonData["name"]
	log.Info("用户请求的图片名称是: %s", imgname)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON 中缺少 name 字段"})
		return
	}

	// 将 interface{} 转换为 string
	nameStr, ok := imgname.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name 字段不是字符串类型"})
		return
	}

	// 现在 nameStr 包含了 JSON 中的 name 值
	fmt.Printf("Name: %s\n", nameStr)

	img_name := nameStr
	user_id := claims.UserID
	println(user_id, img_name)

	imgdetail, err := services.Getimage(user_id, img_name)
	println(imgdetail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "没有找到图片"})
		return
	}
	c.JSON(200, gin.H{"data": imgdetail})
}
