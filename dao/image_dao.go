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
	// 先检查图片是否属于该用户
	var count int64
	if err := db.Model(&models.Image{}).Where("image_id = ? AND user_id = ?", imageID, userID).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("图片不存在或无权限删除")
	}

	// 删除图片标签关联
	if err := db.Where("image_id = ?", imageID).Delete(&models.ImageTag{}).Error; err != nil {
		return err
	}

	// 删除图片
	return db.Delete(&models.Image{}, imageID).Error
}

// CheckImageExists 检查图片是否存在
func CheckImageExists(db *gorm.DB, hashImage string) (bool, error) {
	var count int64
	err := db.Model(&models.Image{}).Where("hash_image = ?", hashImage).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
