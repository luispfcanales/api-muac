package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// PatientHandler maneja las peticiones HTTP relacionadas con pacientes
type PatientHandler struct {
	patientService ports.IPatientService
}

// NewPatientHandler crea una nueva instancia de PatientHandler
func NewPatientHandler(patientService ports.IPatientService) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
	}
}

// RegisterRoutes registra las rutas del manejador
func (h *PatientHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/patients", h.GetAllPatients)
	mux.HandleFunc("POST /api/patients", h.CreatePatient)
	mux.HandleFunc("GET /api/patients/{id}", h.GetPatientByID)
	mux.HandleFunc("PUT /api/patients/{id}", h.UpdatePatient)
	mux.HandleFunc("DELETE /api/patients/{id}", h.DeletePatient)
	mux.HandleFunc("GET /api/patients/dni/{dni}", h.GetPatientByDNI)
	mux.HandleFunc("GET /api/patients/father/{fatherId}", h.GetPatientsByFatherID)
	mux.HandleFunc("GET /api/patients/{id}/measurements", h.GetPatientMeasurements)
	mux.HandleFunc("POST /api/patients/{id}/measurements", h.AddPatientMeasurement)
}

// GetAllPatients obtiene todos los pacientes
func (h *PatientHandler) GetAllPatients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	patients, err := h.patientService.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patients)
}

// GetPatientByID obtiene un paciente por su ID
func (h *PatientHandler) GetPatientByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de paciente no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	patient, err := h.patientService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrPatientNotFound {
			http.Error(w, "Paciente no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
}

// GetPatientByDNI obtiene un paciente por su DNI
func (h *PatientHandler) GetPatientByDNI(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dni := r.PathValue("dni")
	if dni == "" {
		http.Error(w, "DNI no proporcionado", http.StatusBadRequest)
		return
	}

	patient, err := h.patientService.GetByDNI(ctx, dni)
	if err != nil {
		if err == domain.ErrPatientNotFound {
			http.Error(w, "Paciente no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
}

// CreatePatient crea un nuevo paciente
func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Name        string `json:"name"`
		Lastname    string `json:"lastname"`
		Gender      string `json:"gender"`
		Age         int    `json:"age"`
		BirthDate   string `json:"birth_date"`
		ArmSize     string `json:"arm_size"`
		Weight      string `json:"weight"`
		Size        string `json:"size"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	patient := domain.NewPatient(req.Name, req.Lastname)
	patient.Gender = req.Gender
	patient.Age = req.Age
	patient.BirthDate = req.BirthDate
	patient.ArmSize = req.ArmSize
	patient.Weight = req.Weight
	patient.Size = req.Size
	patient.Description = req.Description

	if err := h.patientService.Create(ctx, patient); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(patient)
}

// UpdatePatient actualiza un paciente existente
func (h *PatientHandler) UpdatePatient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de paciente no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		Name         string `json:"name"`
		Lastname     string `json:"lastname"`
		Gender       string `json:"gender"`
		Age          int    `json:"age"`
		BirthDate    string `json:"birth_date"`
		ArmSize      string `json:"arm_size"`
		Weight       string `json:"weight"`
		Size         string `json:"size"`
		ConsentGiven bool   `json:"consent_given"`
		Description  string `json:"description"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	patient, err := h.patientService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrPatientNotFound {
			http.Error(w, "Paciente no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	patient.Update(
		req.Name,
		req.Lastname,
		req.Gender,
		req.BirthDate,
		req.ArmSize,
		req.Weight,
		req.Size,
		req.Description,
		req.Age,
		req.ConsentGiven,
	)

	if err := h.patientService.Update(ctx, patient); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
}

// DeletePatient elimina un paciente por su ID
func (h *PatientHandler) DeletePatient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de paciente no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.patientService.Delete(ctx, id)
	if err != nil {
		if err == domain.ErrPatientNotFound {
			http.Error(w, "Paciente no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetPatientsByFatherID obtiene los pacientes asociados a un padre específico
func (h *PatientHandler) GetPatientsByFatherID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fatherIDStr := r.PathValue("fatherId")
	if fatherIDStr == "" {
		http.Error(w, "ID de padre no proporcionado", http.StatusBadRequest)
		return
	}

	fatherID, err := uuid.Parse(fatherIDStr)
	if err != nil {
		http.Error(w, "ID de padre inválido", http.StatusBadRequest)
		return
	}

	patients, err := h.patientService.GetByFatherID(ctx, fatherID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patients)
}

// GetPatientMeasurements obtiene las mediciones de un paciente específico
func (h *PatientHandler) GetPatientMeasurements(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de paciente no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	measurements, err := h.patientService.GetMeasurements(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurements)
}

// AddPatientMeasurement añade una nueva medición a un paciente
func (h *PatientHandler) AddPatientMeasurement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de paciente no proporcionado", http.StatusBadRequest)
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
		UserID           uuid.UUID `json:"user_id"`
		TagID            uuid.UUID `json:"tag_id"`
		RecommendationID uuid.UUID `json:"recommendation_id"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	measurement := domain.NewMeasurement(
		req.MuacValue,
		req.Description,
		req.Location,
		time.Now(),
		id,
		req.UserID,
		req.TagID,
		req.RecommendationID,
	)

	if err := h.patientService.AddMeasurement(ctx, id, measurement); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(measurement)
}
