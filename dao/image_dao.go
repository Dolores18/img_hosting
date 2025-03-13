package dao

import (
	"fmt"
	"img_hosting/models"

	"gorm.io/gorm"
)

// CreateImage 创建新图片记录
func CreateImage(db *gorm.DB, userID uint, imageURL, imageName, imageExtension, hashImage string, imageSize int64, imageType string) (uint, error) {
	image := models.Image{
		UserID:        userID,
		ImageURL:      imageURL,
		ImageName:     imageName,
		Imageextenion: imageExtension,
		HashImage:     hashImage,
		ImageSize:     imageSize,
		ImageType:     imageType,
	}

	result := db.Create(&image)
	if result.Error != nil {
		fmt.Printf("创建图片记录失败: %v\n", result.Error)
		return 0, result.Error
	}

	return image.ImageID, nil
}

// GetImageByID 根据ID获取图片
func GetImageByID(db *gorm.DB, imageID uint) (*models.Image, error) {
	var image models.Image
	result := db.Preload("Tags").First(&image, imageID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &image, nil
}

// GetImagesByUserID 获取用户的所有图片
func GetImagesByUserID(db *gorm.DB, userID uint, page, pageSize int) ([]models.Image, int64, error) {
	var images []models.Image
	var total int64

	// 获取总数
	if err := db.Model(&models.Image{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	result := db.Where("user_id = ?", userID).
		Preload("Tags").
		Offset(offset).
		Limit(pageSize).
		Find(&images)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return images, total, nil
}

// SearchImages 搜索图片
func SearchImages(db *gorm.DB, userID uint, keyword string, page, pageSize int) ([]models.Image, int64, error) {
	var images []models.Image
	var total int64

	query := db.Model(&models.Image{}).Where("user_id = ?", userID)

	// 如果有关键词，添加搜索条件
	if keyword != "" {
		query = query.Where("image_name LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	result := query.Preload("Tags").
		Offset(offset).
		Limit(pageSize).
		Find(&images)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return images, total, nil
}

// DeleteImage 删除图片
func DeleteImage(db *gorm.DB, imageID, userID uint) error {
	// 使用单个事务处理整个删除过程
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 先查询图片是否存在
		var image models.Image
		if err := tx.Where("image_id = ? AND user_id = ?", imageID, userID).First(&image).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("图片不存在或无权限删除")
			}
			return err
		}

		// 2. 删除图片标签关联
		if err := tx.Where("image_id = ?", imageID).Delete(&models.ImageTag{}).Error; err != nil {
			return err
		}

		// 3. 删除图片记录
		if err := tx.Delete(&image).Error; err != nil {
			return err
		}

		return nil
	})
}

// CheckImageExists 检查图片是否已存在
func CheckImageExists(db *gorm.DB, hashImage string) (bool, error) {
	var count int64
	err := db.Model(&models.Image{}).
		Where("hash_image = ?", hashImage).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("检查图片哈希失败: %w", err)
	}

	return count > 0, nil
}

// ListAllImages 获取所有图片（管理员用）
func ListAllImages(db *gorm.DB, page, pageSize int) ([]models.Image, int64, error) {
	var images []models.Image
	var total int64

	if err := db.Model(&models.Image{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	result := db.Preload("Tags").
		Offset(offset).
		Limit(pageSize).
		Find(&images)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return images, total, nil
}
