package dao

import (
	"img_hosting/models"
)

func UserHasPermissions(userID uint, permissionNames []string) (map[string]bool, error) {
	db := models.GetDB()
	var permissions []struct {
		PermissionName string
	}

	err := db.Table("user_roles").
		Select("DISTINCT permissions.permission_name").
		Joins("JOIN role_permissions ON user_roles.role_id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.permission_id").
		Where("user_roles.user_id = ?", userID).
		Where("permissions.permission_name IN (?)", permissionNames).
		Scan(&permissions).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]bool)
	for _, perm := range permissionNames {
		result[perm] = false
	}
	for _, perm := range permissions {
		result[perm.PermissionName] = true
	}

	return result, nil
}
