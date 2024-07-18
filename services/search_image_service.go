package services

import (
	"fmt"
	"img_hosting/dao"
	"img_hosting/models"
)

func Getimage(user_id uint, img_name string) ([]map[string]interface{}, error) {
	conditions := map[string]interface{}{
		"user_id":    user_id,
		"image_name": img_name,
	}
	db := models.GetDB()
	imageDetails, err := dao.FindByFields(db, models.Image{}, conditions, false)
	if err != nil {
		return nil, fmt.Errorf("查询图像失败: %w", err)
	}
	if len(imageDetails) == 0 {
		return nil, fmt.Errorf("未找到匹配的图像")
	}
	return imageDetails, nil
}
func GetAllimage(user_id uint) ([]map[string]interface{}, error) {

	conditions := map[string]interface{}{
		"user_id": user_id,
	}
	db := models.GetDB()
	imageDetails, err := dao.FindByFields(db, models.Image{}, conditions, false)
	if err != nil {
		return nil, fmt.Errorf("查询图像失败: %w", err)
	}
	if len(imageDetails) == 0 {
		return nil, fmt.Errorf("未找到匹配的图像")
	}
	return imageDetails, nil
}
