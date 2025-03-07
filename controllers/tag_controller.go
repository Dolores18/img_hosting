package controllers

import (
	"img_hosting/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TagController 标签控制器
type TagController struct{}

// NewTagController 创建标签控制器
func NewTagController() *TagController {
	return &TagController{}
}

// GetAllTags 获取所有标签
func (tc *TagController) GetAllTags(c *gin.Context) {
	userID := c.GetUint("user_id")

	tags, err := services.GetAllTag(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

// CreateTag 创建标签
func (tc *TagController) CreateTag(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		TagName string `json:"tag_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	tagID, err := services.CreateTag(userID, req.TagName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "标签创建成功",
		"tag_id":  tagID,
	})
}

// AddTagToImage 为图片添加标签
func (tc *TagController) AddTagToImage(c *gin.Context) {
	userID := c.GetUint("user_id")
	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	var req struct {
		TagID uint `json:"tag_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 添加权限检查
	if err := services.AddTagToImage(uint(imageID), req.TagID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "标签添加成功"})
}

// GetImagesByTag 获取带有特定标签的图片
func (tc *TagController) GetImagesByTag(c *gin.Context) {
	userID := c.GetUint("user_id")
	tagID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	images, total, err := services.GetImagesByTag(userID, uint(tagID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取图片失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"images":    images,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
