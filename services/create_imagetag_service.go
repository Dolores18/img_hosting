package services

import (
	"errors"
	"img_hosting/models"
	"img_hosting/pkg/logger"

	"gorm.io/gorm"
)

func CreateImageTags(imageID uint, userID uint, tagnames []string) error {
	log := logger.GetLogger()
	db := models.GetDB()

	// 检查图片是否存在
	var image models.Image
	if err := db.First(&image, imageID).Error; err != nil {
		log.Error("图片不存在", err)
		return err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, tagname := range tagnames {
		log.Info("处理标签: ", tagname)

		// 在事务中创建标签
		tagID, err := CreateTagWithTx(tx, userID, tagname)
		if err != nil {
			tx.Rollback()
			log.Error("创建标签失败", err)
			return err
		}

		// 检查图片标签关联
		var existingTag models.Tag
		err = tx.Table("tags").
			Joins("JOIN image_tags ON image_tags.tag_id = tags.tag_id").
			Where("tags.tag_id = ? AND image_tags.image_id = ?", tagID, imageID).
			First(&existingTag).Error

		// 忽略"记录不存在"错误
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			log.Error("检查标签关联失败", err)
			return err
		}

		// 如果标签已存在，继续下一个
		if existingTag.TagID != 0 {
			continue
		}

		// 创建新的关联
		if err := tx.Model(&image).Association("Tags").Append(&models.Tag{TagID: tagID}); err != nil {
			tx.Rollback()
			log.Error("创建图片标签映射失败", err)
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("提交事务失败", err)
		return err
	}

	return nil
}

// CreateTagWithTx 在事务中创建标签
func CreateTagWithTx(tx *gorm.DB, userID uint, tagname string) (uint, error) {
	var tag models.Tag

	// 先查找是否已存在相同的标签
	if err := tx.Where("tag_name = ?", tagname).First(&tag).Error; err == nil {
		// 标签已存在，直接返回ID
		return tag.TagID, nil
	}

	// 创建新标签
	tag = models.Tag{
		UserID:  userID,
		TagName: tagname, // 改用 TagName 而不是 Name
	}

	if err := tx.Create(&tag).Error; err != nil {
		return 0, err
	}

	return tag.TagID, nil // 改用 TagID 而不是 ID
}

// GetImage 获取单个图片及其标签
func GetImage(imageID uint) (*models.Image, error) {
	var image models.Image
	db := models.GetDB()
	if err := db.Preload("Tags").First(&image, imageID).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

// GetImages 获取所有图片及其标签
func GetImages() ([]models.Image, error) {
	var images []models.Image
	db := models.GetDB()
	if err := db.Preload("Tags").Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}
