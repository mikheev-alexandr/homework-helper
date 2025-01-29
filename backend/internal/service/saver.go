package service

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type FileSaverStruct struct {}

func (f FileSaverStruct) SaveFile(c *gin.Context, file *multipart.FileHeader, uploadDir string) (string, error) {
	openedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer openedFile.Close()

	buffer := make([]byte, 512)
	if _, err := openedFile.Read(buffer); err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	allowedTypes := map[string]bool{
		"application/octet-stream": true,
		"application/pdf":          true,
		"image/jpeg":               true,
		"image/png":                true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
		"application/vnd.ms-excel": true,
		"application/zip":          true,
	}
	isValid := allowedTypes[contentType]

	if !isValid {
		return "", fmt.Errorf("unsupported file type: %s", contentType)
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filePath := filepath.Join(uploadDir, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}

	return filePath, nil
}
