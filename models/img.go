package models

import (
	"time"
)

// Image 图片结构体
type Image struct {
	ImageID       uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null" json:"user_id"` // 用户名
	ImageURL      string    `json:"image_url"`               // 图片存储路径或URL
	ImageName     string    `json:"image_name"`              // 图片名称
	Imageextenion string    `json:"image_extenion"`          // 图片扩展名
	HashImage     string    `json:"hash_image"`              //图片哈希名
	ImageSize     int64     `json:"image_size"`              // 图片大小（字节）
	ImageType     string    `json:"image_type"`              // 图片格式
	UploadTime    time.Time `gorm:"autoCreateTime"`          // 上传时间
	Description   string    `json:"description"`             // 图片描述（可选）
	Tags          []Tag     `gorm:"many2many:image_tags;foreignKey:ImageID;joinForeignKey:ImageID;references:TagID;joinReferences:TagID"`
}

type ImageResult struct {
	Images []Image `json:"images"`
	Total  int     `json:"total"`
}
