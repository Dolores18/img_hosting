package controllers

import (
	"img_hosting/models"
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

// UploadImage godoc
// @Summary 上传单个图片
// @Description 上传单个图片文件
// @Tags 图片管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "图片文件"
// @Param description formData string false "图片描述"
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /images/upload [post]
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

	c.JSON(http.StatusOK, models.ImageUploadResponse{
		ImageID:  imageID,
		ImageURL: imageURL,
	})
}

// GetImage godoc
// @Summary 获取图片详情
// @Description 获取指定图片的详细信息
// @Tags 图片管理
// @Produce json
// @Param id path int true "图片ID"
// @Security BearerAuth
// @Success 200 {object} models.Response{data=models.Image}
// @Failure 404 {object} models.Response
// @Router /images/{id} [get]
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

// ListImages godoc
// @Summary 获取图片列表
// @Description 获取用户的图片列表，支持分页
// @Tags 图片管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Security BearerAuth
// @Success 200 {object} models.Response{data=models.ImageListResponse}
// @Router /images [get]
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

// SearchImages godoc
// @Summary 搜索图片
// @Description 根据关键词搜索图片
// @Tags 图片管理
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Security BearerAuth
// @Success 200 {object} models.Response{data=models.ImageListResponse}
// @Router /images/search [get]
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

// DeleteImage godoc
// @Summary 删除图片
// @Description 删除指定的图片
// @Tags 图片管理
// @Produce json
// @Param id path int true "图片ID"
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 400,404 {object} models.Response
// @Router /images/{id} [delete]
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

// BatchUploadImages godoc
// @Summary 批量上传图片
// @Description 同时上传多个图片文件
// @Tags 图片管理
// @Accept multipart/form-data
// @Produce json
// @Param files[] formData file true "图片文件数组"
// @Param description formData string false "图片描述"
// @Security BearerAuth
// @Success 200 {object} models.BatchUploadResponse
// @Failure 400 {object} models.Response
// @Router /images/batch-upload [post]
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

	results := make([]models.ImageUploadResponse, 0)
	successCount := 0

	for _, file := range files {
		result := models.ImageUploadResponse{
			ImageID:  0,
			ImageURL: "",
		}

		// 处理上传
		imageID, imageURL, err := services.UploadImage(userID, file, description)
		if err != nil {
			result.ImageID = 0
			result.ImageURL = err.Error()
		} else {
			result.ImageID = imageID
			result.ImageURL = imageURL
			successCount++
		}

		results = append(results, result)
	}

	c.JSON(http.StatusOK, models.BatchUploadResponse{
		Message:      "批量上传完成",
		Results:      results,
		Total:        len(files),
		SuccessCount: successCount,
	})
}
