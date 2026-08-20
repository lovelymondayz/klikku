package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Storage handles file operations on local filesystem
type Storage struct {
	basePath string
}

// NewStorage creates a new storage instance
func NewStorage(basePath string) (*Storage, error) {
	// Ensure base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("create storage dir: %w", err)
	}
	return &Storage{basePath: basePath}, nil
}

// Upload saves data to storage
func (s *Storage) Upload(bucket, objectName string, data []byte, contentType string) error {
	dir := filepath.Join(s.basePath, bucket)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create bucket dir: %w", err)
	}

	path := filepath.Join(dir, objectName)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// Download reads data from storage
func (s *Storage) Download(bucket, objectName string) ([]byte, error) {
	path := filepath.Join(s.basePath, bucket, objectName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return data, nil
}

// Delete removes a file from storage
func (s *Storage) Delete(bucket, objectName string) error {
	path := filepath.Join(s.basePath, bucket, objectName)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("remove file: %w", err)
	}
	return nil
}

// Exists checks if a file exists
func (s *Storage) Exists(bucket, objectName string) bool {
	path := filepath.Join(s.basePath, bucket, objectName)
	_, err := os.Stat(path)
	return err == nil
}

// Serve serves a file to an io.Writer
func (s *Storage) Serve(bucket, objectName string, w io.Writer) error {
	data, err := s.Download(bucket, objectName)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// GetPath returns the full path for a file
func (s *Storage) GetPath(bucket, objectName string) string {
	return filepath.Join(s.basePath, bucket, objectName)
}
