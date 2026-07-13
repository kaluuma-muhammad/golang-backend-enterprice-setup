package storage

import (
	"errors"
	"mime/multipart"
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

func ValidateImage(fileHeader *multipart.FileHeader) error {
	if fileHeader.Size > 5*1024*1024 {
		return errors.New("file too large (max 5MB)")
	}
	return nil
}
