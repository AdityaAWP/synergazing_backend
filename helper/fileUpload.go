package helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GetUrlFile(filepath string) string {
	if filepath == "" {
		return ""
	}
	return "/" + strings.ReplaceAll(filepath, "\\", "/")
}

func UploadFile(file *multipart.FileHeader, uploadType string) (string, error) {
	if !isValidateImageType(file.Filename) {
		return "", fmt.Errorf("invalid file type. Only jpg, jpeg, png, gif are allowed")
	}

	if file.Size > 2*1024*1024 {
		return "", fmt.Errorf("file too large, maximum is only 2MB")
	}

	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)

	var uploadDir string
	switch uploadType {
	case "profile":
		uploadDir = "storage/profiles"
	case "post":
		uploadDir = "storage/posts"
	default:
		uploadDir = "storage/temp"
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	filePath := filepath.Join(uploadDir, newFilename)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open upload file: %v", err)
	}

	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}

	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return filePath, nil
}

func isValidateImageType(Filename string) bool {
	ext := strings.ToLower(filepath.Ext(Filename))
	validextension := []string{".jpg", ".jpeg", ".png"}

	for _, validExt := range validextension {
		if ext == validExt {
			return true
		}
	}
	return false
}

func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(filePath)
}
