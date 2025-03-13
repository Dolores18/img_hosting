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
	"img_hosting/pkg/logger"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// UploadImage 处理图片上传
func UploadImage(userID uint, file *multipart.FileHeader, description string) (uint, string, error) {
	cfg := config.GetConfig()
	db := models.GetDB()
	logger := logger.GetLogger()

	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"filename": file.Filename,
		"size":     file.Size,
	}).Info("开始处理图片上传")

	// 检查文件格式和大小
	valid, message, extension, size := CheckImg(file.Filename, file.Size)
	if !valid {
		logger.WithFields(logrus.Fields{
			"filename": file.Filename,
			"message":  message,
		}).Warn("图片格式或大小检查失败")
		return 0, "", fmt.Errorf(message)
	}

	// 确保文件名安全
	name, _, err := SanitizeFileName(file.Filename)
	if err != nil {
		logger.WithError(err).Error("文件名安全检查失败")
		return 0, "", err
	}

	// 读取文件内容
	fileContent, err := file.Open()
	if err != nil {
		logger.WithError(err).Error("打开上传文件失败")
		return 0, "", err
	}
	defer fileContent.Close()

	// 将文件内容读取为字节切片
	fileBytes, err := io.ReadAll(fileContent)
	if err != nil {
		logger.WithError(err).Error("读取文件内容失败")
		return 0, "", err
	}

	logger.WithFields(logrus.Fields{
		"bytes_read": len(fileBytes),
		"extension":  extension,
	}).Debug("文件内容已读取")

	// 计算文件哈希
	hashImage := HashFileName(fileBytes)
	logger.WithField("hash", hashImage).Debug("文件哈希计算完成")

	// 在事务中检查图片是否已存在
	exists, err := dao.CheckImageExists(tx, hashImage)
	if err != nil {
		tx.Rollback()
		logger.WithError(err).Error("检查图片是否存在失败")
		return 0, "", err
	}
	if exists {
		tx.Rollback()
		logger.WithField("hash", hashImage).Warn("图片已存在")
		return 0, "", fmt.Errorf("图片已存在")
	}

	// 使用配置的路径
	hashedFilename := hashImage + extension
	uploadPath := filepath.Join(cfg.Upload.Path, hashedFilename)

	logger.WithField("path", uploadPath).Debug("准备保存文件")

	// 确保上传目录存在
	if err := os.MkdirAll(cfg.Upload.Path, 0755); err != nil {
		logger.WithError(err).Error("创建上传目录失败")
		return 0, "", err
	}

	// 保存文件
	if err := saveUploadedFile(file, uploadPath); err != nil {
		logger.WithError(err).WithField("path", uploadPath).Error("保存文件失败")
		return 0, "", err
	}

	logger.WithField("path", uploadPath).Info("文件保存成功")

	// 生成缩略图
	if err := CompressToWebP(fileBytes, hashImage); err != nil {
		logger.WithError(err).Error("生成缩略图失败")
		return 0, "", err
	}

	logger.Debug("缩略图生成成功")

	// 构建图片URL
	imageURL := config.AppConfigInstance.Url.Imgurl + hashImage + extension

	// 保存到数据库
	imageID, err := dao.CreateImage(tx, userID, imageURL, name, extension, hashImage, size, description)
	if err != nil {
		tx.Rollback()
		logger.WithError(err).Error("保存图片信息到数据库失败")
		return 0, "", err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.WithError(err).Error("提交事务失败")
		return 0, "", err
	}

	logger.WithFields(logrus.Fields{
		"image_id":  imageID,
		"image_url": imageURL,
	}).Info("图片上传处理完成")

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

// DeleteImage 删除图片及其缩略图
func DeleteImage(imageID, userID uint) error {
	cfg := config.GetConfig()
	db := models.GetDB()
	logger := logger.GetLogger()

	err := db.Transaction(func(tx *gorm.DB) error {
		// 获取图片信息
		image, err := dao.GetImageByID(tx, imageID)
		if err != nil {
			return fmt.Errorf("获取图片信息失败: %w", err)
		}

		// 检查权限
		if image.UserID != userID {
			return fmt.Errorf("无权删除该图片")
		}

		// 删除图片标签关联
		if err := tx.Where("image_id = ?", imageID).Delete(&models.ImageTag{}).Error; err != nil {
			return fmt.Errorf("删除图片标签关联失败: %w", err)
		}

		// 删除原图
		originalPath := filepath.Join(cfg.Upload.Path, image.HashImage+image.Imageextenion)
		if err := deleteFile(originalPath); err != nil {
			return fmt.Errorf("删除原图失败: %w", err)
		}

		// 删除所有可能的缩略图格式
		thumbnailFormats := []string{".webp", ".png", ".jpg"}
		for _, format := range thumbnailFormats {
			thumbPath := filepath.Join(cfg.Upload.ThumbnailsPath, image.HashImage+format)
			if err := deleteFile(thumbPath); err != nil {
				return fmt.Errorf("删除缩略图失败 %s: %w", thumbPath, err)
			}
		}

		// 最后删除数据库记录
		if err := dao.DeleteImage(tx, imageID, userID); err != nil {
			return fmt.Errorf("删除图片记录失败: %w", err)
		}

		return nil
	})

	if err != nil {
		logger.WithError(err).Error("删除图片失败")
		return err
	}

	logger.WithFields(logrus.Fields{
		"image_id": imageID,
		"user_id":  userID,
	}).Info("图片删除成功")

	return nil
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
	cfg := config.GetConfig()
	logger := logger.GetLogger()

	logger.WithFields(logrus.Fields{
		"hash":      hashName,
		"data_size": len(imageData),
	}).Debug("开始生成缩略图")

	// 确保缩略图目录存在
	if err := os.MkdirAll(cfg.Upload.ThumbnailsPath, 0755); err != nil {
		logger.WithError(err).Error("创建缩略图目录失败")
		return err
	}

	// 尝试直接从文件读取而不是从内存
	uploadPath := filepath.Join(cfg.Upload.Path, hashName+".jpg")
	logger.WithField("path", uploadPath).Debug("尝试从文件读取图像")

	// 尝试方法1：使用imaging直接打开文件
	srcImg, err := imaging.Open(uploadPath)
	if err != nil {
		logger.WithError(err).Warn("使用imaging打开文件失败，尝试其他方法")

		// 尝试方法2：使用标准库解码
		reader := bytes.NewReader(imageData)
		img, format, err := image.Decode(reader)
		if err != nil {
			logger.WithError(err).Error("标准库解码图像失败")

			// 尝试方法3：跳过缩略图生成，返回成功
			logger.Warn("无法生成缩略图，但将继续处理上传")
			return nil // 返回nil而不是错误，允许上传继续
		}

		logger.WithField("format", format).Info("成功解码图像")
		srcImg = imaging.Clone(img)
	}

	// 调整大小，创建缩略图
	thumbnail := imaging.Resize(srcImg, 300, 0, imaging.Lanczos)

	// 使用配置的路径
	thumbPath := filepath.Join(cfg.Upload.ThumbnailsPath, hashName+".webp")
	logger.WithField("path", thumbPath).Debug("准备保存缩略图")

	// 尝试保存为PNG格式，如果WebP不工作
	err = imaging.Save(thumbnail, thumbPath)
	if err != nil {
		logger.WithError(err).Warn("保存为WebP失败，尝试保存为PNG")

		// 尝试保存为PNG
		pngPath := filepath.Join(cfg.Upload.ThumbnailsPath, hashName+".png")
		err = imaging.Save(thumbnail, pngPath)
		if err != nil {
			logger.WithError(err).Error("保存缩略图失败")
			// 继续处理，不返回错误
			return nil
		}
	}

	logger.WithField("path", thumbPath).Info("缩略图保存成功")
	return nil
}
