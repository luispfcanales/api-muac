// domain/report.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

// ============= ESTRUCTURAS DE REPORTES SIMPLES =============

// DashboardReport - Reporte principal del dashboard
type DashboardReport struct {
	TotalPatients      int64              `json:"total_patients"`
	TotalMeasurements  int64              `json:"total_measurements"`
	PatientsAtRisk     int64              `json:"patients_at_risk"`
	TotalUsers         int64              `json:"total_users"`
	StatusDistribution StatusDistribution `json:"status_distribution"`
	GeneratedAt        time.Time          `json:"generated_at"`
}

// StatusDistribution - Distribución por estado nutricional
type StatusDistribution struct {
	Normal   StatusCount `json:"normal"`   // Verde ≥ 12.5 cm
	Moderate StatusCount `json:"moderate"` // Amarillo 11.5-12.4 cm
	Severe   StatusCount `json:"severe"`   // Rojo < 11.5 cm
}

type StatusCount struct {
	Total      int64   `json:"total"`
	Percentage float64 `json:"percentage"`
}

// PatientsByLocalityReport - Pacientes agrupados por localidad
type PatientsByLocalityReport struct {
	LocalityData []LocalityData `json:"locality_data"`
	GeneratedAt  time.Time      `json:"generated_at"`
}

type LocalityData struct {
	LocalityID   uuid.UUID          `json:"locality_id"`
	LocalityName string             `json:"locality_name"`
	Total        int                `json:"total"`
	AtRisk       int                `json:"at_risk"`
	Distribution StatusDistribution `json:"distribution"`
}

// RecentMeasurementsReport - Mediciones recientes
type RecentMeasurementsReport struct {
	Measurements []RecentMeasurement `json:"measurements"`
	GeneratedAt  time.Time           `json:"generated_at"`
}

type RecentMeasurement struct {
	ID           uuid.UUID `json:"id"`
	PatientName  string    `json:"patient_name"`
	PatientAge   int       `json:"patient_age"`
	MuacValue    float64   `json:"muac_value"`
	MuacCode     string    `json:"muac_code"`
	ColorCode    string    `json:"color_code"`
	UserName     string    `json:"user_name"`
	LocalityName string    `json:"locality_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// RiskPatientsReport - Pacientes en riesgo
type RiskPatientsReport struct {
	SevereCases   []RiskPatient `json:"severe_cases"`
	ModerateCases []RiskPatient `json:"moderate_cases"`
	GeneratedAt   time.Time     `json:"generated_at"`
}

type RiskPatient struct {
	PatientID    uuid.UUID `json:"patient_id"`
	PatientName  string    `json:"patient_name"`
	Age          int       `json:"age"`
	Gender       string    `json:"gender"`
	MuacValue    float64   `json:"muac_value"`
	MuacCode     string    `json:"muac_code"`
	LocalityName string    `json:"locality_name"`
	UserName     string    `json:"user_name"`
	LastMeasure  time.Time `json:"last_measure"`
	DaysAgo      int       `json:"days_ago"`
}

// UserActivityReport - Actividad de usuarios
type UserActivityReport struct {
	Users       []UserStats `json:"users"`
	GeneratedAt time.Time   `json:"generated_at"`
}

type UserStats struct {
	UserID           uuid.UUID  `json:"user_id"`
	UserName         string     `json:"user_name"`
	LocalityName     string     `json:"locality_name"`
	TotalPatients    int        `json:"total_patients"`
	TotalMeasures    int        `json:"total_measures"`
	LastActivity     *time.Time `json:"last_activity"`
	MeasuresThisWeek int        `json:"measures_this_week"`
}

// ============= FILTROS SIMPLES =============
type ReportFilters struct {
	LocalityID *uuid.UUID `json:"locality_id,omitempty"`
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	Days       int        `json:"days,omitempty"`  // Últimos N días (default: 30)
	Limit      int        `json:"limit,omitempty"` // Límite de resultados (default: 100)
}
