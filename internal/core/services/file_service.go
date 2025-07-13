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
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
	"github.com/xuri/excelize/v2"
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

func (s *FileService) GenerateApoderadosReport(ctx context.Context, users []*domain.User) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	// Crear hojas
	if err := s.createApoderadosSheet(f, users); err != nil {
		return nil, fmt.Errorf("error creando hoja de apoderados: %w", err)
	}

	if err := s.createPatientsSheet(f, users); err != nil {
		return nil, fmt.Errorf("error creando hoja de pacientes: %w", err)
	}

	if err := s.createMeasurementsSheet(f, users); err != nil {
		return nil, fmt.Errorf("error creando hoja de mediciones: %w", err)
	}

	// Eliminar la hoja por defecto
	f.DeleteSheet("Sheet1")

	// Generar el archivo
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("error generando archivo Excel: %w", err)
	}

	return buffer.Bytes(), nil
}

func (s *FileService) createApoderadosSheet(f *excelize.File, users []*domain.User) error {
	sheetName := "Apoderados"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)

	// Headers
	headers := []string{"ID", "Nombre", "Apellido", "Username", "Email", "DNI", "Teléfono", "Activo", "Rol", "Localidad", "Total Pacientes"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Estilo para headers
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
	})
	f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), style)

	// Datos
	for i, user := range users {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), user.ID.String())
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), user.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), user.LastName)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), user.Username)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), user.Email)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), user.DNI)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), user.Phone)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), user.Active)
		if user.Role.ID != uuid.Nil {
			f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), user.Role.Name)
		}
		if user.Locality != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), user.Locality.Name)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), len(user.Patients))
	}

	// Ajustar ancho de columnas
	columns := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}
	for _, col := range columns {
		f.SetColWidth(sheetName, col, col, 15)
	}

	return nil
}

func (s *FileService) createPatientsSheet(f *excelize.File, users []*domain.User) error {
	sheetName := "Pacientes"
	_, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Headers
	headers := []string{"Paciente ID", "Nombre", "Apellido", "Género", "Edad", "DNI", "Fecha Nacimiento",
		"Talla Brazo", "Peso", "Talla", "Consentimiento", "Fecha Consentimiento", "Descripción",
		"Apoderado", "Localidad", "Total Mediciones"}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Estilo para headers
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
	})
	f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), style)

	// Datos
	row := 2
	for _, user := range users {
		for _, patient := range user.Patients {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), patient.ID.String())
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), patient.Name)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), patient.Lastname)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), patient.Gender)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), patient.Age)
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), patient.DNI)
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), patient.BirthDate)
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), patient.ArmSize)
			f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), patient.Weight)
			f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), patient.Size)
			f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), patient.ConsentGiven)
			if !patient.ConsentDate.IsZero() {
				f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), patient.ConsentDate.Format("2006-01-02"))
			}
			f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), patient.Description)
			f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), user.Name+" "+user.LastName)
			if user.Locality != nil {
				f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), user.Locality.Name)
			}
			f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), len(patient.Measurements))
			row++
		}
	}

	// Ajustar ancho de columnas
	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		f.SetColWidth(sheetName, col, col, 15)
	}

	return nil
}

func (s *FileService) createMeasurementsSheet(f *excelize.File, users []*domain.User) error {
	sheetName := "Mediciones"
	_, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Headers
	headers := []string{"Medición ID", "Valor MUAC", "Descripción", "Fecha", "Paciente", "Apoderado",
		"Localidad", "Tag", "Color Tag", "Prioridad Tag", "Recomendación", "Umbral", "Código MUAC"}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Estilo para headers
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
	})
	f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), style)

	// Datos
	row := 2
	for _, user := range users {
		for _, patient := range user.Patients {
			for _, measurement := range patient.Measurements {
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), measurement.ID.String())
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), measurement.MuacValue)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), measurement.Description)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), measurement.CreatedAt.Format("2006-01-02 15:04:05"))
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), patient.Name+" "+patient.Lastname)
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), user.Name+" "+user.LastName)
				if user.Locality != nil {
					f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), user.Locality.Name)
				}
				if measurement.Tag != nil {
					f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), measurement.Tag.Name)
					f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), measurement.Tag.Color)
					f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), measurement.Tag.Priority)
				}
				if measurement.Recommendation != nil {
					f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), measurement.Recommendation.Name)
					f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), measurement.Recommendation.RecommendationUmbral)
					f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), measurement.Recommendation.MuacCode)
				}
				row++
			}
		}
	}

	// Ajustar ancho de columnas
	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		f.SetColWidth(sheetName, col, col, 18)
	}

	return nil
}
