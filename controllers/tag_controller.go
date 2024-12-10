package controllers

import (
	"img_hosting/services"

	"github.com/gin-gonic/gin"
)

// 统一错误信息
const (
	ErrInvalidJSON = "无法解析 JSON 数据"
	ErrNoTagName   = "缺少标签名称"
	ErrCreateTag   = "创建标签失败"
	ErrGetTags     = "获取标签失败"
)

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	var req struct {
		TagName string `json:"tagname" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": ErrInvalidJSON})
		return
	}
	userID := c.GetUint("user_id")
	if tagID, err := services.CreateTag(userID, req.TagName); err != nil {
		c.JSON(500, gin.H{"error": ErrCreateTag})
		return
	} else {
		c.JSON(200, gin.H{"message": "标签创建成功", "tagid": tagID})
	}
}

// GetAllTag 获取所有标签
func GetAllTag(c *gin.Context) {
	userID := c.GetUint("user_id")
	tags, err := services.GetAllTag(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": ErrGetTags})
		return
	}
	println("tags: %v", tags)
	c.JSON(200, gin.H{"data": tags})
}
