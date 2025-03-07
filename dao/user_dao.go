package dao

import (
	"errors"
	"fmt"
	"img_hosting/models"
	"time"

	"gorm.io/gorm"
)

// CreateUser 创建新用户
func CreateUser(db *gorm.DB, user *models.UserInfo) error {
	return db.Create(user).Error
}

// GetUserByID 通过ID获取用户
func GetUserByID(db *gorm.DB, userID uint) (*models.UserInfo, error) {
	var user models.UserInfo
	err := db.Preload("Roles").First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 通过邮箱获取用户
func GetUserByEmail(db *gorm.DB, email string) (*models.UserInfo, error) {
	var user models.UserInfo
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByName 通过用户名获取用户
func GetUserByName(db *gorm.DB, name string) (*models.UserInfo, error) {
	var user models.UserInfo
	err := db.Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(db *gorm.DB, user *models.UserInfo) error {
	return db.Model(user).Updates(map[string]interface{}{
		"name":          user.Name,
		"email":         user.Email,
		"phone":         user.Phone,
		"status":        user.Status,
		"last_login_at": user.LastLoginAt,
		"last_login_ip": user.LastLoginIP,
	}).Error
}

// UpdateUserStatus 更新用户状态
func UpdateUserStatus(db *gorm.DB, userID uint, status string) error {
	return db.Model(&models.UserInfo{}).
		Where("user_id = ?", userID).
		Update("status", status).Error
}

// DeleteUser 删除用户（软删除）
func DeleteUser(db *gorm.DB, userID uint) error {
	return db.Delete(&models.UserInfo{}, userID).Error
}

// ListUsers 获取用户列表（支持分页和搜索）
func ListUsers(db *gorm.DB, page, pageSize int, search string) ([]models.UserInfo, int64, error) {
	var users []models.UserInfo
	var total int64

	query := db.Model(&models.UserInfo{})

	if search != "" {
		fmt.Printf("搜索条件: %s\n", search)
		query = query.Where("name LIKE ? OR email LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		fmt.Printf("获取总数失败: %v\n", err)
		return nil, 0, err
	}
	fmt.Printf("总记录数: %d\n", total)

	// 简化 Preload，让 GORM 使用默认的关联
	err = query.Preload("Roles").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error

	if err != nil {
		fmt.Printf("查询数据失败: %v\n", err)
		return nil, 0, err
	}

	fmt.Printf("查询到 %d 条记录\n", len(users))
	return users, total, err
}

// AssignRoleToUser 为用户分配角色
func AssignRoleToUser(db *gorm.DB, userID uint, roleID uint) error {
	return db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)",
		userID, roleID).Error
}

// RemoveRoleFromUser 移除用户的角色
func RemoveRoleFromUser(db *gorm.DB, userID uint, roleID uint) error {
	return db.Exec("DELETE FROM user_roles WHERE user_id = ? AND role_id = ?",
		userID, roleID).Error
}

// GetUserRoles 获取用户的所有角色
func GetUserRoles(db *gorm.DB, userID uint) ([]models.Roles, error) {
	var user models.UserInfo
	err := db.Preload("Roles").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return user.Roles, nil
}

// UpdateLoginInfo 更新用户的登录信息
func UpdateLoginInfo(db *gorm.DB, userID uint, ip string) error {
	fmt.Printf("更新用户登录信息: userID=%d, ip=%s\n", userID, ip)

	// 使用 time.Now() 替代 NOW()
	result := db.Model(&models.UserInfo{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at": time.Now(), // 使用 time.Now() 而不是 NOW()
			"last_login_ip": ip,
		})

	if result.Error != nil {
		fmt.Printf("更新登录信息失败: %v\n", result.Error)
		return result.Error
	}

	fmt.Printf("更新登录信息成功: 影响行数=%d\n", result.RowsAffected)
	return nil
}
