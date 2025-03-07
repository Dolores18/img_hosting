package services

import (
	"fmt"
	"img_hosting/dao"
	"img_hosting/models"
)

// GetAllTag 获取用户的所有标签
func GetAllTag(userID uint) ([]models.Tag, error) {
	db := models.GetDB()
	return dao.GetAllTags(db, userID)
}

// CreateTag 创建新标签
func CreateTag(userID uint, tagName string) (uint, error) {
	db := models.GetDB()
	tag, err := dao.CreateTag(db, userID, tagName)
	if err != nil {
		return 0, err
	}
	return tag.TagID, nil
}

// AddTagToImage 为图片添加标签
func AddTagToImage(imageID, tagID, userID uint) error {
	db := models.GetDB()

	// 检查图片是否属于该用户
	image, err := dao.GetImageByID(db, imageID)
	if err != nil {
		return err
	}

	if image.UserID != userID {
		return fmt.Errorf("无权操作该图片")
	}

	// 添加标签
	return dao.AddImageTag(db, imageID, tagID)
}

// GetImagesByTag 获取带有特定标签的图片
func GetImagesByTag(userID, tagID uint, page, pageSize int) ([]models.Image, int64, error) {
	db := models.GetDB()
	return dao.GetImagesByTag(db, userID, tagID, page, pageSize)
}
