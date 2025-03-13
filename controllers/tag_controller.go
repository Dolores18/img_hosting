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

// GetAllTags godoc
// @Summary 获取所有标签
// @Description 获取用户的所有标签列表
// @Tags 标签管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response{data=[]models.Tag}
// @Failure 500 {object} models.Response
// @Router /tags [get]
func (tc *TagController) GetAllTags(c *gin.Context) {
	userID := c.GetUint("user_id")

	tags, err := services.GetAllTag(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	TagName string `json:"tag_name" binding:"required"`
}

// AddTagRequest 添加标签请求
type AddTagRequest struct {
	TagID uint `json:"tag_id" binding:"required"`
}

// CreateTag godoc
// @Summary 创建标签
// @Description 创建新的标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param request body CreateTagRequest true "标签名称"
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 400,500 {object} models.Response
// @Router /tags [post]
func (tc *TagController) CreateTag(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateTagRequest

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

// AddTagToImage godoc
// @Summary 为图片添加标签
// @Description 为指定图片添加标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path int true "图片ID"
// @Param request body AddTagRequest true "标签ID"
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 400,500 {object} models.Response
// @Router /tags/image/{id} [post]
func (tc *TagController) AddTagToImage(c *gin.Context) {
	userID := c.GetUint("user_id")
	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	var req AddTagRequest

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

// GetImagesByTag godoc
// @Summary 获取标签下的图片
// @Description 获取指定标签下的所有图片
// @Tags 标签管理
// @Produce json
// @Param id path int true "标签ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Security BearerAuth
// @Success 200 {object} models.Response{data=models.ImageListResponse}
// @Failure 400,500 {object} models.Response
// @Router /tags/{id}/images [get]
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
