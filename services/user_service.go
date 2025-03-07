package services

import (
	"errors"
	"fmt"
	"img_hosting/dao"
	"img_hosting/models"

	"gorm.io/gorm"
)

type UserService struct{}

// GetUserProfile 获取用户信息
func (s *UserService) GetUserProfile(userID uint) (*models.UserInfo, error) {
	db := models.GetDB()
	return dao.GetUserByID(db, userID)
}

// UpdateUserProfile 更新用户信息
func (s *UserService) UpdateUserProfile(userID uint, updates map[string]interface{}) error {
	db := models.GetDB()
	user, err := dao.GetUserByID(db, userID)
	if err != nil {
		return err
	}

	// 只更新允许的字段
	allowedFields := map[string]bool{
		"name":   true,
		"email":  true,
		"phone":  true,
		"age":    true,
		"status": true,
	}

	updateData := make(map[string]interface{})
	for k, v := range updates {
		if allowedFields[k] {
			updateData[k] = v
		}
	}

	user.Name = updateData["name"].(string)
	user.Email = updateData["email"].(string)
	if phone, ok := updateData["phone"].(string); ok {
		user.Phone = phone
	}
	if age, ok := updateData["age"].(int); ok {
		user.Age = age
	}
	if status, ok := updateData["status"].(string); ok {
		user.Status = status
	}

	return dao.UpdateUser(db, user)
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(page, pageSize int, search string) ([]models.UserInfo, int64, error) {
	db := models.GetDB()
	return dao.ListUsers(db, page, pageSize, search)
}

// DeleteUser 删除用户及其所有关联数据
func (s *UserService) DeleteUser(userID uint) error {
	db := models.GetDB()

	// 开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除用户角色关联
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
			return err
		}

		// 2. 删除用户的Token
		if err := tx.Where("user_id = ?", userID).Delete(&models.Token{}).Error; err != nil {
			return err
		}

		// 3. 删除用户的私人文件
		if err := tx.Where("user_id = ?", userID).Delete(&models.PrivateFile{}).Error; err != nil {
			return err
		}

		// 4. 删除用户的图片标签关联
		var userImages []models.Image
		if err := tx.Where("user_id = ?", userID).Find(&userImages).Error; err != nil {
			return err
		}

		var imageIDs []uint
		for _, img := range userImages {
			imageIDs = append(imageIDs, img.ImageID)
		}

		if len(imageIDs) > 0 {
			if err := tx.Where("image_id IN ?", imageIDs).Delete(&models.ImageTag{}).Error; err != nil {
				return err
			}
		}

		// 5. 删除用户的图片
		if err := tx.Where("user_id = ?", userID).Delete(&models.Image{}).Error; err != nil {
			return err
		}

		// 6. 删除用户创建的标签
		if err := tx.Where("user_id = ?", userID).Delete(&models.Tag{}).Error; err != nil {
			return err
		}

		// 7. 最后删除用户
		return tx.Delete(&models.UserInfo{}, userID).Error
	})
}

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(userID uint, status string) error {
	if status != models.UserStatusActive &&
		status != models.UserStatusInactive &&
		status != models.UserStatusBanned {
		return errors.New("无效的用户状态")
	}

	db := models.GetDB()
	return dao.UpdateUserStatus(db, userID, status)
}

// ManageUserRoles 管理用户角色
func (s *UserService) ManageUserRoles(userID uint, roleID uint, isAdd bool) error {
	db := models.GetDB()
	if isAdd {
		return dao.AssignRoleToUser(db, userID, roleID)
	}
	return dao.RemoveRoleFromUser(db, userID, roleID)
}

// GetUserRoles 获取用户角色
func (s *UserService) GetUserRoles(userID uint) ([]models.Roles, error) {
	db := models.GetDB()
	return dao.GetUserRoles(db, userID)
}

// SearchUsers 搜索用户
func (s *UserService) SearchUsers(keyword string, page, pageSize int) ([]models.UserInfo, int64, error) {
	db := models.GetDB()
	return dao.ListUsers(db, page, pageSize, keyword)
}

// AssignRoles 为用户分配角色
func (s *UserService) AssignRoles(userID uint, roleNames []string) error {
	fmt.Printf("开始分配角色: userID=%d, roles=%v\n", userID, roleNames)
	db := models.GetDB()

	// 开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 先删除用户现有的所有角色
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
			fmt.Printf("删除现有角色失败: %v\n", err)
			return err
		}

		// 2. 获取角色ID
		for _, roleName := range roleNames {
			var role models.Roles
			// 打印SQL查询
			fmt.Printf("查询角色SQL: SELECT * FROM roles WHERE role_name = '%s'\n", roleName)

			// 使用 role_name 而不是 name
			if err := tx.Where("role_name = ?", roleName).First(&role).Error; err != nil {
				fmt.Printf("获取角色失败: role_name=%s, error=%v\n", roleName, err)
				return fmt.Errorf("角色不存在: %s", roleName)
			}

			fmt.Printf("找到角色: %+v\n", role)

			// 3. 分配新角色
			if err := dao.AssignRoleToUser(tx, userID, role.RoleID); err != nil {
				fmt.Printf("分配角色失败: userID=%d, roleID=%d, error=%v\n",
					userID, role.RoleID, err)
				return err
			}
			fmt.Printf("成功分配角色: userID=%d, roleName=%s, roleID=%d\n",
				userID, roleName, role.RoleID)
		}

		return nil
	})
}
