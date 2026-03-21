package util

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func IsValidImageType(contentType string) bool {
	allowed := map[string]struct{}{
		"image/jpeg": {},
		"image/png":  {},
		"image/gif": {},
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	mediaType = strings.ToLower(strings.TrimSpace(mediaType))

	_, ok := allowed[mediaType]
	return ok
}

func SaveThumbnail(file multipart.File, header *multipart.FileHeader) (string, error) {
    ext := filepath.Ext(header.Filename)
    filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dir := filepath.Join("static", "assets", "thumbnails")
    path := filepath.Join("static", "assets", "thumbnails", filename)

	if err := os.MkdirAll(dir, 0755); err != nil {
        return "", fmt.Errorf("failed to create upload directory: %w", err)
    }

    dst, err := os.Create(path)
    if err != nil {
        return "", err
    }
    defer dst.Close()

     if _, err := io.Copy(dst, file); err != nil {
        return "", fmt.Errorf("failed to write file: %w", err)
    }

    return path, nil
}