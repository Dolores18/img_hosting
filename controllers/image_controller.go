package controllers

import (
	"img_hosting/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ImageController 图片控制器
type ImageController struct{}

// NewImageController 创建图片控制器
func NewImageController() *ImageController {
	return &ImageController{}
}

// UploadImage 处理图片上传
func (ic *ImageController) UploadImage(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的图片"})
		return
	}

	// 获取描述信息
	description := c.PostForm("description")

	// 处理上传
	imageID, imageURL, err := services.UploadImage(userID, file, description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "图片上传成功",
		"data": gin.H{
			"image_id":  imageID,
			"image_url": imageURL,
		},
	})
}

// GetImage 获取图片详情
func (ic *ImageController) GetImage(c *gin.Context) {
	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	image, err := services.GetImageByID(uint(imageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"image": image})
}

// ListImages 获取用户图片列表
func (ic *ImageController) ListImages(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	images, total, err := services.GetUserImages(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取图片列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"images":    images,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// SearchImages 搜索图片
func (ic *ImageController) SearchImages(c *gin.Context) {
	userID := c.GetUint("user_id")
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	images, total, err := services.SearchImages(userID, keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索图片失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"images":    images,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// DeleteImage 删除图片
func (ic *ImageController) DeleteImage(c *gin.Context) {
	userID := c.GetUint("user_id")
	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片ID"})
		return
	}

	if err := services.DeleteImage(uint(imageID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "图片已删除"})
}

// BatchUploadImages 批量上传图片
func (ic *ImageController) BatchUploadImages(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 获取表单中的多个文件
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析表单数据"})
		return
	}

	files := form.File["files[]"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的图片"})
		return
	}

	// 获取描述信息
	description := c.PostForm("description")

	results := make([]map[string]interface{}, 0)
	successCount := 0

	for _, file := range files {
		result := map[string]interface{}{
			"file_name": file.Filename,
			"success":   false,
		}

		// 处理上传
		imageID, imageURL, err := services.UploadImage(userID, file, description)
		if err != nil {
			result["error"] = err.Error()
		} else {
			result["success"] = true
			result["image_id"] = imageID
			result["image_url"] = imageURL
			successCount++
		}

		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "批量上传完成",
		"results":       results,
		"total":         len(files),
		"success_count": successCount,
	})
}
