package services

import (
	"fmt"
	"img_hosting/config"
	"img_hosting/models"

	"gorm.io/gorm"
)

type PermissionService struct{}

// GetAllPermissions 获取所有权限
func (s *PermissionService) GetAllPermissions() ([]models.Permissions, error) {
	db := models.GetDB()
	var permissions []models.Permissions
	if err := db.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetAllRoles 获取所有角色
func (s *PermissionService) GetAllRoles() ([]models.Roles, error) {
	db := models.GetDB()
	var roles []models.Roles
	if err := db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRolePermissions 获取角色的权限
func (s *PermissionService) GetRolePermissions(roleName string) ([]models.Permissions, error) {
	db := models.GetDB()
	var role models.Roles
	if err := db.Where("role_name = ?", roleName).First(&role).Error; err != nil {
		return nil, fmt.Errorf("角色不存在: %s", roleName)
	}

	var permissions []models.Permissions
	if err := db.Table("permissions").
		Joins("JOIN role_permissions ON permissions.permission_id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", role.RoleID).
		Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

// UpdateRolePermissions 更新角色的权限
func (s *PermissionService) UpdateRolePermissions(roleName string, permissionNames []string) error {
	db := models.GetDB()

	// 开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 获取角色ID
		var role models.Roles
		if err := tx.Where("role_name = ?", roleName).First(&role).Error; err != nil {
			return fmt.Errorf("角色不存在: %s", roleName)
		}

		// 2. 删除角色现有的所有权限
		if err := tx.Where("role_id = ?", role.RoleID).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}

		// 3. 添加新的权限
		for _, permName := range permissionNames {
			var perm models.Permissions
			if err := tx.Where("permission_name = ?", permName).First(&perm).Error; err != nil {
				return fmt.Errorf("权限不存在: %s", permName)
			}

			// 创建角色-权限关联
			rolePermission := models.RolePermission{
				RoleID:       role.RoleID,
				PermissionID: perm.PermissionID,
			}

			if err := tx.Create(&rolePermission).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// CreatePermission 创建新权限
func (s *PermissionService) CreatePermission(name, description string) error {
	db := models.GetDB()

	// 检查权限是否已存在
	var count int64
	db.Model(&models.Permissions{}).Where("permission_name = ?", name).Count(&count)
	if count > 0 {
		return fmt.Errorf("权限已存在: %s", name)
	}

	// 创建新权限
	permission := models.Permissions{
		Name:        name,
		Description: description,
	}

	return db.Create(&permission).Error
}

// CreateRole 创建新角色
func (s *PermissionService) CreateRole(name, description string) error {
	db := models.GetDB()

	// 检查角色是否已存在
	var count int64
	db.Model(&models.Roles{}).Where("role_name = ?", name).Count(&count)
	if count > 0 {
		return fmt.Errorf("角色已存在: %s", name)
	}

	// 创建新角色
	role := models.Roles{
		RoleName:    name,
		Description: description,
		IsActive:    true,
	}

	return db.Create(&role).Error
}

// SyncConfigPermissions 同步配置文件中的权限到数据库
func (s *PermissionService) SyncConfigPermissions() error {
	db := models.GetDB()

	// 从配置获取角色-权限映射
	rolePermissions := config.AppConfigInstance.Permissions.Roles

	// 开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 遍历所有角色
		for roleName, permissions := range rolePermissions {
			// 1. 确保角色存在
			var role models.Roles
			err := tx.Where("role_name = ?", roleName).First(&role).Error
			if err != nil {
				// 角色不存在，创建新角色
				role = models.Roles{
					RoleName:    roleName,
					Description: fmt.Sprintf("%s角色", roleName),
					IsActive:    true,
				}
				if err := tx.Create(&role).Error; err != nil {
					return err
				}
			}

			// 2. 确保所有权限存在
			for _, permName := range permissions {
				var perm models.Permissions
				err := tx.Where("permission_name = ?", permName).First(&perm).Error
				if err != nil {
					// 权限不存在，创建新权限
					perm = models.Permissions{
						Name:        permName,
						Description: fmt.Sprintf("%s权限", permName),
					}
					if err := tx.Create(&perm).Error; err != nil {
						return err
					}
				}

				// 3. 创建角色-权限关联（如果不存在）
				var rolePermCount int64
				tx.Model(&models.RolePermission{}).
					Where("role_id = ? AND permission_id = ?", role.RoleID, perm.PermissionID).
					Count(&rolePermCount)

				if rolePermCount == 0 {
					rolePermission := models.RolePermission{
						RoleID:       role.RoleID,
						PermissionID: perm.PermissionID,
					}
					if err := tx.Create(&rolePermission).Error; err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}

// GetUserPermissions 获取用户的所有权限
func (s *PermissionService) GetUserPermissions(userID uint) ([]models.Permissions, error) {
	db := models.GetDB()

	// 查询用户通过角色获得的所有权限
	var permissions []models.Permissions
	err := db.Table("permissions").
		Select("DISTINCT permissions.*").
		Joins("JOIN role_permissions ON permissions.permission_id = role_permissions.permission_id").
		Joins("JOIN roles ON role_permissions.role_id = roles.role_id").
		Joins("JOIN user_roles ON roles.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.is_active = ?", userID, true).
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	// 移除直接分配权限的查询，因为没有user_permissions表
	// 只通过角色获取权限

	return permissions, nil
}

// GetUserPermissionMap 获取用户权限映射（权限名称 -> 是否拥有）
func (s *PermissionService) GetUserPermissionMap(userID uint) (map[string]bool, error) {
	permissions, err := s.GetUserPermissions(userID)
	if err != nil {
		return nil, err
	}

	// 转换为映射
	permMap := make(map[string]bool)
	for _, p := range permissions {
		permMap[p.Name] = true
	}

	return permMap, nil
}
