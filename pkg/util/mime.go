package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 判断文件是否为支持的类型
func IsSupportedFileType(file multipart.File) (bool, string) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return false, ""
	}
	filetype := http.DetectContentType(buffer)
	file.Seek(0, 0) // Reset the file pointer
	return isValidFileType(filetype), filetype
}

func isValidFileType(filetype string) bool {
	supportedTypes := []string{
		"image/jpeg",
		"image/png",
		"video/mp4",
	}
	for _, t := range supportedTypes {
		if t == filetype {
			return true
		}
	}
	return false
}

// IsVideoFile 检查文件是否为视频文件
func IsVideoFile(filename string) bool {
	// 获取文件扩展名并转换为小写
	ext := strings.ToLower(filepath.Ext(filename))

	// 视频文件扩展名列表
	videoExtensions := map[string]bool{
		".mp4":  true,
		".mov":  true,
		".avi":  true,
		".wmv":  true,
		".flv":  true,
		".mkv":  true,
		".webm": true,
		".mpeg": true,
		".mpg":  true,
		".3gp":  true,
		".ogv":  true,
		".m4v":  true,
	}

	// 检查扩展名是否在视频扩展名列表中
	if _, exists := videoExtensions[ext]; exists {
		return true
	}

	// 检测文件的 MIME 类型（可选）
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false
	}
	defer file.Close()
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return false
	}
	mimeType := http.DetectContentType(buf)
	return strings.HasPrefix(mimeType, "video/")
}

// IsImageFile 检查文件是否为图片文件
func IsImageFile(filename string) bool {
	// 获取文件扩展名并转换为小写
	ext := strings.ToLower(filepath.Ext(filename))
	// 图片文件扩展名列表
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
		".heic": true,
	}
	// 检查扩展名是否在图片扩展名列表中
	if _, exists := imageExtensions[ext]; exists {
		return true
	}
	return false
}

// SaveUploadedFile uploads the form file to specific dst.
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// IsAudioFile 检查文件是否为音频文件
func IsAudioFile(filename string) bool {
	// 获取文件扩展名并转换为小写
	ext := strings.ToLower(filepath.Ext(filename))
	// 音频文件扩展名列表
	audioExtensions := map[string]bool{
		".mp3":  true,
		".wav":  true,
		".ogg":  true,
		".flac": true,
		".aac":  true,
		".m4a":  true,
		".wma":  true,
	}
	// 检查扩展名是否在音频扩展名列表中
	if _, exists := audioExtensions[ext]; exists {
		return true
	}

	// 如果扩展名检查不通过，进一步检查文件内容
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false
	}
	defer file.Close()

	// 读取文件前512字节
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return false
	}

	// 使用http.DetectContentType检测MIME类型
	contentType := http.DetectContentType(buffer)
	return strings.HasPrefix(contentType, "audio/")
}
