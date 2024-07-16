package models

type RolePermission struct {
	RoleID       uint        `gorm:"primaryKey"`
	PermissionID uint        `gorm:"primaryKey"`
	Role         Roles       `gorm:"foreignKey:RoleID"`
	Permission   Permissions `gorm:"foreignKey:PermissionID"`
}
