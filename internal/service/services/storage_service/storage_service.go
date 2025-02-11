package storage_service

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type StorageServiceV1 struct {
	uploadDir string
}

func NewStorageServiceV1() *StorageServiceV1 {
	uploadDir := os.Getenv("UPLOAD_DIR")

	if uploadDir == "" {
		log.Fatal("UPLOAD_DIR environment variable is not set")
	}

	dirExist := ensureDirExists(uploadDir)

	if dirExist != nil {
		log.Fatal("Folder does not exist")
	}

	return &StorageServiceV1{
		uploadDir: uploadDir,
	}
}

// Method to save a file
func (s *StorageServiceV1) SaveFile(attachment *multipart.FileHeader) (string, error) {
	// Generate a unique file name using a UUID
	fileName := s.generateUniqueFileName(attachment.Filename)
	filePath := filepath.Join(s.uploadDir, fileName)

	// Save the uploaded file to disk
	if err := s.saveUploadedFile(attachment, filePath); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filePath, nil
}

// Helper to generate a unique file name using UUID
func (s *StorageServiceV1) generateUniqueFileName(originalFileName string) string {
	ext := filepath.Ext(originalFileName) // Get the file extension
	base := strings.TrimSuffix(originalFileName, ext)
	uuid := uuid.New().String() // Generate a unique identifier

	return fmt.Sprintf("%s_%s%s", base, uuid, ext)
}

// Helper to save the uploaded file to the file system
func (s *StorageServiceV1) saveUploadedFile(attachment *multipart.FileHeader, destinationPath string) error {
	// Open the uploaded file
	srcFile, err := attachment.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy the content from the uploaded file to the destination
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

func ensureDirExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755) // Set appropriate permissions
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}
