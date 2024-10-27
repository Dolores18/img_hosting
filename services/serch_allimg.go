package services

import (
	"fmt"
	"img_hosting/dao"
	"img_hosting/models"
	"strings"
	"time"
)

type ImageFilter struct {
	TimeRange string   `json:"timeRange"` // e.g., "last_week", "last_month", "all"
	Tags      []string `json:"tags"`
	IsPublic  *bool    `json:"isPublic"`
	SortBy    string   `json:"sortBy"`
	Order     string   `json:"order"`
	Limit     int      `json:"limit"`
	Offset    int      `json:"offset"`
}

// ImageResult 用于返回图片查询结果
type ImageResult = models.ImageResult // 使用类型别名

func GetImgs(userID uint, filter ImageFilter) (*ImageResult, error) {
	db := models.GetDB()
	query := db.Where("user_id = ?", userID)

	// 应用时间范围过滤
	switch filter.TimeRange {
	case "last_week":
		query = query.Where("created_at > ?", time.Now().AddDate(0, 0, -7))
	case "last_month":
		query = query.Where("created_at > ?", time.Now().AddDate(0, -1, 0))
	case "all":
		// 不添加时间限制
	}

	// 应用标签过滤（通用方法）
	if len(filter.Tags) > 0 {
		// 假设 tags 是以逗号分隔的字符串存储在数据库中
		tagConditions := make([]string, len(filter.Tags))
		tagValues := make([]interface{}, len(filter.Tags))
		for i, tag := range filter.Tags {
			tagConditions[i] = "tags LIKE ?"
			tagValues[i] = fmt.Sprintf("%%,%s,%%", tag)
		}
		query = query.Where(strings.Join(tagConditions, " OR "), tagValues...)
	}

	// 应用公开/私密过滤
	if filter.IsPublic != nil {
		query = query.Where("is_public = ?", *filter.IsPublic)
	}

	// 应用排序
	if filter.SortBy != "" {
		order := "ASC"
		if filter.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", filter.SortBy, order))
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("计算总数失败: %w", err)
	}

	// 应用分页
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var images []models.Image
	if err := query.Find(&images).Error; err != nil {
		return nil, fmt.Errorf("查询图像失败: %w", err)
	}

	return &ImageResult{
		Images: images,
		Total:  int(total),
	}, nil
}

// GetImagesByTags 根据标签搜索图片
func GetImagesByTags(userID uint, tagNames []string) (*ImageResult, error) {
	db := models.GetDB()
	return dao.GetImagesByTags(db, userID, tagNames, false, 0, 0)
}
func GetAllTag(userID uint) ([]models.Tag, error) {
	db := models.GetDB()
	tags, err := dao.GetAllTag(db, userID)
	if err != nil {
		return nil, fmt.Errorf("查询标签失败: %w", err)
	}
	return tags, nil
}
func CreateTag(userID uint, tagName string) (uint, error) {
	db := models.GetDB()
	tagID, err := dao.CreateTag(db, userID, tagName)
	if err != nil {
		return 0, fmt.Errorf("创建标签失败: %w", err)
	}
	return tagID, nil
}
func GetTagIDByName(userID uint, tagName string) (uint, error) {
	db := models.GetDB()
	return dao.GetTagIDByName(db, userID, tagName)
}
