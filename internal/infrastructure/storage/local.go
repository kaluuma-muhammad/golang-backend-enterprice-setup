package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type LocalStorage struct {
	BasePath string
	BaseURL  string
}

func NewLocalStorage(basePath, baseURL string) *LocalStorage {
	return &LocalStorage{
		BasePath: basePath,
		BaseURL:  baseURL,
	}
}

func (s *LocalStorage) Upload(input UploadInput) (*File, error) {
	ext := filepath.Ext(input.Filename)

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Normalize folder
	folder := strings.Trim(input.Folder, "/")

	// Full relative path (store in DB)
	relativePath := filepath.Join(folder, filename)

	// Absolute path (disk)
	fullPath := filepath.Join(s.BasePath, relativePath)

	// Ensure directory exists
	err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
	if err != nil {
		return nil, err
	}

	out, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	size, err := io.Copy(out, input.File)
	if err != nil {
		return nil, err
	}

	return &File{
		Path:        relativePath,
		URL:         s.GetURL(relativePath),
		Filename:    filename,
		Size:        size,
		ContentType: input.ContentType,
	}, nil
}

func (s *LocalStorage) Delete(path string) error {
	if path == "" {
		return nil
	}

	fullPath := filepath.Join(s.BasePath, path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(fullPath)
}

func (s *LocalStorage) GetURL(path string) string {
	return fmt.Sprintf("%s/%s", s.BaseURL, path)
}
