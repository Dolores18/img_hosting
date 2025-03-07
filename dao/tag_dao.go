package dao

import (
	"fmt"
	"img_hosting/models"

	"gorm.io/gorm"
)

// CreateTag 创建新标签
func CreateTag(db *gorm.DB, userID uint, tagName string) (*models.Tag, error) {
	tag := models.Tag{
		UserID:  userID,
		TagName: tagName,
	}

	// 检查标签是否已存在
	var existingTag models.Tag
	result := db.Where("user_id = ? AND tag_name = ?", userID, tagName).First(&existingTag)
	if result.Error == nil {
		// 标签已存在，返回已存在的标签
		return &existingTag, nil
	} else if result.Error != gorm.ErrRecordNotFound {
		// 查询出错
		return nil, result.Error
	}

	// 创建新标签
	result = db.Create(&tag)
	if result.Error != nil {
		fmt.Printf("创建标签失败: %v\n", result.Error)
		return nil, result.Error
	}

	return &tag, nil
}

// GetAllTags 获取用户的所有标签
func GetAllTags(db *gorm.DB, userID uint) ([]models.Tag, error) {
	var tags []models.Tag
	result := db.Where("user_id = ?", userID).Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}
	return tags, nil
}

// AddImageTag 为图片添加标签
func AddImageTag(db *gorm.DB, imageID, tagID uint) error {
	// 检查图片是否存在
	var image models.Image
	if err := db.First(&image, imageID).Error; err != nil {
		return fmt.Errorf("图片不存在")
	}

	// 检查标签是否存在
	var tag models.Tag
	if err := db.First(&tag, tagID).Error; err != nil {
		return fmt.Errorf("标签不存在")
	}

	// 检查关联是否已存在
	var count int64
	if err := db.Model(&models.ImageTag{}).Where("image_id = ? AND tag_id = ?", imageID, tagID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("标签已添加到该图片")
	}

	// 创建关联
	imageTag := models.ImageTag{
		ImageID: imageID,
		TagID:   tagID,
	}
	return db.Create(&imageTag).Error
}

// GetImagesByTag 根据标签获取图片
func GetImagesByTag(db *gorm.DB, userID, tagID uint, page, pageSize int) ([]models.Image, int64, error) {
	var images []models.Image
	var total int64

	// 构建查询
	query := db.Table("images").
		Joins("JOIN image_tags ON images.image_id = image_tags.image_id").
		Where("images.user_id = ? AND image_tags.tag_id = ?", userID, tagID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Preload("Tags").
		Offset(offset).
		Limit(pageSize).
		Find(&images).Error

	if err != nil {
		return nil, 0, err
	}

	return images, total, nil
}
