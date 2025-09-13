package utils

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ValidateImageFile(header *multipart.FileHeader, maxSize int64, allowedTypes []string) error {
	if header.Size > maxSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d bytes", header.Size, maxSize)
	}

	if header.Size == 0 {
		return fmt.Errorf("file is empty")
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		mimeType = GetMimeTypeFromExtension(ext)
	}

	if !IsAllowedMimeType(mimeType, allowedTypes) {
		return fmt.Errorf("file type %s is not allowed. Allowed types: %v", mimeType, allowedTypes)
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		return fmt.Errorf("file must have an extension")
	}

	return nil
}

func DetectMimeType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read file for mime type detection: %w", err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("failed to seek file: %w", err)
	}

	return http.DetectContentType(buffer), nil
}

func IsAllowedMimeType(mimeType string, allowedTypes []string) bool {
	for _, allowed := range allowedTypes {
		if strings.EqualFold(strings.TrimSpace(mimeType), strings.TrimSpace(allowed)) {
			return true
		}
	}
	return false
}

func GetMimeTypeFromExtension(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".tiff", ".tif":
		return "image/tiff"
	default:
		return "application/octet-stream"
	}
}

func GetExtensionFromMimeType(mimeType string) string {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/jpeg":
		return ".jpg"
	case "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/bmp":
		return ".bmp"
	case "image/tiff":
		return ".tiff"
	default:
		return ".jpg"
	}
}

func EnsureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func SanitizeFilename(filename string) string {
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")

	if len(filename) > 255 {
		ext := filepath.Ext(filename)
		name := filename[:255-len(ext)]
		filename = name + ext
	}

	return filename
}