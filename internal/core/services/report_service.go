// services/report_service.go
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// reportService implementa la lógica de negocio para reportes
type reportService struct {
	reportRepo   ports.IReportRepository
	excelService ports.IFileService
}

// NewReportService crea una nueva instancia de ReportService
func NewReportService(reportRepo ports.IReportRepository, excelService ports.IFileService) ports.IReportService {
	return &reportService{
		reportRepo:   reportRepo,
		excelService: excelService,
	}
}

// GetDashboardReport obtiene los datos principales del dashboard
func (s *reportService) GetDashboardReport(ctx context.Context, filters *domain.ReportFilters) (*domain.DashboardReport, error) {
	if err := s.ValidateFilters(filters); err != nil {
		return nil, err
	}

	report, err := s.reportRepo.GetDashboardData(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error al generar reporte de dashboard: %w", err)
	}

	report.GeneratedAt = time.Now()
	return report, nil
}

// GetPatientsByLocalityReport obtiene pacientes agrupados por localidad
func (s *reportService) GetPatientsByLocalityReport(ctx context.Context, filters *domain.ReportFilters) (*domain.PatientsByLocalityReport, error) {
	if err := s.ValidateFilters(filters); err != nil {
		return nil, err
	}

	report, err := s.reportRepo.GetPatientsByLocality(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error al generar reporte por localidad: %w", err)
	}

	report.GeneratedAt = time.Now()
	return report, nil
}

// GetRecentMeasurementsReport obtiene las mediciones más recientes
func (s *reportService) GetRecentMeasurementsReport(ctx context.Context, filters *domain.ReportFilters) (*domain.RecentMeasurementsReport, error) {
	if err := s.ValidateFilters(filters); err != nil {
		return nil, err
	}

	report, err := s.reportRepo.GetRecentMeasurements(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error al generar reporte de mediciones recientes: %w", err)
	}

	report.GeneratedAt = time.Now()
	return report, nil
}

// GetRiskPatientsReport obtiene pacientes en riesgo
func (s *reportService) GetRiskPatientsReport(ctx context.Context, filters *domain.ReportFilters) (*domain.RiskPatientsReport, error) {
	if err := s.ValidateFilters(filters); err != nil {
		return nil, err
	}

	report, err := s.reportRepo.GetRiskPatients(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error al generar reporte de pacientes en riesgo: %w", err)
	}

	report.GeneratedAt = time.Now()
	return report, nil
}

// GetRiskPatientsReportExcel obtiene pacientes en riesgo y genera reporte Excel
func (s *reportService) GetRiskPatientsReportExcel(ctx context.Context, filters *domain.ReportFilters) ([]byte, error) {
	if err := s.ValidateFilters(filters); err != nil {
		return nil, err
	}

	// Obtener datos de pacientes en riesgo
	report, err := s.reportRepo.GetRiskPatients(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error al generar reporte de pacientes en riesgo: %w", err)
	}

	report.GeneratedAt = time.Now()

	// Generar archivo Excel
	excelData, err := s.excelService.GenerateRiskPatientsReport(ctx, report)
	if err != nil {
		return nil, fmt.Errorf("error al generar archivo Excel: %w", err)
	}

	return excelData, nil
}

// GetRiskPatientsCoordinates obtiene coordenadas de pacientes en riesgo
func (s *reportService) GetRiskPatientsCoordinates(ctx context.Context, filters *domain.ReportFilters) ([][]float64, error) {
	if err := s.ValidateFilters(filters); err != nil {
		return nil, err
	}

	coordinates, err := s.reportRepo.GetRiskPatientsCoordinates(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error al obtener coordenadas de pacientes en riesgo: %w", err)
	}

	return coordinates, nil
}

// GetUserActivityReport obtiene la actividad de usuarios
func (s *reportService) GetUserActivityReport(ctx context.Context, filters *domain.ReportFilters) (*domain.UserActivityReport, error) {
	if err := s.ValidateFilters(filters); err != nil {
		return nil, err
	}

	report, err := s.reportRepo.GetUserActivity(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error al generar reporte de actividad de usuarios: %w", err)
	}

	report.GeneratedAt = time.Now()
	return report, nil
}

// ValidateFilters valida los filtros de entrada
func (s *reportService) ValidateFilters(filters *domain.ReportFilters) error {
	if filters == nil {
		return nil // Los filtros son opcionales
	}

	// Validar días (máximo 365)
	if filters.Days > 365 {
		return fmt.Errorf("el filtro de días no puede ser mayor a 365")
	}
	if filters.Days < 0 {
		return fmt.Errorf("el filtro de días no puede ser negativo")
	}

	// Validar límite (máximo 1000)
	if filters.Limit > 1000 {
		return fmt.Errorf("el límite no puede ser mayor a 1000")
	}
	if filters.Limit < 0 {
		return fmt.Errorf("el límite no puede ser negativo")
	}

	return nil
}
