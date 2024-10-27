package services

import (
	"fmt"
	"img_hosting/models"
)

func Getimage(user_id uint, img_name string) ([]models.Image, error) {
	var images []models.Image
	db := models.GetDB()

	// 使用 Preload 加载标签，并添加查询条件
	if err := db.Preload("Tags").
		Where("user_id = ? AND image_name = ?", user_id, img_name).
		Find(&images).Error; err != nil {
		return nil, fmt.Errorf("查询图像失败: %w", err)
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("未找到匹配的图像")
	}

	return images, nil
}

func GetAllimage(user_id uint) ([]models.Image, error) {
	var images []models.Image
	db := models.GetDB()

	// 使用 Preload 加载标签，并添加用户ID条件
	if err := db.Preload("Tags").
		Where("user_id = ?", user_id).
		Find(&images).Error; err != nil {
		return nil, fmt.Errorf("查询图像失败: %w", err)
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("未找到匹配的图像")
	}

	return images, nil
}
