package controllers

import (
	"fmt"
	"img_hosting/dao"
	"img_hosting/models"
	"img_hosting/services"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PrivateFileController struct{}

// UploadFile 上传私人文件
func (pfc *PrivateFileController) UploadFile(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}

	// 获取加密选项
	isEncrypted := c.PostForm("is_encrypted") == "true"
	password := c.PostForm("password")

	// 如果设置了加密但没有提供密码
	if isEncrypted && password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "加密文件必须提供密码"})
		return
	}

	// 上传文件
	privateFile, err := services.UploadPrivateFile(file, userID, isEncrypted, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文件上传成功",
		"file":    privateFile,
	})
}

// GetFile 获取文件信息
func (pfc *PrivateFileController) GetFile(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件ID"})
		return
	}

	password := c.Query("password")

	file, err := services.GetPrivateFile(uint(fileID), userID, password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file": file})
}

// DownloadFile 下载文件
func (pfc *PrivateFileController) DownloadFile(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件ID"})
		return
	}

	// 获取密码（如果提供）
	password := c.Query("password")

	// 获取解密后的文件路径
	filePath, err := services.GetDecryptedFilePath(uint(fileID), userID, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取文件信息"})
		return
	}

	// 获取文件名
	db := models.GetDB()
	file, err := dao.GetPrivateFileByID(db, uint(fileID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 设置文件名和内容类型
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.QueryEscape(file.FileName)))
	c.Header("Content-Type", file.FileType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// 提供文件下载
	c.File(filePath)

	// 如果是临时解密文件，下载后删除
	if file.IsEncrypted {
		defer os.Remove(filePath)
	}
}

// ListFiles 获取文件列表
func (pfc *PrivateFileController) ListFiles(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	files, total, err := services.ListPrivateFiles(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files":     files,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// DeleteFile 删除文件
func (pfc *PrivateFileController) DeleteFile(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件ID"})
		return
	}

	if err := services.DeletePrivateFile(uint(fileID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件已删除"})
}

// SearchFiles 搜索文件
func (pfc *PrivateFileController) SearchFiles(c *gin.Context) {
	userID := c.GetUint("user_id")
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	files, total, err := services.SearchPrivateFiles(userID, keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files":     files,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateFile 更新文件信息
func (pfc *PrivateFileController) UpdateFile(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件ID"})
		return
	}

	// 定义请求体结构
	type UpdateFileRequest struct {
		FileName    string `json:"file_name"`
		IsEncrypted bool   `json:"is_encrypted"`
		Password    string `json:"password"`
	}

	var req UpdateFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 更新文件信息
	updatedFile, err := services.UpdatePrivateFileInfo(uint(fileID), userID, req.FileName, req.IsEncrypted, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文件信息已更新",
		"file":    updatedFile,
	})
}
