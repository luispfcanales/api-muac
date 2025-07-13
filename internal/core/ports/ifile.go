package ports

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// FileInfo contiene información sobre un archivo subido
type FileInfo struct {
	ID           string `json:"id"`
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	ContentType  string `json:"content_type"`
	Path         string `json:"path"`
	URL          string `json:"url"`
	UploadedAt   string `json:"uploaded_at"`
}

// IFileService define las operaciones del servicio de archivos
type IFileService interface {
	// UploadFile sube un archivo al servidor
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (*FileInfo, error)

	// GetFile obtiene información de un archivo por su ID
	GetFile(ctx context.Context, fileID string) (*FileInfo, error)

	// GetFileContent obtiene el contenido de un archivo
	GetFileContent(ctx context.Context, fileID string) (io.ReadCloser, error)

	// DeleteFile elimina un archivo del servidor
	DeleteFile(ctx context.Context, fileID string) error

	// GetFilesByFolder obtiene todos los archivos de una carpeta
	GetFilesByFolder(ctx context.Context, folder string) ([]*FileInfo, error)

	// ValidateFile valida si un archivo es válido (tipo, tamaño, etc.)
	ValidateFile(header *multipart.FileHeader) error

	GenerateApoderadosReport(ctx context.Context, users []*domain.User) ([]byte, error)
}
