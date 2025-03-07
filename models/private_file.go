package models

import (
	"time"

	"gorm.io/gorm"
)

// PrivateFile 私人文件模型
type PrivateFile struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`          // 所属用户
	FileName    string         `gorm:"size:255;not null" json:"file_name"`     // 文件名
	FileHash    string         `gorm:"size:64;not null" json:"file_hash"`      // 文件哈希值
	FileSize    int64          `gorm:"not null" json:"file_size"`              // 文件大小(字节)
	FileType    string         `gorm:"size:50" json:"file_type"`               // 文件类型(MIME类型)
	StoragePath string         `gorm:"size:512;not null" json:"storage_path"`  // 存储路径
	IsEncrypted bool           `gorm:"default:false" json:"is_encrypted"`      // 是否加密
	Password    string         `gorm:"size:255" json:"-"`                      // 加密密码(如果有)
	ViewCount   int64          `gorm:"default:0" json:"view_count"`            // 查看次数
	Status      string         `gorm:"size:20;default:'active'" json:"status"` // 文件状态(active/deleted)
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // 软删除

	User UserInfo `gorm:"foreignKey:UserID" json:"-"` // 关联用户
}

// FileStatus 定义文件状态常量
const (
	FileStatusActive  = "active"  // 正常
	FileStatusDeleted = "deleted" // 已删除
)

// BeforeCreate 创建前的钩子
func (pf *PrivateFile) BeforeCreate(tx *gorm.DB) error {
	if pf.Status == "" {
		pf.Status = FileStatusActive
	}
	return nil
}
