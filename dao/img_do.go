package dao

import (
	"img_hosting/models"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// 创建图片
func CreateImg(db2 *gorm.DB, userid uint, image_url string, image_name string, image_extension string, hash_image string, image_size int64, image_type string) (uint, error) {
	// 验证图片名称长度
	if len(image_name) > 20 {
		return 0, fmt.Errorf("图片名称长度不能超过20个字符")
	}
	
	img := models.Image{
		UserID:        userid,
		ImageURL:      image_url,
		ImageName:     image_name,
		Imageextenion: image_extension,
		HashImage:     hash_image,
		ImageSize:     image_size,
		ImageType:     image_type,
		UploadTime:    time.Now(),
	}
	
	err := db2.Create(&img).Error
	if err != nil {
		log.Printf("创建图片记录失败: %v", err)
		return 0, err
	}
	
	log.Printf("图片创建成功，ID: %d", img.ImageID)
	return img.ImageID, nil
}
