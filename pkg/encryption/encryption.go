package encryption

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/yeka/zip"
)

// EncryptFile 加密文件（使用ZIP加密）
func EncryptFile(inputPath, outputPath, password string) error {
	fmt.Println("【加密】开始加密文件:", inputPath, "->", outputPath)

	// 检查输入文件是否存在
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		fmt.Println("【加密错误】输入文件不存在:", inputPath)
		return fmt.Errorf("输入文件不存在: %s", inputPath)
	}

	// 创建ZIP文件
	zipFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("【加密错误】创建ZIP文件失败:", err)
		return err
	}
	defer zipFile.Close()

	// 创建带密码的ZIP writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 获取原始文件名
	originalFileName := filepath.Base(inputPath)

	// 创建带密码的writer
	writer, err := zipWriter.Encrypt(originalFileName, password, zip.StandardEncryption)
	if err != nil {
		fmt.Println("【加密错误】创建加密writer失败:", err)
		return err
	}

	// 读取原始文件
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		fmt.Println("【加密错误】读取原始文件失败:", err)
		return err
	}

	// 写入加密ZIP
	_, err = writer.Write(data)
	if err != nil {
		fmt.Println("【加密错误】写入加密数据失败:", err)
		return err
	}

	fmt.Println("【加密】文件加密成功:", outputPath)
	return nil
}

// DecryptFile 解密文件（从ZIP中提取）
func DecryptFile(inputPath, outputPath, password string) error {
	fmt.Println("【解密】开始解密文件:", inputPath, "->", outputPath)

	// 检查输入文件是否存在
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		fmt.Println("【解密错误】输入文件不存在:", inputPath)
		return fmt.Errorf("输入文件不存在: %s", inputPath)
	}

	// 打开ZIP文件
	zipReader, err := zip.OpenReader(inputPath)
	if err != nil {
		fmt.Println("【解密错误】打开ZIP文件失败:", err)
		return err
	}
	defer zipReader.Close()

	// 检查ZIP文件中是否有文件
	if len(zipReader.File) == 0 {
		fmt.Println("【解密错误】ZIP文件为空")
		return fmt.Errorf("ZIP文件为空")
	}

	// 获取第一个文件（我们只加密一个文件）
	zippedFile := zipReader.File[0]

	// 设置密码
	zippedFile.SetPassword(password)

	// 打开加密文件
	fileInZip, err := zippedFile.Open()
	if err != nil {
		fmt.Println("【解密错误】打开ZIP中的文件失败（可能密码错误）:", err)
		return fmt.Errorf("解密失败（可能密码错误）: %w", err)
	}
	defer fileInZip.Close()

	// 创建输出文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("【解密错误】创建输出文件失败:", err)
		return err
	}
	defer outputFile.Close()

	// 复制解密后的内容到输出文件
	_, err = io.Copy(outputFile, fileInZip)
	if err != nil {
		fmt.Println("【解密错误】写入解密数据失败:", err)
		return err
	}

	fmt.Println("【解密】文件解密成功:", outputPath)
	return nil
}

// EncryptFileInPlace 原地加密文件
func EncryptFileInPlace(filePath, password string) error {
	fmt.Println("\n【原地加密】开始原地加密文件:", filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("【原地加密错误】文件不存在:", filePath)
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	// 获取文件的绝对路径
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Println("【原地加密错误】获取绝对路径失败:", err)
		return err
	}
	fmt.Println("【原地加密】文件绝对路径:", absPath)

	// 创建临时文件
	tempFile := filePath + ".zip.tmp"
	fmt.Println("【原地加密】临时文件路径:", tempFile)

	// 加密到临时文件
	if err := EncryptFile(filePath, tempFile, password); err != nil {
		fmt.Println("【原地加密错误】加密到临时文件失败:", err)
		os.Remove(tempFile) // 清理临时文件
		return err
	}

	// 删除原文件
	fmt.Println("【原地加密】删除原文件:", filePath)
	if err := os.Remove(filePath); err != nil {
		fmt.Println("【原地加密错误】删除原文件失败:", err)
		os.Remove(tempFile) // 清理临时文件
		return err
	}

	// 生成新的文件路径，添加.zip后缀
	newFilePath := filePath + ".zip"
	fmt.Println("【原地加密】重命名临时文件:", tempFile, "->", newFilePath)

	// 重命名临时文件为带.zip后缀的文件
	if err := os.Rename(tempFile, newFilePath); err != nil {
		fmt.Println("【原地加密错误】重命名临时文件失败:", err)
		return err
	}

	fmt.Println("【原地加密】文件原地加密完成:", newFilePath)

	// 返回新的文件路径，以便调用者更新数据库
	return nil
}

// DecryptFileInPlace 原地解密文件
func DecryptFileInPlace(filePath, password string) error {
	fmt.Println("\n【原地解密】开始原地解密文件:", filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("【原地解密错误】文件不存在:", filePath)
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	// 获取文件的绝对路径
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Println("【原地解密错误】获取绝对路径失败:", err)
		return err
	}
	fmt.Println("【原地解密】文件绝对路径:", absPath)

	// 创建临时文件
	tempFile := filePath + ".tmp"
	fmt.Println("【原地解密】临时文件路径:", tempFile)

	// 解密到临时文件
	if err := DecryptFile(filePath, tempFile, password); err != nil {
		fmt.Println("【原地解密错误】解密到临时文件失败:", err)
		os.Remove(tempFile) // 清理临时文件
		return err
	}

	// 删除原文件
	fmt.Println("【原地解密】删除原文件:", filePath)
	if err := os.Remove(filePath); err != nil {
		fmt.Println("【原地解密错误】删除原文件失败:", err)
		os.Remove(tempFile) // 清理临时文件
		return err
	}

	// 生成新的文件路径，移除.zip后缀
	newFilePath := filePath
	if strings.HasSuffix(filePath, ".zip") {
		newFilePath = strings.TrimSuffix(filePath, ".zip")
	}

	fmt.Println("【原地解密】重命名临时文件:", tempFile, "->", newFilePath)

	// 重命名临时文件为不带.zip后缀的文件
	if err := os.Rename(tempFile, newFilePath); err != nil {
		fmt.Println("【原地解密错误】重命名临时文件失败:", err)
		return err
	}

	fmt.Println("【原地解密】文件原地解密完成:", newFilePath)
	return nil
}
