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

func GetAllimage(user_id uint, enablePaging bool, page, pageSize int, order string) (*models.ImageResult, error) {
	db := models.GetDB()
	var result models.ImageResult

	// 基础查询
	query := db.Model(&models.Image{}).
		Where("user_id = ?", user_id).
		Preload("Tags")

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("查询总数失败: %w", err)
	}
	result.Total = int(total)

	// 添加排序
	if order == "asc" {
		query = query.Order("upload_time asc")
	} else {
		query = query.Order("upload_time desc")
	}

	// 如果启用分页，添加分页条件
	if enablePaging {
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	// 执行查询
	var images []models.Image
	if err := query.Find(&images).Error; err != nil {
		return nil, fmt.Errorf("查询图像失败: %w", err)
	}

	// 如果不需要分页，保持原有的返回格式
	if !enablePaging {
		return &models.ImageResult{
			Images: images,
			Total:  len(images),
		}, nil
	}

	// 返回带分页的结果
	result.Images = images
	return &result, nil
}
