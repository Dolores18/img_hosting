package controllers

import (
	"fmt"
	"img_hosting/pkg/logger"
	"img_hosting/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchImage(c *gin.Context) {
	log := logger.GetLogger() //必须实例化
	// 添加调试日志
	for key, value := range c.Keys {
		log.Infof("Context key: %s, value: %v", key, value)
	}
	//获取用户信息
	user_id, exists := c.Get("user_id")
	log.Infof("user_id exists: %v, value: %v", exists, user_id)
	if !exists {
		c.JSON(400, gin.H{"error": "没有找到user_id"})
		log.Info("没有找到user_id")

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
	log.Info("用户请求的图片名称是: ", imgname)
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

	// nameStr 包含了 JSON 中的 name 值
	fmt.Printf("Name: %s\n", nameStr)

	img_name := nameStr

	println(user_id, img_name)

	imgdetail, err := services.Getimage(user_id.(uint), img_name)
	println(imgdetail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "没有找到图片"})
		return
	}
	c.JSON(200, gin.H{"data": imgdetail})
}

// 获取上传图片的全部信息
func GetAllimage(c *gin.Context) {
	//log := logger.GetLogger()
	user_id, exists := c.Get("user_id")
	println("handler获得用户id:", uint64(user_id.(uint))) // 正确格式化输出
	if !exists {
		c.JSON(400, gin.H{"error": "没有找到user_id"})
		return
	}

	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析 JSON 数据"})
		return
	}

	allimg, exists := jsonData["allimg"]
	if !exists {
		c.JSON(400, gin.H{"error": "json中缺少allimg字段"})
		return // 添加 return
	}

	allimgbool, ok := allimg.(bool)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "allimg字段不是bool类型"})
		return
	}

	// 只有当 allimg 为 true 时才获取图片
	if !allimgbool {
		c.JSON(400, gin.H{"error": "allimg必须为true"})
		return
	}

	imgdetails, err := services.GetAllimage(user_id.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "没有找到图片"})
		return
	}

	c.JSON(200, gin.H{"data": imgdetails})
}

// 根据标签搜索图片
func SearchImageByTags(c *gin.Context) {
	log := logger.GetLogger()
	// 修改key为 "user_id"
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "该用户未授权"})
		return
	}

	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析 JSON 数据"})
		log.Info("无法解析数据")
		return
	}

	tags, exists := jsonData["tags"].([]interface{})
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON 中缺少 tags 字段或格式错误"})
		return
	}

	// 转换 tags 为 []string
	tagsList := make([]string, len(tags))
	for i, tag := range tags {
		if tagStr, ok := tag.(string); ok {
			tagsList[i] = tagStr
		}
	}

	// 使用中间件已验证的 userID
	imgdetails, err := services.GetImagesByTags(userID.(uint), tagsList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索图片失败"})
		return
	}

	c.JSON(200, gin.H{"data": imgdetails})
}
