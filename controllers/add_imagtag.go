package controllers

import (
	"img_hosting/pkg/logger"
	"img_hosting/services"

	"net/http"

	"github.com/gin-gonic/gin"
)

// AddImageTagRequest 添加标签请求
type AddImageTagRequest struct {
	ImageID  uint     `json:"imageid"`
	TagNames []string `json:"tagnames"`
}

// AddImageTag godoc
// @Summary 为图片添加多个标签
// @Description 为指定图片添加多个标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param request body AddImageTagRequest true "添加标签请求"
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "请求数据无效"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /images/tags [post]
func AddImageTag(c *gin.Context) {
	log := logger.GetLogger()
	log.Info("开始处理添加标签请求")

	var req AddImageTagRequest

	// 直接使用 ShouldBindJSON
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("解析请求数据失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析请求数据"})
		return
	}

	// 使用 Printf 打印接收到的数据
	log.Printf("接收到的数据 - ImageID: %d", req.ImageID)
	log.Printf("接收到的标签列表: %+v", req.TagNames)
	// 更详细的标签信息打印
	for i, tag := range req.TagNames {
		log.Printf("标签[%d]: %s", i, tag)
	}

	if len(req.TagNames) == 0 {
		log.Warn("标签列表为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "标签列表不能为空"})
		return
	}

	if err := services.CreateImageTags(req.ImageID, c.GetUint("user_id"), req.TagNames); err != nil {
		log.Error("添加标签失败: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加标签失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "标签添加成功",
		"tags":    req.TagNames, // 返回添加的标签列表
	})
}
