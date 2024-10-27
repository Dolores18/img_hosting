package services

import (
	"img_hosting/models"
	"img_hosting/pkg/logger"
)

func CreateImageTag(imageID uint, tagID uint) error {
	log := logger.GetLogger()
	log.Info("创建图片标签映射")
	db := models.GetDB()

	// 检查图片是否存在
	var image models.Image
	if err := db.First(&image, imageID).Error; err != nil {
		log.Error("图片不存在", err)
		return err
	}

	// 检查标签是否存在
	var tag models.Tag
	if err := db.First(&tag, tagID).Error; err != nil {
		log.Error("标签不存在", err)
		return err
	}

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 使用 GORM 关联添加标签
	if err := tx.Model(&image).Association("Tags").Append(&tag); err != nil {
		tx.Rollback()
		log.Error("创建图片标签映射失败", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("提交事务失败", err)
		return err
	}

	return nil
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
