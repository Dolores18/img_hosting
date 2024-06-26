package services

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

//检查上传图片的格式以及大小

// CheckImg 检查上传图片的格式以及大小
func CheckImg(filename string, filesize int64) (bool, string, string, int64) {
	imgExtensions := []string{"jpg", "jpeg", "png", "gif", "webp"}

	// 检查文件名是否以任何一个扩展名结尾
	validFormat := false
	var imgExt string
	for _, extension := range imgExtensions {
		if strings.HasSuffix(strings.ToLower(filename), extension) {
			validFormat = true
			imgExt = extension

			break
		}
	}
	if !validFormat {
		return false, "图片后缀名不对", "", filesize
	}
	// 检查文件大小
	const maxSize = 10485760 // 这里的大小单位是字节，可以根据需要调整
	if filesize > maxSize {
		return false, "File size exceeds the limit", "", filesize
	}

	// 文件格式和大小都合法
	return true, "图片合法", imgExt, filesize
}

// 预防注入攻击
func SanitizeFileName(filename string) (string, error) {
	// 获取文件扩展名
	extension := filepath.Ext(filename)
	// 获取文件名（不包括扩展名）
	name := strings.TrimSuffix(filename, extension)
	// 只允许字母、数字、下划线和横线
	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name)
	if !validName {
		return "", fmt.Errorf("无效的文件名: %s", name)
	}
	return name + extension, nil
}
