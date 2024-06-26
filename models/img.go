package models

import (
	"time"
)

// Image 图片结构体
type Image struct {
	ImageID     uint      `gorm:"primaryKey" json:"id"`
	UserName    string    `gorm:"not null" json:"username"` // 用户名
	ImageURL    string    `json:"image_url"`                // 图片存储路径或URL
	ImageName   string    `json:"image_name"`               // 图片名称
	ImageSize   int64     `json:"image_size"`               // 图片大小（字节）
	ImageType   string    `json:"image_type"`               // 图片格式
	UploadTime  time.Time `gorm:"autoCreateTime"`           // 上传时间
	Description string    `json:"description"`              // 图片描述（可选）
}
