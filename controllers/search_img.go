package controllers

import (
	"fmt"
	"img_hosting/pkg/logger"
	"img_hosting/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAllImageRequest struct {
	AllImg   bool   `json:"allimg"`
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"pageSize,omitempty"`
	Order    string `json:"order,omitempty"`
}

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
	log := logger.GetLogger()
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"error": "没有找到user_id"})
		return
	}

	var req GetAllImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析 JSON 数据"})
		return
	}

	log.Infof("收到的排序参数: %s", req.Order)

	if !req.AllImg {
		c.JSON(400, gin.H{"error": "json中缺少allimg字段"})
		return
	}

	// 检查是否需要分页
	enablePaging := req.Page > 0 && req.PageSize > 0

	// 如果启用分页但参数无效，使用默认值
	if enablePaging {
		if req.Page <= 0 {
			req.Page = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = 10 // 默认每页10条
		}
	}

	imgdetails, err := services.GetAllimage(user_id.(uint), enablePaging, req.Page, req.PageSize, req.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "没有找到图片"})
		return
	}

	c.JSON(200, gin.H{"data": imgdetails})
}

// 添加标签搜索请求结构体
type TagSearchRequest struct {
	Tags     []string `json:"tags"`
	Page     int      `json:"page,omitempty"`
	PageSize int      `json:"pageSize,omitempty"`
	Order    string   `json:"order,omitempty"`
}

// 修改 SearchImageByTags 函数
func SearchImageByTags(c *gin.Context) {
	log := logger.GetLogger()
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "该用户未授权"})
		return
	}

	var req TagSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析 JSON 数据"})
		log.Info("无法解析数据")
		return
	}

	// 打印日志
	fmt.Printf("收到的标签搜索请求: %+v", req)
	println(req.Order)

	// 检查是否需要分页
	enablePaging := req.Page > 0 && req.PageSize > 0

	// 如果启用分页但参数无效，使用默认值
	if enablePaging {
		if req.Page <= 0 {
			req.Page = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = 10 // 默认每页10条
		}
	}

	// 修改函数调用，添加排序参数
	imgResult, err := services.GetImagesByTags(userID.(uint), req.Tags, enablePaging, req.Page, req.PageSize, req.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索图片失败"})
		return
	}

	c.JSON(200, gin.H{
		"data": gin.H{
			"total":  imgResult.Total,
			"images": imgResult.Images,
		},
	})
}
