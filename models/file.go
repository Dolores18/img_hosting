package models

import (
	"time"
)

// Image 图片结构体
type File struct {
	FileID       uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null" json:"user_id"` // 用户名
	FileURL      string    `json:"file_url"`               // 图片存储路径或URL
	FileName     string    `json:"file_name"`              // 图片名称
	Fileextenion string    `json:"file_extenion"`          // 图片扩展名
	HashFile     string    `json:"hash_file"`              //图片哈希名
	FileSize     int64     `json:"file_size"`              // 图片大小（字节）
	FileType     string    `json:"file_type"`              // 图片格式
	UploadTime    time.Time `gorm:"autoCreateTime"`          // 上传时间
	Description   string    `json:"description"`             // 图片描述（可选）
	Tags          []Tag     `gorm:"many2many:file_tags;foreignKey:FileID;joinForeignKey:FileID;references:TagID;joinReferences:TagID"`
}