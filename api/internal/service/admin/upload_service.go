package admin

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type UploadService struct {
	uploadDir string
	maxSize   int64
}

func NewUploadService() *UploadService {
	return &UploadService{
		uploadDir: "storage/uploads/avatars",
		maxSize:   2 * 1024 * 1024, // 2MB
	}
}

// UploadAvatar uploads an avatar file
func (s *UploadService) UploadAvatar(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID string) (string, error) {
	// Validate file size
	if header.Size > s.maxSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.maxSize)
	}

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		return "", fmt.Errorf("invalid file type: only jpg, jpeg, png allowed")
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Create user directory
	userDir := filepath.Join(s.uploadDir, userID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return "", err
	}

	// Create destination file
	destPath := filepath.Join(userDir, filename)
	dest, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dest.Close()

	// Copy file
	if _, err := io.Copy(dest, file); err != nil {
		return "", err
	}

	// Return relative path
	return filepath.Join(userID, filename), nil
}

// DeleteAvatar deletes an avatar file
func (s *UploadService) DeleteAvatar(ctx context.Context, avatarPath string) error {
	fullPath := filepath.Join(s.uploadDir, avatarPath)
	return os.Remove(fullPath)
}

// isValidImageType checks if content type is valid image
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
	}

	contentType = strings.ToLower(contentType)
	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}

	return false
}
