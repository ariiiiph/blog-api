package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

const MaxImageSize = 5 << 20 // 5 MB

var allowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// ValidateImageFile checks the uploaded file's extension and size.
func ValidateImageFile(fileHeader *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	if !allowedImageExtensions[ext] {
		return "", errors.New("unsupported file type, allowed: jpg, jpeg, png, webp")
	}

	if fileHeader.Size > MaxImageSize {
		return "", errors.New("file too large, max 5MB")
	}

	return ext, nil
}

// GenerateImageFilename builds a timestamp-based filename with the given extension.
func GenerateImageFilename(ext string) string {
	return fmt.Sprintf("%d%s", time.Now().UnixMilli(), ext)
}
