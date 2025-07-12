// ports/report.go
package ports

import (
	"context"

	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// IReportRepository define las consultas para generar reportes
type IReportRepository interface {
	// Dashboard principal
	GetDashboardData(ctx context.Context, filters *domain.ReportFilters) (*domain.DashboardReport, error)

	// Pacientes por localidad
	GetPatientsByLocality(ctx context.Context, filters *domain.ReportFilters) (*domain.PatientsByLocalityReport, error)

	// Mediciones recientes
	GetRecentMeasurements(ctx context.Context, filters *domain.ReportFilters) (*domain.RecentMeasurementsReport, error)

	// Pacientes en riesgo
	GetRiskPatients(ctx context.Context, filters *domain.ReportFilters) (*domain.RiskPatientsReport, error)

	// Actividad de usuarios
	GetUserActivity(ctx context.Context, filters *domain.ReportFilters) (*domain.UserActivityReport, error)
}

// IReportService define las operaciones del servicio para reportes
type IReportService interface {
	// Reportes principales
	GetDashboardReport(ctx context.Context, filters *domain.ReportFilters) (*domain.DashboardReport, error)
	GetPatientsByLocalityReport(ctx context.Context, filters *domain.ReportFilters) (*domain.PatientsByLocalityReport, error)
	GetRecentMeasurementsReport(ctx context.Context, filters *domain.ReportFilters) (*domain.RecentMeasurementsReport, error)
	GetRiskPatientsReport(ctx context.Context, filters *domain.ReportFilters) (*domain.RiskPatientsReport, error)
	GetUserActivityReport(ctx context.Context, filters *domain.ReportFilters) (*domain.UserActivityReport, error)

	// Validación
	ValidateFilters(filters *domain.ReportFilters) error
}
