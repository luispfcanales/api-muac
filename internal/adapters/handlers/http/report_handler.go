// http/report_handler.go
package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// ReportHandler maneja las peticiones HTTP relacionadas con reportes
type ReportHandler struct {
	reportService ports.IReportService
}

// NewReportHandler crea una nueva instancia de ReportHandler
func NewReportHandler(reportService ports.IReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// RegisterRoutes registra las rutas del manejador
func (h *ReportHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/reports/dashboard", h.GetDashboard)
	mux.HandleFunc("GET /api/reports/patients-by-locality", h.GetPatientsByLocality)
	mux.HandleFunc("GET /api/reports/recent-measurements", h.GetRecentMeasurements)
	mux.HandleFunc("GET /api/reports/risk-patients", h.GetRiskPatients)
	mux.HandleFunc("GET /api/reports/user-activity", h.GetUserActivity)
	mux.HandleFunc("GET /api/reports/risk-patients-coordinates", h.GetRiskPatientsCoordinates)
}

// GetDashboard godoc
// @Summary Obtener datos del dashboard principal
// @Description Obtiene las estadísticas principales del dashboard (total pacientes, mediciones, etc.)
// @Tags reports
// @Accept json
// @Produce json
// @Param locality_id query string false "ID de la localidad para filtrar"
// @Param days query int false "Número de días hacia atrás (default: 30)"
// @Success 200 {object} domain.DashboardReport
// @Failure 400 {object} map[string]string "Parámetros inválidos"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/reports/dashboard [get]
func (h *ReportHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filters, err := h.parseFilters(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	report, err := h.reportService.GetDashboardReport(ctx, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GetPatientsByLocality godoc
// @Summary Obtener pacientes agrupados por localidad
// @Description Obtiene estadísticas de pacientes organizadas por localidad
// @Tags reports
// @Accept json
// @Produce json
// @Param locality_id query string false "ID de la localidad para filtrar"
// @Param days query int false "Número de días hacia atrás (default: 30)"
// @Param limit query int false "Límite de resultados (default: 100)"
// @Success 200 {object} domain.PatientsByLocalityReport
// @Failure 400 {object} map[string]string "Parámetros inválidos"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/reports/patients-by-locality [get]
func (h *ReportHandler) GetPatientsByLocality(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filters, err := h.parseFilters(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	report, err := h.reportService.GetPatientsByLocalityReport(ctx, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GetRecentMeasurements godoc
// @Summary Obtener mediciones recientes
// @Description Obtiene las mediciones más recientes registradas en el sistema
// @Tags reports
// @Accept json
// @Produce json
// @Param locality_id query string false "ID de la localidad para filtrar"
// @Param user_id query string false "ID del usuario para filtrar"
// @Param days query int false "Número de días hacia atrás (default: 7)"
// @Param limit query int false "Límite de resultados (default: 50)"
// @Success 200 {object} domain.RecentMeasurementsReport
// @Failure 400 {object} map[string]string "Parámetros inválidos"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/reports/recent-measurements [get]
func (h *ReportHandler) GetRecentMeasurements(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filters, err := h.parseFilters(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Para mediciones recientes, usar por defecto 7 días si no se especifica
	if filters.Days == 0 {
		filters.Days = 7
	}

	// Límite por defecto para mediciones recientes
	if filters.Limit == 0 {
		filters.Limit = 50
	}

	report, err := h.reportService.GetRecentMeasurementsReport(ctx, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GetRiskPatients godoc
// @Summary Obtener pacientes en riesgo
// @Description Obtiene la lista de pacientes en riesgo nutricional (casos moderados y severos)
// @Tags reports
// @Accept json
// @Produce json
// @Param locality_id query string false "ID de la localidad para filtrar"
// @Param user_id query string false "ID del usuario para filtrar"
// @Param limit query int false "Límite de resultados (default: 100)"
// @Success 200 {object} domain.RiskPatientsReport
// @Failure 400 {object} map[string]string "Parámetros inválidos"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/reports/risk-patients [get]
func (h *ReportHandler) GetRiskPatients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filters, err := h.parseFilters(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Límite por defecto para pacientes en riesgo
	if filters.Limit == 0 {
		filters.Limit = 100
	}

	report, err := h.reportService.GetRiskPatientsReport(ctx, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GetRiskPatientsCoordinates obtiene coordenadas para mapa de calor
func (h *ReportHandler) GetRiskPatientsCoordinates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filters, err := h.parseFilters(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	coordinates, err := h.reportService.GetRiskPatientsCoordinates(ctx, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respuesta simple: solo el array de coordenadas
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coordinates)
}

// GetUserActivity godoc
// @Summary Obtener actividad de usuarios
// @Description Obtiene estadísticas de actividad de los usuarios del sistema
// @Tags reports
// @Accept json
// @Produce json
// @Param locality_id query string false "ID de la localidad para filtrar"
// @Param user_id query string false "ID del usuario para filtrar"
// @Param days query int false "Número de días hacia atrás (default: 30)"
// @Param limit query int false "Límite de resultados (default: 50)"
// @Success 200 {object} domain.UserActivityReport
// @Failure 400 {object} map[string]string "Parámetros inválidos"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/reports/user-activity [get]
func (h *ReportHandler) GetUserActivity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filters, err := h.parseFilters(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Límite por defecto para actividad de usuarios
	if filters.Limit == 0 {
		filters.Limit = 50
	}

	report, err := h.reportService.GetUserActivityReport(ctx, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// parseFilters parsea los query parameters a filtros
func (h *ReportHandler) parseFilters(r *http.Request) (*domain.ReportFilters, error) {
	filters := &domain.ReportFilters{}

	// Locality ID
	if localityIDStr := r.URL.Query().Get("locality_id"); localityIDStr != "" {
		localityID, err := uuid.Parse(localityIDStr)
		if err != nil {
			return nil, fmt.Errorf("locality_id inválido: %v", err)
		}
		filters.LocalityID = &localityID
	}

	// User ID
	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("user_id inválido: %v", err)
		}
		filters.UserID = &userID
	}

	// Days
	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return nil, fmt.Errorf("days debe ser un número válido: %v", err)
		}
		if days < 0 {
			return nil, fmt.Errorf("days no puede ser negativo")
		}
		if days > 365 {
			return nil, fmt.Errorf("days no puede ser mayor a 365")
		}
		filters.Days = days
	} else {
		filters.Days = 30 // Por defecto últimos 30 días
	}

	// Limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("limit debe ser un número válido: %v", err)
		}
		if limit < 0 {
			return nil, fmt.Errorf("limit no puede ser negativo")
		}
		if limit > 1000 {
			return nil, fmt.Errorf("limit no puede ser mayor a 1000")
		}
		filters.Limit = limit
	}

	return filters, nil
}
