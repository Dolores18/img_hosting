package dao

import (
	"errors"
	"img_hosting/models"

	"gorm.io/gorm"
)

// CreatePrivateFile 创建私人文件记录
func CreatePrivateFile(db *gorm.DB, file *models.PrivateFile) error {
	return db.Create(file).Error
}

// GetPrivateFileByID 通过ID获取私人文件
func GetPrivateFileByID(db *gorm.DB, fileID uint, userID uint) (*models.PrivateFile, error) {
	var file models.PrivateFile
	err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文件不存在或无权访问")
		}
		return nil, err
	}
	return &file, nil
}

// ListUserPrivateFiles 获取用户的所有私人文件
func ListUserPrivateFiles(db *gorm.DB, userID uint, page, pageSize int) ([]models.PrivateFile, int64, error) {
	var files []models.PrivateFile
	var total int64

	// 计算总数
	if err := db.Model(&models.PrivateFile{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&files).Error

	return files, total, err
}

// UpdatePrivateFile 更新私人文件信息
func UpdatePrivateFile(db *gorm.DB, file *models.PrivateFile) error {
	return db.Model(file).Updates(map[string]interface{}{
		"file_name":    file.FileName,
		"is_encrypted": file.IsEncrypted,
		"password":     file.Password,
		"status":       file.Status,
	}).Error
}

// DeletePrivateFile 删除私人文件（软删除）
func DeletePrivateFile(db *gorm.DB, fileID uint, userID uint) error {
	result := db.Where("id = ? AND user_id = ?", fileID, userID).Delete(&models.PrivateFile{})
	if result.RowsAffected == 0 {
		return errors.New("文件不存在或无权删除")
	}
	return result.Error
}

// IncrementViewCount 增加文件查看次数
func IncrementViewCount(db *gorm.DB, fileID uint) error {
	return db.Model(&models.PrivateFile{}).
		Where("id = ?", fileID).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).
		Error
}

// GetFilesByHash 通过文件哈希查找文件（用于查重）
func GetFilesByHash(db *gorm.DB, fileHash string, userID uint) (*models.PrivateFile, error) {
	var file models.PrivateFile
	err := db.Where("file_hash = ? AND user_id = ?", fileHash, userID).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 返回 nil 表示没有重复文件
		}
		return nil, err
	}
	return &file, nil
}

// SearchPrivateFiles 搜索私人文件
func SearchPrivateFiles(db *gorm.DB, userID uint, keyword string, page, pageSize int) ([]models.PrivateFile, int64, error) {
	var files []models.PrivateFile
	var total int64

	query := db.Model(&models.PrivateFile{}).Where("user_id = ?", userID)

	if keyword != "" {
		query = query.Where("file_name LIKE ?", "%"+keyword+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&files).Error

	return files, total, err
}
