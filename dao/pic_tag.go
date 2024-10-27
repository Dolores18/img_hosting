package dao

import (
	"img_hosting/models"

	"gorm.io/gorm"
)

func ImageHasTags(db *gorm.DB, userID uint, imageID uint, tagNames []string) (map[string]bool, error) {
	// 首先验证图片是否存在且属于该用户
	var image models.Image
	if err := db.Where("image_id = ? AND user_id = ?", imageID, userID).First(&image).Error; err != nil {
		return nil, err
	}

	var tags []struct {
		TagName string
	}

	err := db.Table("image_tags").
		Select("DISTINCT tags.tag_name").
		Joins("JOIN tags ON image_tags.tag_id = tags.tag_id").
		Joins("JOIN images ON image_tags.image_id = images.image_id").
		Where("images.user_id = ?", userID).
		Where("image_tags.image_id = ?", imageID).
		Where("tags.tag_name IN (?)", tagNames).
		Scan(&tags).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]bool)
	for _, tagName := range tagNames {
		result[tagName] = false
	}
	for _, tag := range tags {
		result[tag.TagName] = true
	}

	return result, nil
}

func GetImagesByTags(db *gorm.DB, userID uint, tagNames []string, enablePaging bool, page, pageSize int) (*models.ImageResult, error) {
	var result models.ImageResult

	// 使用 Model 而不是 Table
	query := db.Model(&models.Image{}).
		Distinct().
		Joins("JOIN image_tags ON images.image_id = image_tags.image_id").
		Joins("JOIN tags ON image_tags.tag_id = tags.tag_id").
		Where("images.user_id = ?", userID).
		Preload("Tags") // 现在可以使用 Preload 了

	if len(tagNames) == 1 {
		query = query.Where("tags.tag_name = ?", tagNames[0])
	} else {
		for _, tagName := range tagNames {
			subQuery := db.Table("image_tags").
				Select("image_id").
				Joins("JOIN tags ON image_tags.tag_id = tags.tag_id").
				Where("tags.tag_name = ?", tagName)
			query = query.Where("images.image_id IN (?)", subQuery)
		}
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	result.Total = int(total)

	// 排序
	query = query.Order("images.upload_time DESC")

	// 分页
	if enablePaging {
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	// 执行查询
	var images []models.Image
	if err := query.Find(&images).Error; err != nil {
		return nil, err
	}
	result.Images = images

	return &result, nil
}

// 创建标签
func CreateTag(db *gorm.DB, userID uint, tagName string) (uint, error) {
	tag := models.Tag{
		UserID:  userID, // 添加用户ID
		TagName: tagName,
	}
	err := db.Create(&tag).Error
	if err != nil {
		return 0, err
	}
	return tag.TagID, nil
}

// 查询用户的所有标签
func GetAllTag(db *gorm.DB, userID uint) ([]models.Tag, error) {
	var tags []models.Tag
	if err := db.Where("user_id = ?", userID).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// CreateImageTag 创建图片标签的映射
func CreateImageTag(db *gorm.DB, imageID uint, tagID uint) error {
	imageTag := models.ImageTag{ImageID: imageID, TagID: tagID}
	err := db.Create(&imageTag).Error
	if err != nil {
		return err
	}
	return nil
}

// 根据标签名查询标签ID
func GetTagIDByName(db *gorm.DB, userID uint, tagName string) (uint, error) {
	var tag models.Tag
	if err := db.Where("tag_name = ? AND user_id = ?", tagName, userID).First(&tag).Error; err != nil {
		return 0, err
	}
	return tag.TagID, nil
}
