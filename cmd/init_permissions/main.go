package main

import (
	"img_hosting/models"
	"log"
)

func main() {
	log.Println("开始初始化权限...")
	db := models.GetDB()

	// 定义基础权限
	permissions := []models.Permissions{
		{Name: "manage_private_files", Description: "管理私人文件的权限"},
		{Name: "upload_file", Description: "上传文件权限"},
		{Name: "download_file", Description: "下载文件权限"},
		{Name: "delete_file", Description: "删除文件权限"},
		{Name: "view_file", Description: "查看文件权限"},
	}

	// 定义角色
	roles := []models.Roles{
		{RoleName: "admin", Description: "管理员"},
		{RoleName: "user", Description: "普通用户"},
	}

	// 创建权限
	for _, perm := range permissions {
		if err := db.FirstOrCreate(&perm, models.Permissions{Name: perm.Name}).Error; err != nil {
			log.Printf("创建权限失败 %s: %v\n", perm.Name, err)
		}
	}

	// 创建角色
	for _, role := range roles {
		if err := db.FirstOrCreate(&role, models.Roles{RoleName: role.RoleName}).Error; err != nil {
			log.Printf("创建角色失败 %s: %v\n", role.RoleName, err)
		}
	}

	// 为角色分配权限
	var adminRole models.Roles
	var userRole models.Roles
	db.Where("name = ?", "admin").First(&adminRole)
	db.Where("name = ?", "user").First(&userRole)

	// 获取所有权限
	var allPermissions []models.Permissions
	db.Find(&allPermissions)

	// 为管理员分配所有权限
	for _, perm := range allPermissions {
		rolePermission := models.RolePermission{
			RoleID:       adminRole.RoleID,
			PermissionID: perm.PermissionID,
		}
		db.FirstOrCreate(&rolePermission, rolePermission)
	}

	// 为普通用户分配基本权限
	userPermissions := []string{
		"manage_private_files",
		"upload_file",
		"download_file",
		"view_file",
	}

	for _, permName := range userPermissions {
		var perm models.Permissions
		if err := db.Where("name = ?", permName).First(&perm).Error; err != nil {
			continue
		}
		rolePermission := models.RolePermission{
			RoleID:       userRole.RoleID,
			PermissionID: perm.PermissionID,
		}
		db.FirstOrCreate(&rolePermission, rolePermission)
	}

	log.Println("权限初始化完成")
}
