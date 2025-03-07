package services

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"img_hosting/config"
	"img_hosting/dao"
	"img_hosting/models"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// UploadImage 处理图片上传
func UploadImage(userID uint, file *multipart.FileHeader, description string) (uint, string, error) {
	db := models.GetDB()

	// 检查文件格式和大小
	valid, message, extension, size := CheckImg(file.Filename, file.Size)
	if !valid {
		return 0, "", fmt.Errorf(message)
	}

	// 确保文件名安全
	name, _, err := SanitizeFileName(file.Filename)
	if err != nil {
		return 0, "", err
	}

	// 读取文件内容
	fileContent, err := file.Open()
	if err != nil {
		return 0, "", err
	}
	defer fileContent.Close()

	// 将文件内容读取为字节切片
	fileBytes, err := io.ReadAll(fileContent)
	if err != nil {
		return 0, "", err
	}

	// 计算文件哈希
	hashImage := HashFileName(fileBytes)

	// 检查图片是否已存在
	exists, err := dao.CheckImageExists(db, hashImage)
	if err != nil {
		return 0, "", err
	}
	if exists {
		return 0, "", fmt.Errorf("图片已存在")
	}

	// 构建文件名和路径
	hashedFilename := hashImage + extension
	uploadPath := "./statics/uploads/" + hashedFilename

	// 确保上传目录存在
	if err := os.MkdirAll("./statics/uploads/", 0755); err != nil {
		return 0, "", err
	}

	// 保存文件
	if err := saveUploadedFile(file, uploadPath); err != nil {
		return 0, "", err
	}

	// 生成缩略图
	if err := CompressToWebP(fileBytes, hashImage); err != nil {
		return 0, "", err
	}

	// 构建图片URL
	imageURL := config.AppConfigInstance.Url.Imgurl + hashImage + extension

	// 保存到数据库
	imageID, err := dao.CreateImage(db, userID, imageURL, name, extension, hashImage, size, extension)
	if err != nil {
		return 0, "", err
	}

	return imageID, imageURL, nil
}

// GetImageByID 获取图片详情
func GetImageByID(imageID uint) (*models.Image, error) {
	db := models.GetDB()
	return dao.GetImageByID(db, imageID)
}

// GetUserImages 获取用户的所有图片
func GetUserImages(userID uint, page, pageSize int) ([]models.Image, int64, error) {
	db := models.GetDB()
	return dao.GetImagesByUserID(db, userID, page, pageSize)
}

// SearchImages 搜索图片
func SearchImages(userID uint, keyword string, page, pageSize int) ([]models.Image, int64, error) {
	db := models.GetDB()
	return dao.SearchImages(db, userID, keyword, page, pageSize)
}

// DeleteImage 删除图片
func DeleteImage(imageID, userID uint) error {
	db := models.GetDB()

	// 获取图片信息
	image, err := dao.GetImageByID(db, imageID)
	if err != nil {
		return err
	}

	// 检查权限
	if image.UserID != userID {
		return fmt.Errorf("无权删除该图片")
	}

	// 删除物理文件
	filePath := "./statics/uploads/" + image.HashImage + image.Imageextenion
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// 删除缩略图
	thumbPath := "./statics/thumbnails/" + image.HashImage + ".webp"
	if err := os.Remove(thumbPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// 从数据库删除
	return dao.DeleteImage(db, imageID, userID)
}

// 辅助函数：保存上传的文件
func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// HashFileName 计算文件哈希
func HashFileName(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

// SanitizeFileName 确保文件名安全
func SanitizeFileName(filename string) (string, string, error) {
	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(filename))

	// 获取文件名（不含扩展名）
	name := strings.TrimSuffix(filepath.Base(filename), ext)

	// 移除不安全字符
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, name)

	if name == "" {
		return "", "", fmt.Errorf("无效的文件名")
	}

	return name, ext, nil
}

// CheckImg 检查图片格式和大小
func CheckImg(filename string, size int64) (bool, string, string, int64) {
	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(filename))

	// 检查文件类型
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return false, "不支持的文件类型", "", 0
	}

	// 检查文件大小（10MB限制）
	maxSize := int64(10 * 1024 * 1024)
	if size > maxSize {
		return false, "文件大小超过限制", "", 0
	}

	return true, "", ext, size
}

// CompressToWebP 将图片压缩为缩略图
func CompressToWebP(imageData []byte, hashName string) error {
	// 确保缩略图目录存在
	if err := os.MkdirAll("./statics/thumbnails/", 0755); err != nil {
		return err
	}

	// 解码图像
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return err
	}

	// 调整大小，创建缩略图
	thumbnail := imaging.Resize(img, 300, 0, imaging.Lanczos)

	// 保存缩略图
	thumbPath := "./statics/thumbnails/" + hashName + ".jpg"
	return imaging.Save(thumbnail, thumbPath)
}
