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
		baseURL:      baseURL,          // Asegúrate de pasar https://nutriradar.unamad.edu.pe aquí
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

	// Crear FileInfo con la URL correcta
	info := &ports.FileInfo{
		ID:           fileID,
		FileName:     fileName,
		OriginalName: header.Filename,
		Size:         fileInfo.Size(),
		ContentType:  header.Header.Get("Content-Type"),
		Path:         filePath,
		URL:          fmt.Sprintf("%s/files/%s/%s", fs.baseURL, folder, fileName), // Aquí se usa baseURL
		UploadedAt:   time.Now().Format(time.RFC3339),
	}

	// Guardar metadata del archivo
	if err := fs.saveFileMetadata(info, folder); err != nil {
		return nil, fmt.Errorf("error al guardar metadata: %v", err)
	}

	return info, nil
}

// GetFile obtiene información de un archivo por su ID - MEJORADO
func (fs *FileService) GetFile(ctx context.Context, fileID string) (*ports.FileInfo, error) {
	// Estructura específica para tu caso: uploads/patients/dni/metadata/
	metadataPath := filepath.Join(fs.uploadPath, "patients", "dni", "metadata", fmt.Sprintf("%s.json", fileID))

	if info, err := fs.loadFileMetadata(metadataPath); err == nil {
		return info, nil
	}

	// Si no se encuentra, buscar en otras ubicaciones posibles como fallback
	folders := []string{
		"patients/dni",
		"patients/documents",
		"patients/images",
		"documents",
		"images",
		"uploads",
	}

	for _, folder := range folders {
		metadataPath := filepath.Join(fs.uploadPath, folder, "metadata", fmt.Sprintf("%s.json", fileID))
		if info, err := fs.loadFileMetadata(metadataPath); err == nil {
			return info, nil
		}
	}

	return nil, fmt.Errorf("archivo no encontrado: %s", fileID)
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

// DeleteFile elimina un archivo del servidor - MEJORADO
func (fs *FileService) DeleteFile(ctx context.Context, fileID string) error {
	// Obtener información del archivo
	info, err := fs.GetFile(ctx, fileID)
	if err != nil {
		return fmt.Errorf("archivo no encontrado para eliminar: %s", fileID)
	}

	// Eliminar archivo físico
	if err := os.Remove(info.Path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error al eliminar archivo físico %s: %v", info.Path, err)
	}

	// Construir ruta de metadata basada en la estructura conocida
	// Para uploads/patients/dni/archivo.jpg -> uploads/patients/dni/metadata/uuid.json
	var metadataPath string

	// Detectar el tipo de archivo basado en la ruta
	if filepath.Dir(info.Path) == filepath.Join(fs.uploadPath, "patients", "dni") {
		metadataPath = filepath.Join(fs.uploadPath, "patients", "dni", "metadata", fmt.Sprintf("%s.json", fileID))
	} else {
		// Para otros tipos de archivos, intentar extraer la carpeta padre
		parentDir := filepath.Dir(info.Path)
		metadataPath = filepath.Join(parentDir, "metadata", fmt.Sprintf("%s.json", fileID))
	}

	// Eliminar metadata (no fallar si no existe)
	if err := os.Remove(metadataPath); err != nil && !os.IsNotExist(err) {
		// Log pero no fallar por metadata
		fmt.Printf("Warning: no se pudo eliminar metadata %s: %v\n", metadataPath, err)
	}

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

// FileExists verifica si un archivo existe - MÉTODO NUEVO
func (fs *FileService) FileExists(ctx context.Context, fileID string) bool {
	_, err := fs.GetFile(ctx, fileID)
	return err == nil
}

// DeleteFileIfExists elimina un archivo si existe - MÉTODO NUEVO
func (fs *FileService) DeleteFileIfExists(ctx context.Context, fileID string) error {
	if !fs.FileExists(ctx, fileID) {
		return nil // No hacer nada si el archivo no existe
	}
	return fs.DeleteFile(ctx, fileID)
}

//para generar el excel
//----------------------------------------------------------------------------------------------

// ============= AGREGAR AL FILE SERVICE =============

// GenerateRiskPatientsReport genera un reporte Excel de pacientes en riesgo
func (s *FileService) GenerateRiskPatientsReport(ctx context.Context, report *domain.RiskPatientsReport) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	// Crear hojas
	if err := s.createRiskSummarySheet(f, report); err != nil {
		return nil, fmt.Errorf("error creando hoja resumen: %w", err)
	}

	if err := s.createSevereCasesSheet(f, report.SevereCases); err != nil {
		return nil, fmt.Errorf("error creando hoja casos severos: %w", err)
	}

	if err := s.createModerateCasesSheet(f, report.ModerateCases); err != nil {
		return nil, fmt.Errorf("error creando hoja casos moderados: %w", err)
	}

	if err := s.createAllRiskPatientsSheet(f, report); err != nil {
		return nil, fmt.Errorf("error creando hoja todos los pacientes: %w", err)
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

// createRiskSummarySheet crea la hoja de resumen
func (s *FileService) createRiskSummarySheet(f *excelize.File, report *domain.RiskPatientsReport) error {
	sheetName := "Resumen"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)

	// Estilo para títulos
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	// Estilo para datos críticos
	criticalStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"DC3545"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	// Estilo para datos moderados
	moderateStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "000000"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"FFC107"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	// Título principal
	f.SetCellValue(sheetName, "A1", "REPORTE DE PACIENTES EN RIESGO NUTRICIONAL")
	f.MergeCell(sheetName, "A1", "D1")
	f.SetCellStyle(sheetName, "A1", "D1", titleStyle)

	// Información del reporte
	f.SetCellValue(sheetName, "A3", "Fecha de generación:")
	f.SetCellValue(sheetName, "B3", report.GeneratedAt.Format("2006-01-02 15:04:05"))

	// Resumen estadístico
	f.SetCellValue(sheetName, "A5", "RESUMEN ESTADÍSTICO")
	f.SetCellStyle(sheetName, "A5", "A5", titleStyle)

	f.SetCellValue(sheetName, "A7", "Casos Severos (MUAC < 11.5 cm)")
	f.SetCellValue(sheetName, "B7", len(report.SevereCases))
	f.SetCellStyle(sheetName, "B7", "B7", criticalStyle)

	f.SetCellValue(sheetName, "A8", "Casos Moderados (MUAC 11.5-12.4 cm)")
	f.SetCellValue(sheetName, "B8", len(report.ModerateCases))
	f.SetCellStyle(sheetName, "B8", "B8", moderateStyle)

	f.SetCellValue(sheetName, "A9", "Total Pacientes en Riesgo")
	f.SetCellValue(sheetName, "B9", len(report.SevereCases)+len(report.ModerateCases))

	// Distribución por localidad
	f.SetCellValue(sheetName, "A11", "DISTRIBUCIÓN POR LOCALIDAD")
	f.SetCellStyle(sheetName, "A11", "A11", titleStyle)

	// Contar por localidad
	localityCount := make(map[string]int)
	allPatients := append(report.SevereCases, report.ModerateCases...)
	for _, patient := range allPatients {
		localityCount[patient.LocalityName]++
	}

	row := 13
	f.SetCellValue(sheetName, "A12", "Localidad")
	f.SetCellValue(sheetName, "B12", "Pacientes")
	f.SetCellStyle(sheetName, "A12", "B12", titleStyle)

	for locality, count := range localityCount {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), locality)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), count)
		row++
	}

	// Ajustar ancho de columnas
	f.SetColWidth(sheetName, "A", "A", 25)
	f.SetColWidth(sheetName, "B", "B", 15)

	return nil
}

// createSevereCasesSheet crea la hoja de casos severos
func (s *FileService) createSevereCasesSheet(f *excelize.File, severeCases []domain.RiskPatient) error {
	sheetName := "Casos Severos"
	_, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Headers
	headers := []string{"ID Paciente", "Nombre Paciente", "Edad", "Género", "Valor MUAC",
		"Código MUAC", "Localidad", "Apoderado", "Última Medición", "Días Transcurridos"}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Estilo para headers (crítico)
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"DC3545"}, Pattern: 1},
	})
	f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), style)

	// Datos
	for i, patient := range severeCases {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), patient.PatientID.String())
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), patient.PatientName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), patient.Age)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), patient.Gender)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), patient.MuacValue)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), patient.MuacCode)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), patient.LocalityName)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), patient.UserName)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), patient.LastMeasure.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), patient.DaysAgo)
	}

	// Ajustar ancho de columnas
	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		f.SetColWidth(sheetName, col, col, 15)
	}

	return nil
}

// createModerateCasesSheet crea la hoja de casos moderados
func (s *FileService) createModerateCasesSheet(f *excelize.File, moderateCases []domain.RiskPatient) error {
	sheetName := "Casos Moderados"
	_, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Headers
	headers := []string{"ID Paciente", "Nombre Paciente", "Edad", "Género", "Valor MUAC",
		"Código MUAC", "Localidad", "Apoderado", "Última Medición", "Días Transcurridos"}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Estilo para headers (moderado)
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "000000"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FFC107"}, Pattern: 1},
	})
	f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), style)

	// Datos
	for i, patient := range moderateCases {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), patient.PatientID.String())
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), patient.PatientName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), patient.Age)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), patient.Gender)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), patient.MuacValue)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), patient.MuacCode)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), patient.LocalityName)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), patient.UserName)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), patient.LastMeasure.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), patient.DaysAgo)
	}

	// Ajustar ancho de columnas
	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		f.SetColWidth(sheetName, col, col, 15)
	}

	return nil
}

// createAllRiskPatientsSheet crea una hoja con todos los pacientes en riesgo
func (s *FileService) createAllRiskPatientsSheet(f *excelize.File, report *domain.RiskPatientsReport) error {
	sheetName := "Todos los Pacientes"
	_, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Headers
	headers := []string{"ID Paciente", "Nombre Paciente", "Edad", "Género", "Valor MUAC",
		"Código MUAC", "Nivel Riesgo", "Localidad", "Apoderado", "Última Medición", "Días Transcurridos"}

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

	// Combinar datos y agregar tipo de riesgo
	allPatients := []domain.RiskPatient{}

	// Agregar casos severos
	for _, patient := range report.SevereCases {
		allPatients = append(allPatients, patient)
	}

	// Agregar casos moderados
	for _, patient := range report.ModerateCases {
		allPatients = append(allPatients, patient)
	}

	// Estilos para diferentes niveles de riesgo
	criticalRowStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FFEBEE"}, Pattern: 1},
	})
	moderateRowStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FFF8E1"}, Pattern: 1},
	})

	// Datos
	for i, patient := range allPatients {
		row := i + 2

		// Determinar nivel de riesgo
		riskLevel := "Moderado"
		if patient.MuacValue < 11.5 {
			riskLevel = "Severo"
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), patient.PatientID.String())
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), patient.PatientName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), patient.Age)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), patient.Gender)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), patient.MuacValue)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), patient.MuacCode)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), riskLevel)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), patient.LocalityName)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), patient.UserName)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), patient.LastMeasure.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), patient.DaysAgo)

		// Aplicar estilo según nivel de riesgo
		if riskLevel == "Severo" {
			f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("K%d", row), criticalRowStyle)
		} else {
			f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("K%d", row), moderateRowStyle)
		}
	}

	// Ajustar ancho de columnas
	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		f.SetColWidth(sheetName, col, col, 15)
	}

	return nil
}
