package services

import (
	"fmt"
	"img_hosting/config"
	"img_hosting/models"
	"img_hosting/pkg/cache"

	"gorm.io/gorm"
)

// PermissionService 权限服务
type PermissionService struct {
	cache *cache.Service
}

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

// UpdateRolePermissions 更新角色权限并同步
func (ps *PermissionService) UpdateRolePermissions(roleName string, permissions []string) error {
	db := models.GetDB()

	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 获取角色
	var role models.Roles
	if err := tx.Where("role_name = ?", roleName).First(&role).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("角色不存在: %s", roleName)
	}

	// 2. 删除现有的角色权限关联
	if err := tx.Where("role_id = ?", role.RoleID).Delete(&models.RolePermission{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除现有权限失败: %w", err)
	}

	// 3. 添加新的权限关联
	for _, permName := range permissions {
		var perm models.Permissions
		if err := tx.Where("permission_name = ?", permName).First(&perm).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("权限不存在: %s", permName)
		}

		rolePermission := models.RolePermission{
			RoleID:       role.RoleID,
			PermissionID: perm.PermissionID,
		}

		if err := tx.Create(&rolePermission).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建权限关联失败: %w", err)
		}
	}

	// 4. 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	// 5. 清除权限缓存
	cache.ClearUserPermissionCache(role.RoleID)

	return nil
}

// CreatePermission 创建新权限并同步
func (ps *PermissionService) CreatePermission(name, description string) error {
	db := models.GetDB()
	tx := db.Begin()

	// 1. 创建权限
	perm := models.Permissions{
		Name:        name,
		Description: description,
	}
	if err := tx.Create(&perm).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. 同步配置到路由
	if err := ps.SyncConfigPermissions(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
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

// SyncConfigPermissions 同步权限配置到路由
func (ps *PermissionService) SyncConfigPermissions() error {
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

// DeletePermission 删除权限并同步
func (ps *PermissionService) DeletePermission(name string) error {
	db := models.GetDB()
	tx := db.Begin()

	// 1. 删除权限
	if err := tx.Where("permission_name = ?", name).Delete(&models.Permissions{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. 同步配置到路由
	if err := ps.SyncConfigPermissions(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// updateRolePermissionsInTx 在事务中更新角色权限
func (ps *PermissionService) updateRolePermissionsInTx(tx *gorm.DB, roleName string, permissions []string) error {
	// 获取角色信息
	var role models.Roles
	if err := tx.Where("role_name = ?", roleName).First(&role).Error; err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}

	// 删除现有权限
	if err := tx.Where("role_id = ?", role.RoleID).Delete(&models.RolePermission{}).Error; err != nil {
		return fmt.Errorf("删除现有权限失败: %w", err)
	}

	// 添加新权限
	for _, permName := range permissions {
		// 获取权限ID
		var perm models.Permissions
		if err := tx.Where("permission_name = ?", permName).First(&perm).Error; err != nil {
			return fmt.Errorf("权限不存在 %s: %w", permName, err)
		}

		rp := models.RolePermission{
			RoleID:       role.RoleID,
			PermissionID: perm.PermissionID,
		}
		if err := tx.Create(&rp).Error; err != nil {
			return fmt.Errorf("添加新权限失败: %w", err)
		}
	}

	return nil
}
