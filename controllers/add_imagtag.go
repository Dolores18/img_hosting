package controllers

import (
	"img_hosting/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddImageTag(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析 JSON 数据"})
		return
	}

	imageID := uint(jsonData["imageid"].(float64))
	tagname := jsonData["tagname"].(string)
	userID := c.GetUint("user_id")

	// 在创建标签之前添加日志
	log.Printf("准备创建标签，userID: %d, tagname: %s", userID, tagname)
	tagID, err := services.CreateTag(userID, tagname)
	if err != nil {
		log.Printf("创建标签失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败"})
		return
	}

	// 创建图片和标签的关联
	if err := services.CreateImageTag(imageID, tagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "关联图片和标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "标签添加成功"})
}
