package services

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"img_hosting/dao"
	"img_hosting/models"
	"img_hosting/pkg/encryption"
	"img_hosting/pkg/logger"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	MaxFileSize  = 100 * 1024 * 1024                                      // 100MB
	UploadPath   = "./uploads/private"                                    // 私人文件上传路径
	AllowedTypes = ".jpg,.jpeg,.png,.gif,.pdf,.doc,.docx,.xls,.xlsx,.txt" // 允许的文件类型
)

// UploadPrivateFile 上传私人文件
func UploadPrivateFile(file *multipart.FileHeader, userID uint, isEncrypted bool, password string) (*models.PrivateFile, error) {
	log := logger.GetLogger()

	// 检查文件大小
	if file.Size > MaxFileSize {
		return nil, errors.New("文件大小超过限制")
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !strings.Contains(AllowedTypes, ext) {
		return nil, errors.New("不支持的文件类型")
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// 计算文件哈希
	hash := md5.New()
	if _, err := io.Copy(hash, src); err != nil {
		return nil, err
	}
	fileHash := hex.EncodeToString(hash.Sum(nil))

	// 检查文件是否已存在
	db := models.GetDB()
	existingFile, err := dao.GetFilesByHash(db, fileHash, userID)
	if err != nil {
		return nil, err
	}
	if existingFile != nil {
		return nil, errors.New("文件已存在")
	}

	// 创建上传目录
	uploadDir := filepath.Join(UploadPath, fmt.Sprintf("user_%d", userID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	// 生成存储路径
	storagePath := filepath.Join(uploadDir, fileHash+ext)

	// 重新打开文件（因为之前的文件指针已经到达末尾）
	src.Seek(0, 0)

	// 创建目标文件
	dst, err := os.Create(storagePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	// 如果需要加密，则加密文件
	if isEncrypted {
		log.WithField("file", storagePath).Info("开始加密文件")
		if err := encryption.EncryptFileInPlace(storagePath, password); err != nil {
			// 如果加密失败，删除文件
			os.Remove(storagePath)
			return nil, fmt.Errorf("文件加密失败: %w", err)
		}
	}

	// 创建文件记录
	privateFile := &models.PrivateFile{
		UserID:      userID,
		FileName:    file.Filename,
		FileHash:    fileHash,
		FileSize:    file.Size,
		FileType:    file.Header.Get("Content-Type"),
		StoragePath: storagePath,
		IsEncrypted: isEncrypted,
		Password:    password,
		Status:      models.FileStatusActive,
	}

	if err := dao.CreatePrivateFile(db, privateFile); err != nil {
		// 如果数据库创建失败，删除已上传的文件
		os.Remove(storagePath)
		return nil, err
	}

	return privateFile, nil
}

// GetPrivateFile 获取私人文件信息
func GetPrivateFile(fileID, userID uint, password string) (*models.PrivateFile, error) {
	db := models.GetDB()
	file, err := dao.GetPrivateFileByID(db, fileID, userID)
	if err != nil {
		return nil, err
	}

	// 检查加密文件的密码
	if file.IsEncrypted && file.Password != password {
		return nil, errors.New("密码错误")
	}

	// 增加查看次数
	dao.IncrementViewCount(db, fileID)

	return file, nil
}

// ListPrivateFiles 获取用户的私人文件列表
func ListPrivateFiles(userID uint, page, pageSize int) ([]models.PrivateFile, int64, error) {
	return dao.ListUserPrivateFiles(models.GetDB(), userID, page, pageSize)
}

// DeletePrivateFile 删除私人文件
func DeletePrivateFile(fileID, userID uint) error {
	db := models.GetDB()

	// 获取文件信息
	file, err := dao.GetPrivateFileByID(db, fileID, userID)
	if err != nil {
		return err
	}

	// 删除物理文件
	if err := os.Remove(file.StoragePath); err != nil {
		return err
	}

	// 删除数据库记录
	return dao.DeletePrivateFile(db, fileID, userID)
}

// SearchPrivateFiles 搜索私人文件
func SearchPrivateFiles(userID uint, keyword string, page, pageSize int) ([]models.PrivateFile, int64, error) {
	return dao.SearchPrivateFiles(models.GetDB(), userID, keyword, page, pageSize)
}

// GetDecryptedFilePath 下载并解密文件
func GetDecryptedFilePath(fileID, userID uint, password string) (string, error) {
	log := logger.GetLogger()

	// 获取文件信息
	file, err := GetPrivateFile(fileID, userID, password)
	if err != nil {
		return "", err
	}

	// 如果文件不是加密的，直接返回原路径
	if !file.IsEncrypted {
		return file.StoragePath, nil
	}

	// 为解密文件创建临时目录
	tempDir := os.TempDir()
	decryptedPath := filepath.Join(tempDir, file.FileName)

	// 解密文件
	log.WithFields(logrus.Fields{
		"source": file.StoragePath,
		"target": decryptedPath,
	}).Info("开始解密文件")

	if err := encryption.DecryptFile(file.StoragePath, decryptedPath, password); err != nil {
		return "", fmt.Errorf("文件解密失败: %w", err)
	}

	return decryptedPath, nil
}

// UpdatePrivateFileInfo 更新私人文件信息
func UpdatePrivateFileInfo(fileID, userID uint, fileName string, isEncrypted bool, password string) (*models.PrivateFile, error) {
	fmt.Println("\n【更新文件】开始更新文件信息:", fileID)
	fmt.Println("【更新文件】参数:", "userID=", userID, "fileName=", fileName, "isEncrypted=", isEncrypted, "hasPassword=", password != "")

	db := models.GetDB()

	// 获取文件信息
	file, err := dao.GetPrivateFileByID(db, fileID, userID)
	if err != nil {
		fmt.Println("【更新文件错误】获取文件信息失败:", err)
		return nil, err
	}

	fmt.Println("【更新文件】当前文件信息:", "path=", file.StoragePath, "isEncrypted=", file.IsEncrypted)

	// 检查文件是否存在，如果不存在，尝试查找带.zip后缀的文件
	if _, err := os.Stat(file.StoragePath); os.IsNotExist(err) {
		zipPath := file.StoragePath + ".zip"
		if _, err := os.Stat(zipPath); os.IsNotExist(err) {
			fmt.Println("【更新文件错误】文件不存在:", file.StoragePath, "也不存在:", zipPath)
			return nil, fmt.Errorf("文件不存在: %s", file.StoragePath)
		} else {
			// 找到了带.zip后缀的文件，更新路径
			fmt.Println("【更新文件】找到带.zip后缀的文件:", zipPath)
			file.StoragePath = zipPath
			file.IsEncrypted = true
			// 更新数据库中的路径
			if err := dao.UpdatePrivateFile(db, file); err != nil {
				fmt.Println("【更新文件错误】更新数据库记录失败:", err)
				return nil, err
			}
		}
	}
	

	fmt.Println("【更新文件】文件存在:", file.StoragePath)

	// 处理加密状态变更
	if file.IsEncrypted != isEncrypted {
		if isEncrypted {
			// 从未加密变为加密
			if password == "" {
				fmt.Println("【更新文件错误】加密文件必须提供密码")
				return nil, errors.New("加密文件必须提供密码")
			}
			fmt.Println("【更新文件】将文件从未加密状态变为加密状态:", file.StoragePath)
			if err := encryption.EncryptFileInPlace(file.StoragePath, password); err != nil {
				fmt.Println("【更新文件错误】文件加密失败:", err)
				return nil, fmt.Errorf("文件加密失败: %w", err)
			}
			// 更新文件路径，添加.zip后缀
			file.StoragePath = file.StoragePath + ".zip"
			fmt.Println("【更新文件】文件加密成功，新路径:", file.StoragePath)
		} else {
			// 从加密变为未加密
			fmt.Println("【更新文件】将文件从加密状态变为未加密状态:", file.StoragePath)
			if err := encryption.DecryptFileInPlace(file.StoragePath, file.Password); err != nil {
				fmt.Println("【更新文件错误】文件解密失败:", err)
				return nil, fmt.Errorf("文件解密失败: %w", err)
			}
			// 更新文件路径，移除.zip后缀
			if strings.HasSuffix(file.StoragePath, ".zip") {
				file.StoragePath = strings.TrimSuffix(file.StoragePath, ".zip")
				fmt.Println("【更新文件】文件解密成功，新路径:", file.StoragePath)
			}
		}
	} else if isEncrypted && password != "" {
		// 如果文件已标记为加密，且提供了密码，则强制重新加密
		fmt.Println("【更新文件】强制重新加密文件:", file.StoragePath)

		// 如果密码变更，先尝试用旧密码解密
		if file.Password != password && file.Password != "" {
			fmt.Println("【更新文件】尝试用旧密码解密文件")
			tempPath := file.StoragePath + ".temp"
			if err := encryption.DecryptFile(file.StoragePath, tempPath, file.Password); err != nil {
				// 如果解密失败，可能文件实际上并未加密，继续处理
				fmt.Println("【更新文件】用旧密码解密失败，可能文件未实际加密:", err)
				os.Remove(tempPath)
			} else {
				fmt.Println("【更新文件】用旧密码解密成功")
				// 删除原加密文件
				os.Remove(file.StoragePath)
				// 重命名解密后的临时文件
				os.Rename(tempPath, file.StoragePath)
			}
		}

		// 使用新密码加密
		if err := encryption.EncryptFileInPlace(file.StoragePath, password); err != nil {
			fmt.Println("【更新文件错误】文件加密失败:", err)
			return nil, fmt.Errorf("文件加密失败: %w", err)
		}
		// 更新文件路径，确保有.zip后缀
		if !strings.HasSuffix(file.StoragePath, ".zip") {
			file.StoragePath = file.StoragePath + ".zip"
		}
		fmt.Println("【更新文件】文件强制加密成功，新路径:", file.StoragePath)
	} else {
		fmt.Println("【更新文件】文件加密状态未变更，无需处理文件内容")
	}

	// 更新文件信息
	file.FileName = fileName
	file.IsEncrypted = isEncrypted

	// 只有在提供了新密码时才更新密码
	if password != "" {
		file.Password = password
	}

	// 保存更新
	fmt.Println("【更新文件】保存文件信息到数据库")
	if err := dao.UpdatePrivateFile(db, file); err != nil {
		fmt.Println("【更新文件错误】更新数据库记录失败:", err)
		return nil, err
	}

	fmt.Println("【更新文件】文件信息更新成功:", fileID)
	return file, nil
}
