package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

//检查上传图片的格式以及大小

// CheckImg 检查上传图片的格式以及大小
func CheckImg(filename string, filesize int64) (bool, string, string, int64) {
	imgExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	ext := strings.ToLower(filepath.Ext(filename))
	const maxSize = 10485760 // 10MB

	// 检查文件扩展名
	validExtension := false
	for _, validExt := range imgExtensions {
		if ext == validExt {
			validExtension = true
			break
		}
	}
	if !validExtension {
		return false, "图片后缀名不对", "", filesize
	}

	// 检查文件大小
	if filesize > maxSize {
		return false, "文件大小超出限制", "", filesize
	}

	// 文件格式和大小都合法
	return true, "图片合法", ext[1:], filesize
}

// SanitizeFileName 验证并清理文件名
func SanitizeFileName(filename string) (string, string, error) {
	// 获取文件扩展名
	extension := filepath.Ext(filename)
	// 获取文件名（不包括扩展名）
	name := strings.TrimSuffix(filename, extension)

	// 检查文件名长度（不包括扩展名）
	if len([]rune(name)) > 20 { // 使用 []rune 来正确计算中文字符长度
		return "", "", fmt.Errorf("文件名长度不能超过20个字符: %s", name)
	}

	// 允许中文、字母、数字、下划线和横线
	validName := regexp.MustCompile(`^[\p{Han}a-zA-Z0-9_-]+$`).MatchString(name)
	if !validName {
		return "", "", fmt.Errorf("文件名只能包含中文、字母、数字、下划线和横线: %s", name)
	}

	return name, extension, nil
}

// HashFileName 计算文件内容的MD5哈希，并返回哈希值的前16个字符
func HashFileName(file []byte) string {
	hash := md5.Sum(file)
	hashString := hex.EncodeToString(hash[:])
	return hashString[:16]
}
