package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// MeasurementHandler maneja las peticiones HTTP relacionadas con mediciones
type MeasurementHandler struct {
	measurementService ports.IMeasurementService
}

// NewMeasurementHandler crea una nueva instancia de MeasurementHandler
func NewMeasurementHandler(measurementService ports.IMeasurementService) *MeasurementHandler {
	return &MeasurementHandler{
		measurementService: measurementService,
	}
}

// RegisterRoutes registra las rutas del manejador
func (h *MeasurementHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/measurements", h.GetAllMeasurements)
	mux.HandleFunc("POST /api/measurements", h.CreateMeasurement)
	mux.HandleFunc("GET /api/measurements/{id}", h.GetMeasurementByID)
	mux.HandleFunc("PUT /api/measurements/{id}", h.UpdateMeasurement)
	mux.HandleFunc("DELETE /api/measurements/{id}", h.DeleteMeasurement)
	mux.HandleFunc("GET /api/measurements/patient/{patientId}", h.GetMeasurementsByPatientID)
	mux.HandleFunc("GET /api/measurements/user/{userId}", h.GetMeasurementsByUserID)
	mux.HandleFunc("GET /api/measurements/tag/{tagId}", h.GetMeasurementsByTagID)
	mux.HandleFunc("GET /api/measurements/recommendation/{recommendationId}", h.GetMeasurementsByRecommendationID)
	mux.HandleFunc("GET /api/measurements/date-range", h.GetMeasurementsByDateRange)
	mux.HandleFunc("PUT /api/measurements/{id}/tag/{tagId}", h.AssignTag)
	mux.HandleFunc("PUT /api/measurements/{id}/recommendation/{recommendationId}", h.AssignRecommendation)
}

// GetAllMeasurements obtiene todas las mediciones
func (h *MeasurementHandler) GetAllMeasurements(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	measurements, err := h.measurementService.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurements)
}

// GetMeasurementByID obtiene una medición por su ID
func (h *MeasurementHandler) GetMeasurementByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de medición no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	measurement, err := h.measurementService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrMeasurementNotFound {
			http.Error(w, "Medición no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurement)
}

// GetMeasurementsByPatientID obtiene mediciones por ID de paciente
func (h *MeasurementHandler) GetMeasurementsByPatientID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	patientIDStr := r.PathValue("patientId")
	if patientIDStr == "" {
		http.Error(w, "ID de paciente no proporcionado", http.StatusBadRequest)
		return
	}

	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		http.Error(w, "ID de paciente inválido", http.StatusBadRequest)
		return
	}

	measurements, err := h.measurementService.GetByPatientID(ctx, patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurements)
}

// GetMeasurementsByUserID obtiene mediciones por ID de usuario
func (h *MeasurementHandler) GetMeasurementsByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDStr := r.PathValue("userId")
	if userIDStr == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	measurements, err := h.measurementService.GetByUserID(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurements)
}

// GetMeasurementsByTagID obtiene mediciones por ID de etiqueta
func (h *MeasurementHandler) GetMeasurementsByTagID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tagIDStr := r.PathValue("tagId")
	if tagIDStr == "" {
		http.Error(w, "ID de etiqueta no proporcionado", http.StatusBadRequest)
		return
	}

	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		http.Error(w, "ID de etiqueta inválido", http.StatusBadRequest)
		return
	}

	measurements, err := h.measurementService.GetByTagID(ctx, tagID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurements)
}

// GetMeasurementsByRecommendationID obtiene mediciones por ID de recomendación
func (h *MeasurementHandler) GetMeasurementsByRecommendationID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	recommendationIDStr := r.PathValue("recommendationId")
	if recommendationIDStr == "" {
		http.Error(w, "ID de recomendación no proporcionado", http.StatusBadRequest)
		return
	}

	recommendationID, err := uuid.Parse(recommendationIDStr)
	if err != nil {
		http.Error(w, "ID de recomendación inválido", http.StatusBadRequest)
		return
	}

	measurements, err := h.measurementService.GetByRecommendationID(ctx, recommendationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurements)
}

// GetMeasurementsByDateRange obtiene mediciones dentro de un rango de fechas
func (h *MeasurementHandler) GetMeasurementsByDateRange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "Fechas de inicio y fin son requeridas", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		http.Error(w, "Formato de fecha de inicio inválido. Use RFC3339", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		http.Error(w, "Formato de fecha de fin inválido. Use RFC3339", http.StatusBadRequest)
		return
	}

	measurements, err := h.measurementService.GetByDateRange(ctx, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurements)
}

// CreateMeasurement crea una nueva medición
func (h *MeasurementHandler) CreateMeasurement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		MuacValue        float64   `json:"muac_value"`
		Description      string    `json:"description"`
		Location         string    `json:"location"`
		Timestamp        time.Time `json:"timestamp"`
		PatientID        uuid.UUID `json:"patient_id"`
		UserID           uuid.UUID `json:"user_id"`
		TagID            uuid.UUID `json:"tag_id"`
		RecommendationID uuid.UUID `json:"recommendation_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	// Si no se proporciona una marca de tiempo, usar la hora actual
	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}

	measurement := domain.NewMeasurement(
		req.MuacValue,
		req.Description,
		req.Location,
		req.Timestamp,
		req.PatientID,
		req.UserID,
		req.TagID,
		req.RecommendationID,
	)

	if err := h.measurementService.Create(ctx, measurement); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(measurement)
}

// UpdateMeasurement actualiza una medición existente
func (h *MeasurementHandler) UpdateMeasurement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de medición no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		MuacValue        float64   `json:"muac_value"`
		Description      string    `json:"description"`
		Location         string    `json:"location"`
		Timestamp        time.Time `json:"timestamp"`
		TagID            uuid.UUID `json:"tag_id"`
		RecommendationID uuid.UUID `json:"recommendation_id"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	measurement, err := h.measurementService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrMeasurementNotFound {
			http.Error(w, "Medición no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Si no se proporciona una marca de tiempo, mantener la existente
	if req.Timestamp.IsZero() {
		req.Timestamp = measurement.Timestamp
	}

	measurement.Update(
		req.MuacValue,
		req.Description,
		req.Location,
		req.Timestamp,
		req.TagID,
		req.RecommendationID,
	)

	if err := h.measurementService.Update(ctx, measurement); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurement)
}

// DeleteMeasurement elimina una medición por su ID
func (h *MeasurementHandler) DeleteMeasurement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de medición no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.measurementService.Delete(ctx, id)
	if err != nil {
		if err == domain.ErrMeasurementNotFound {
			http.Error(w, "Medición no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AssignTag asigna una etiqueta a una medición
func (h *MeasurementHandler) AssignTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de medición no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de medición inválido", http.StatusBadRequest)
		return
	}

	tagIDStr := r.PathValue("tagId")
	if tagIDStr == "" {
		http.Error(w, "ID de etiqueta no proporcionado", http.StatusBadRequest)
		return
	}

	var tagID uuid.UUID
	if tagIDStr == "null" {
		tagID = uuid.Nil
	} else {
		tagID, err = uuid.Parse(tagIDStr)
		if err != nil {
			http.Error(w, "ID de etiqueta inválido", http.StatusBadRequest)
			return
		}
	}

	err = h.measurementService.AssignTag(ctx, id, tagID)
	if err != nil {
		if err == domain.ErrMeasurementNotFound {
			http.Error(w, "Medición no encontrada", http.StatusNotFound)
			return
		}
		if err == domain.ErrTagNotFound {
			http.Error(w, "Etiqueta no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AssignRecommendation asigna una recomendación a una medición
func (h *MeasurementHandler) AssignRecommendation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de medición no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de medición inválido", http.StatusBadRequest)
		return
	}

	recommendationIDStr := r.PathValue("recommendationId")
	if recommendationIDStr == "" {
		http.Error(w, "ID de recomendación no proporcionado", http.StatusBadRequest)
		return
	}

	var recommendationID uuid.UUID
	if recommendationIDStr == "null" {
		recommendationID = uuid.Nil
	} else {
		recommendationID, err = uuid.Parse(recommendationIDStr)
		if err != nil {
			http.Error(w, "ID de recomendación inválido", http.StatusBadRequest)
			return
		}
	}

	err = h.measurementService.AssignRecommendation(ctx, id, recommendationID)
	if err != nil {
		if err == domain.ErrMeasurementNotFound {
			http.Error(w, "Medición no encontrada", http.StatusNotFound)
			return
		}
		if err == domain.ErrRecommendationNotFound {
			http.Error(w, "Recomendación no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}