// postgres/report_repository.go
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
	"gorm.io/gorm"
)

// reportRepository implementa la interfaz IReportRepository usando GORM
type reportRepository struct {
	db *gorm.DB
}

// NewReportRepository crea una nueva instancia de ReportRepository
func NewReportRepository(db *gorm.DB) ports.IReportRepository {
	return &reportRepository{
		db: db,
	}
}

// GetDashboardData obtiene los datos principales del dashboard
func (r *reportRepository) GetDashboardData(ctx context.Context, filters *domain.ReportFilters) (*domain.DashboardReport, error) {
	report := &domain.DashboardReport{}

	// Aplicar filtros de tiempo
	query := r.db.WithContext(ctx)
	if filters != nil && filters.Days > 0 {
		since := time.Now().AddDate(0, 0, -filters.Days)
		query = query.Where("created_at >= ?", since)
	}

	// Total de pacientes
	if err := query.Model(&domain.Patient{}).Count(&report.TotalPatients).Error; err != nil {
		return nil, fmt.Errorf("error al contar pacientes: %w", err)
	}

	// Total de mediciones
	measureQuery := r.db.WithContext(ctx)
	if filters != nil && filters.Days > 0 {
		since := time.Now().AddDate(0, 0, -filters.Days)
		measureQuery = measureQuery.Where("measurements.created_at >= ?", since)
	}

	if filters != nil && filters.LocalityID != nil {
		measureQuery = measureQuery.Joins("JOIN patients p ON measurements.patient_id = p.id").
			Joins("JOIN users u ON p.user_id = u.id").
			Where("u.locality_id = ?", *filters.LocalityID)
	}

	if err := measureQuery.Model(&domain.Measurement{}).Count(&report.TotalMeasurements).Error; err != nil {
		return nil, fmt.Errorf("error al contar mediciones: %w", err)
	}

	// Total de usuarios
	userQuery := r.db.WithContext(ctx).Model(&domain.User{})
	if filters != nil && filters.LocalityID != nil {
		userQuery = userQuery.Where("locality_id = ?", *filters.LocalityID)
	}
	if err := userQuery.Count(&report.TotalUsers).Error; err != nil {
		return nil, fmt.Errorf("error al contar usuarios: %w", err)
	}

	// Distribución por estado nutricional
	distribution, err := r.getStatusDistribution(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error al obtener distribución: %w", err)
	}
	report.StatusDistribution = *distribution

	// Pacientes en riesgo (moderado + severo)
	report.PatientsAtRisk = distribution.Moderate.Total + distribution.Severe.Total

	return report, nil
}

// GetPatientsByLocality obtiene pacientes agrupados por localidad
func (r *reportRepository) GetPatientsByLocality(ctx context.Context, filters *domain.ReportFilters) (*domain.PatientsByLocalityReport, error) {
	var localities []struct {
		LocalityID   uuid.UUID
		LocalityName string
		Total        int64
		Normal       int64
		Moderate     int64
		Severe       int64
	}

	query := r.db.WithContext(ctx).
		Select(`
			l.id as locality_id,
			l.name as locality_name,
			COUNT(DISTINCT p.id) as total,
			COUNT(CASE WHEN m.muac_value >= 12.5 THEN 1 END) as normal,
			COUNT(CASE WHEN m.muac_value >= 11.5 AND m.muac_value < 12.5 THEN 1 END) as moderate,
			COUNT(CASE WHEN m.muac_value < 11.5 THEN 1 END) as severe
		`).
		Table("localities l").
		Joins("LEFT JOIN users u ON l.id = u.locality_id").
		Joins("LEFT JOIN patients p ON u.id = p.user_id").
		Joins(`LEFT JOIN measurements m ON p.id = m.patient_id AND m.id = (
			SELECT id FROM measurements m2 
			WHERE m2.patient_id = p.id 
			ORDER BY m2.created_at DESC 
			LIMIT 1
		)`).
		Group("l.id, l.name").
		Order("l.name")

	if filters != nil {
		if filters.LocalityID != nil {
			query = query.Where("l.id = ?", *filters.LocalityID)
		}
		if filters.Days > 0 {
			since := time.Now().AddDate(0, 0, -filters.Days)
			query = query.Where("m.created_at >= ?", since)
		}
	}

	if err := query.Scan(&localities).Error; err != nil {
		return nil, fmt.Errorf("error al obtener datos por localidad: %w", err)
	}

	// Convertir a estructura de respuesta
	report := &domain.PatientsByLocalityReport{
		LocalityData: make([]domain.LocalityData, len(localities)),
	}

	for i, loc := range localities {
		total := float64(loc.Total)
		atRisk := int(loc.Moderate + loc.Severe)

		report.LocalityData[i] = domain.LocalityData{
			LocalityID:   loc.LocalityID,
			LocalityName: loc.LocalityName,
			Total:        int(loc.Total),
			AtRisk:       atRisk,
			Distribution: domain.StatusDistribution{
				Normal: domain.StatusCount{
					Total:      loc.Normal,
					Percentage: r.calculatePercentage(int(loc.Normal), total),
				},
				Moderate: domain.StatusCount{
					Total:      loc.Moderate,
					Percentage: r.calculatePercentage(int(loc.Moderate), total),
				},
				Severe: domain.StatusCount{
					Total:      loc.Severe,
					Percentage: r.calculatePercentage(int(loc.Severe), total),
				},
			},
		}
	}

	return report, nil
}

// GetRecentMeasurements obtiene las mediciones más recientes
func (r *reportRepository) GetRecentMeasurements(ctx context.Context, filters *domain.ReportFilters) (*domain.RecentMeasurementsReport, error) {
	var measurements []domain.RecentMeasurement

	query := r.db.WithContext(ctx).
		Select(`
			m.id,
			CONCAT(p.name, ' ', p.lastname) as patient_name,
			p.age as patient_age,
			m.muac_value,
			CASE 
				WHEN m.muac_value >= 12.5 THEN 'MUAC-G1'
				WHEN m.muac_value >= 11.5 THEN 'MUAC-Y1'
				ELSE 'MUAC-R1'
			END as muac_code,
			CASE 
				WHEN m.muac_value >= 12.5 THEN '#28a745'
				WHEN m.muac_value >= 11.5 THEN '#ffc107'
				ELSE '#dc3545'
			END as color_code,
			CONCAT(u.name, ' ', u.lastname) as user_name,
			l.name as locality_name,
			m.created_at
		`).
		Table("measurements m").
		Joins("JOIN patients p ON m.patient_id = p.id").
		Joins("JOIN users u ON m.user_id = u.id").
		Joins("LEFT JOIN localities l ON u.locality_id = l.id").
		Order("m.created_at DESC")

	// Aplicar filtros
	if filters != nil {
		if filters.LocalityID != nil {
			query = query.Where("u.locality_id = ?", *filters.LocalityID)
		}
		if filters.UserID != nil {
			query = query.Where("m.user_id = ?", *filters.UserID)
		}
		if filters.Days > 0 {
			since := time.Now().AddDate(0, 0, -filters.Days)
			query = query.Where("m.created_at >= ?", since)
		}
		if filters.Limit > 0 {
			query = query.Limit(filters.Limit)
		} else {
			query = query.Limit(50) // Límite por defecto
		}
	} else {
		query = query.Limit(50)
	}

	if err := query.Scan(&measurements).Error; err != nil {
		return nil, fmt.Errorf("error al obtener mediciones recientes: %w", err)
	}

	return &domain.RecentMeasurementsReport{
		Measurements: measurements,
	}, nil
}

// GetRiskPatients obtiene pacientes en riesgo
func (r *reportRepository) GetRiskPatients(ctx context.Context, filters *domain.ReportFilters) (*domain.RiskPatientsReport, error) {
	var patients []struct {
		PatientID    uuid.UUID
		PatientName  string
		Age          int
		Gender       string
		MuacValue    float64
		MuacCode     string
		LocalityName string
		UserName     string
		LastMeasure  time.Time
	}

	query := r.db.WithContext(ctx).
		Select(`
			p.id as patient_id,
			CONCAT(p.name, ' ', p.lastname) as patient_name,
			p.age,
			p.gender,
			m.muac_value,
			CASE 
				WHEN m.muac_value >= 11.5 AND m.muac_value < 12.5 THEN 'MUAC-Y1'
				WHEN m.muac_value < 11.5 THEN 'MUAC-R1'
			END as muac_code,
			l.name as locality_name,
			CONCAT(u.name, ' ', u.lastname) as user_name,
			m.created_at as last_measure
		`).
		Table("patients p").
		Joins(`JOIN measurements m ON p.id = m.patient_id AND m.id = (
			SELECT id FROM measurements m2 
			WHERE m2.patient_id = p.id 
			ORDER BY m2.created_at DESC 
			LIMIT 1
		)`).
		Joins("JOIN users u ON p.user_id = u.id").
		Joins("LEFT JOIN localities l ON u.locality_id = l.id").
		Where("m.muac_value < 12.5"). // Solo pacientes en riesgo
		Order("m.muac_value ASC")

	// Aplicar filtros
	if filters != nil {
		if filters.LocalityID != nil {
			query = query.Where("u.locality_id = ?", *filters.LocalityID)
		}
		if filters.UserID != nil {
			query = query.Where("p.user_id = ?", *filters.UserID)
		}
		if filters.Limit > 0 {
			query = query.Limit(filters.Limit)
		} else {
			query = query.Limit(100)
		}
	} else {
		query = query.Limit(100)
	}

	if err := query.Scan(&patients).Error; err != nil {
		return nil, fmt.Errorf("error al obtener pacientes en riesgo: %w", err)
	}

	// Separar por severidad
	var severeCases, moderateCases []domain.RiskPatient
	now := time.Now()

	for _, p := range patients {
		daysAgo := int(now.Sub(p.LastMeasure).Hours() / 24)

		riskPatient := domain.RiskPatient{
			PatientID:    p.PatientID,
			PatientName:  p.PatientName,
			Age:          p.Age,
			Gender:       p.Gender,
			MuacValue:    p.MuacValue,
			MuacCode:     p.MuacCode,
			LocalityName: p.LocalityName,
			UserName:     p.UserName,
			LastMeasure:  p.LastMeasure,
			DaysAgo:      daysAgo,
		}

		if p.MuacValue < domain.MuacThresholdSevere {
			severeCases = append(severeCases, riskPatient)
		} else {
			moderateCases = append(moderateCases, riskPatient)
		}
	}

	return &domain.RiskPatientsReport{
		SevereCases:   severeCases,
		ModerateCases: moderateCases,
	}, nil
}

// GetUserActivity obtiene la actividad de usuarios
func (r *reportRepository) GetUserActivity(ctx context.Context, filters *domain.ReportFilters) (*domain.UserActivityReport, error) {
	var users []domain.UserStats

	query := r.db.WithContext(ctx).
		Select(`
			u.id as user_id,
			CONCAT(u.name, ' ', u.lastname) as user_name,
			l.name as locality_name,
			COUNT(DISTINCT p.id) as total_patients,
			COUNT(m.id) as total_measures,
			MAX(m.created_at) as last_activity,
			COUNT(CASE WHEN m.created_at >= ? THEN 1 END) as measures_this_week
		`, time.Now().AddDate(0, 0, -7)).
		Table("users u").
		Joins("LEFT JOIN localities l ON u.locality_id = l.id").
		Joins("LEFT JOIN patients p ON u.id = p.user_id").
		Joins("LEFT JOIN measurements m ON u.id = m.user_id").
		Group("u.id, u.name, u.lastname, l.name").
		Order("total_measures DESC")

	// Aplicar filtros
	if filters != nil {
		if filters.LocalityID != nil {
			query = query.Where("u.locality_id = ?", *filters.LocalityID)
		}
		if filters.UserID != nil {
			query = query.Where("u.id = ?", *filters.UserID)
		}
		if filters.Days > 0 {
			since := time.Now().AddDate(0, 0, -filters.Days)
			query = query.Where("m.created_at >= ?", since)
		}
		if filters.Limit > 0 {
			query = query.Limit(filters.Limit)
		}
	}

	if err := query.Scan(&users).Error; err != nil {
		return nil, fmt.Errorf("error al obtener actividad de usuarios: %w", err)
	}

	return &domain.UserActivityReport{
		Users: users,
	}, nil
}

// Funciones helper

func (r *reportRepository) getStatusDistribution(ctx context.Context, filters *domain.ReportFilters) (*domain.StatusDistribution, error) {
	var result struct {
		Total    int64
		Normal   int64
		Moderate int64
		Severe   int64
	}

	query := r.db.WithContext(ctx).
		Select(`
			COUNT(*) as total,
			COUNT(CASE WHEN m.muac_value >= 12.5 THEN 1 END) as normal,
			COUNT(CASE WHEN m.muac_value >= 11.5 AND m.muac_value < 12.5 THEN 1 END) as moderate,
			COUNT(CASE WHEN m.muac_value < 11.5 THEN 1 END) as severe
		`).
		Table("patients p").
		Joins(`JOIN measurements m ON p.id = m.patient_id AND m.id = (
			SELECT id FROM measurements m2 
			WHERE m2.patient_id = p.id 
			ORDER BY m2.created_at DESC 
			LIMIT 1
		)`)

	if filters != nil {
		if filters.LocalityID != nil {
			query = query.Joins("JOIN users u ON p.user_id = u.id").
				Where("u.locality_id = ?", *filters.LocalityID)
		}
		if filters.Days > 0 {
			since := time.Now().AddDate(0, 0, -filters.Days)
			query = query.Where("m.created_at >= ?", since)
		}
	}

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}

	total := float64(result.Total)
	return &domain.StatusDistribution{
		Normal: domain.StatusCount{
			Total:      result.Normal,
			Percentage: r.calculatePercentage(int(result.Normal), total),
		},
		Moderate: domain.StatusCount{
			Total:      result.Moderate,
			Percentage: r.calculatePercentage(int(result.Moderate), total),
		},
		Severe: domain.StatusCount{
			Total:      result.Severe,
			Percentage: r.calculatePercentage(int(result.Severe), total),
		},
	}, nil
}

func (r *reportRepository) calculatePercentage(count int, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (float64(count) / total) * 100
}
