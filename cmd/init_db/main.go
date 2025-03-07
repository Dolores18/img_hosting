package main

import (
	"img_hosting/config"
	"img_hosting/models"
	"log"

	"gorm.io/gorm"
)

// InitPermissions 初始化权限和角色
func InitPermissions(db *gorm.DB) error {
	// 从配置获取角色和权限
	roles := config.AppConfigInstance.Permissions.Roles

	for roleName, permissions := range roles {
		// 创建角色
		role := &models.Roles{RoleName: roleName}
		db.FirstOrCreate(role, models.Roles{RoleName: roleName})

		// 创建权限并关联
		for _, permName := range permissions {
			perm := &models.Permissions{Name: permName}
			db.FirstOrCreate(perm, models.Permissions{Name: permName})

			// 关联角色和权限
			db.FirstOrCreate(&models.RolePermission{
				RoleID:       role.RoleID,
				PermissionID: perm.PermissionID,
			})
		}
	}
	return nil
}

func main() {
	log.Println("开始初始化数据库...")

	// 加载配置
	config.LoadConfig()

	// 获取数据库连接
	db := models.GetDB()

	// 执行迁移
	err := db.AutoMigrate(
		&models.UserInfo{},
		&models.Image{},
		&models.Roles{},
		&models.Permissions{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.Tag{},
		&models.ImageTag{},
		&models.Token{},
		&models.PrivateFile{},
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化权限和角色
	if err := InitPermissions(db); err != nil {
		log.Fatalf("权限初始化失败: %v", err)
	}

	log.Println("数据库初始化完成")
}
