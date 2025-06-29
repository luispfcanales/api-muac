package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

type FileService struct {
	uploadPath   string
	baseURL      string
	maxSize      int64
	allowedTypes map[string]bool
}

// NewFileService crea una nueva instancia del servicio de archivos
func NewFileService(uploadPath, baseURL string) ports.IFileService {
	// Tipos de archivo permitidos
	allowedTypes := map[string]bool{
		"image/jpeg":      true,
		"image/jpg":       true,
		"image/png":       true,
		"image/gif":       true,
		"application/pdf": true,
		"text/plain":      true,
	}

	return &FileService{
		uploadPath:   uploadPath,
		baseURL:      baseURL,
		maxSize:      10 * 1024 * 1024, // 10MB máximo
		allowedTypes: allowedTypes,
	}
}

// UploadFile sube un archivo al servidor
func (fs *FileService) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (*ports.FileInfo, error) {
	// Validar archivo
	if err := fs.ValidateFile(header); err != nil {
		return nil, err
	}

	// Crear directorio si no existe
	folderPath := filepath.Join(fs.uploadPath, folder)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return nil, fmt.Errorf("error al crear directorio: %v", err)
	}

	// Generar nombre único para el archivo
	fileID := uuid.New().String()
	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s%s", fileID, ext)
	filePath := filepath.Join(folderPath, fileName)

	// Crear archivo en el servidor
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("error al crear archivo: %v", err)
	}
	defer dst.Close()

	// Copiar contenido del archivo
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath) // Limpiar en caso de error
		return nil, fmt.Errorf("error al copiar archivo: %v", err)
	}

	// Obtener información del archivo
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("error al obtener información del archivo: %v", err)
	}

	// Crear FileInfo
	info := &ports.FileInfo{
		ID:           fileID,
		FileName:     fileName,
		OriginalName: header.Filename,
		Size:         fileInfo.Size(),
		ContentType:  header.Header.Get("Content-Type"),
		Path:         filePath,
		URL:          fmt.Sprintf("%s/files/%s/%s", fs.baseURL, folder, fileName),
		UploadedAt:   time.Now().Format(time.RFC3339),
	}

	// Guardar metadata del archivo
	if err := fs.saveFileMetadata(info, folder); err != nil {
		return nil, fmt.Errorf("error al guardar metadata: %v", err)
	}

	return info, nil
}

// GetFile obtiene información de un archivo por su ID
func (fs *FileService) GetFile(ctx context.Context, fileID string) (*ports.FileInfo, error) {
	// Buscar en todas las carpetas posibles
	folders := []string{"patients", "documents", "images", "uploads"}

	for _, folder := range folders {
		metadataPath := filepath.Join(fs.uploadPath, folder, "metadata", fmt.Sprintf("%s.json", fileID))
		if info, err := fs.loadFileMetadata(metadataPath); err == nil {
			return info, nil
		}
	}

	return nil, fmt.Errorf("archivo no encontrado")
}

// GetFileContent obtiene el contenido de un archivo
func (fs *FileService) GetFileContent(ctx context.Context, fileID string) (io.ReadCloser, error) {
	info, err := fs.GetFile(ctx, fileID)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(info.Path)
	if err != nil {
		return nil, fmt.Errorf("error al abrir archivo: %v", err)
	}

	return file, nil
}

// DeleteFile elimina un archivo del servidor
func (fs *FileService) DeleteFile(ctx context.Context, fileID string) error {
	info, err := fs.GetFile(ctx, fileID)
	if err != nil {
		return err
	}

	// Eliminar archivo físico
	if err := os.Remove(info.Path); err != nil {
		return fmt.Errorf("error al eliminar archivo: %v", err)
	}

	// Eliminar metadata
	folder := filepath.Dir(filepath.Dir(info.Path))
	folder = filepath.Base(folder)
	metadataPath := filepath.Join(fs.uploadPath, folder, "metadata", fmt.Sprintf("%s.json", fileID))
	os.Remove(metadataPath)

	return nil
}

// GetFilesByFolder obtiene todos los archivos de una carpeta
func (fs *FileService) GetFilesByFolder(ctx context.Context, folder string) ([]*ports.FileInfo, error) {
	metadataDir := filepath.Join(fs.uploadPath, folder, "metadata")

	files, err := os.ReadDir(metadataDir)
	if err != nil {
		return nil, fmt.Errorf("error al leer directorio: %v", err)
	}

	var fileInfos []*ports.FileInfo
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			metadataPath := filepath.Join(metadataDir, file.Name())
			if info, err := fs.loadFileMetadata(metadataPath); err == nil {
				fileInfos = append(fileInfos, info)
			}
		}
	}

	return fileInfos, nil
}

// ValidateFile valida si un archivo es válido
func (fs *FileService) ValidateFile(header *multipart.FileHeader) error {
	// Validar tamaño
	if header.Size > fs.maxSize {
		return fmt.Errorf("archivo demasiado grande. Máximo permitido: %d bytes", fs.maxSize)
	}

	// Validar tipo de contenido
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		// Intentar determinar por extensión
		ext := strings.ToLower(filepath.Ext(header.Filename))
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		case ".pdf":
			contentType = "application/pdf"
		case ".txt":
			contentType = "text/plain"
		default:
			return fmt.Errorf("tipo de archivo no soportado: %s", ext)
		}
	}

	if !fs.allowedTypes[contentType] {
		return fmt.Errorf("tipo de archivo no permitido: %s", contentType)
	}

	return nil
}

// saveFileMetadata guarda la metadata del archivo
func (fs *FileService) saveFileMetadata(info *ports.FileInfo, folder string) error {
	metadataDir := filepath.Join(fs.uploadPath, folder, "metadata")
	if err := os.MkdirAll(metadataDir, 0755); err != nil {
		return err
	}

	metadataPath := filepath.Join(metadataDir, fmt.Sprintf("%s.json", info.ID))
	file, err := os.Create(metadataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(info)
}

// loadFileMetadata carga la metadata del archivo
func (fs *FileService) loadFileMetadata(metadataPath string) (*ports.FileInfo, error) {
	file, err := os.Open(metadataPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var info ports.FileInfo
	if err := json.NewDecoder(file).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}
