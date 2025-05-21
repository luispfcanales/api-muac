package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// FatherHandler maneja las peticiones HTTP relacionadas con padres
type FatherHandler struct {
	fatherService ports.IFatherService
}

// NewFatherHandler crea una nueva instancia de FatherHandler
func NewFatherHandler(fatherService ports.IFatherService) *FatherHandler {
	return &FatherHandler{
		fatherService: fatherService,
	}
}

// RegisterRoutes registra las rutas del manejador
func (h *FatherHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/fathers", h.GetAllFathers)
	mux.HandleFunc("POST /api/fathers", h.CreateFather)
	mux.HandleFunc("GET /api/fathers/{id}", h.GetFatherByID)
	mux.HandleFunc("PUT /api/fathers/{id}", h.UpdateFather)
	mux.HandleFunc("DELETE /api/fathers/{id}", h.DeleteFather)
	mux.HandleFunc("GET /api/fathers/email/{email}", h.GetFatherByEmail)
	mux.HandleFunc("GET /api/fathers/dni/{dni}", h.GetFatherByDNI)
	mux.HandleFunc("GET /api/fathers/patient/{patientId}", h.GetFathersByPatientID)
	mux.HandleFunc("GET /api/fathers/locality/{localityId}", h.GetFathersByLocalityID)
	mux.HandleFunc("PUT /api/fathers/{id}/password", h.UpdatePassword)
	mux.HandleFunc("PUT /api/fathers/{id}/active", h.UpdateActive)
}

// GetAllFathers obtiene todos los padres
func (h *FatherHandler) GetAllFathers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fathers, err := h.fatherService.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fathers)
}

// GetFatherByID obtiene un padre por su ID
func (h *FatherHandler) GetFatherByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de padre no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	father, err := h.fatherService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrFatherNotFound {
			http.Error(w, "Padre no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(father)
}

// GetFatherByEmail obtiene un padre por su email
func (h *FatherHandler) GetFatherByEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	email := r.PathValue("email")
	if email == "" {
		http.Error(w, "Email no proporcionado", http.StatusBadRequest)
		return
	}

	father, err := h.fatherService.GetByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrFatherNotFound {
			http.Error(w, "Padre no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(father)
}

// GetFatherByDNI obtiene un padre por su DNI
func (h *FatherHandler) GetFatherByDNI(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dniStr := r.PathValue("dni")
	if dniStr == "" {
		http.Error(w, "DNI no proporcionado", http.StatusBadRequest)
		return
	}

	dni, err := strconv.Atoi(dniStr)
	if err != nil {
		http.Error(w, "DNI inválido", http.StatusBadRequest)
		return
	}

	father, err := h.fatherService.GetByDNI(ctx, dni)
	if err != nil {
		if err == domain.ErrFatherNotFound {
			http.Error(w, "Padre no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(father)
}

// GetFathersByPatientID obtiene padres por ID de paciente
func (h *FatherHandler) GetFathersByPatientID(w http.ResponseWriter, r *http.Request) {
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

	fathers, err := h.fatherService.GetByPatientID(ctx, patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fathers)
}

// GetFathersByLocalityID obtiene padres por ID de localidad
func (h *FatherHandler) GetFathersByLocalityID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	localityIDStr := r.PathValue("localityId")
	if localityIDStr == "" {
		http.Error(w, "ID de localidad no proporcionado", http.StatusBadRequest)
		return
	}

	localityID, err := uuid.Parse(localityIDStr)
	if err != nil {
		http.Error(w, "ID de localidad inválido", http.StatusBadRequest)
		return
	}

	fathers, err := h.fatherService.GetByLocalityID(ctx, localityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fathers)
}

// CreateFather crea un nuevo padre
func (h *FatherHandler) CreateFather(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Name         string    `json:"name"`
		LastName     string    `json:"lastname"`
		Email        string    `json:"email"`
		DNI          int       `json:"dni"`
		Phone        string    `json:"phone"`
		PasswordHash string    `json:"password_hash"`
		Active       bool      `json:"active"`
		RoleID       uuid.UUID `json:"role_id"`
		LocalityID   uuid.UUID `json:"locality_id"`
		PatientID    uuid.UUID `json:"patient_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	father := domain.NewFather(
		req.Name,
		req.LastName,
		req.Email,
		req.PasswordHash,
		req.RoleID,
		req.LocalityID,
		req.PatientID,
	)
	
	father.DNI = req.DNI
	father.Phone = req.Phone
	father.Active = req.Active

	if err := h.fatherService.Create(ctx, father); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(father)
}

// UpdateFather actualiza un padre existente
func (h *FatherHandler) UpdateFather(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de padre no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		Name       string    `json:"name"`
		LastName   string    `json:"lastname"`
		Email      string    `json:"email"`
		DNI        int       `json:"dni"`
		Phone      string    `json:"phone"`
		Active     bool      `json:"active"`
		RoleID     uuid.UUID `json:"role_id"`
		LocalityID uuid.UUID `json:"locality_id"`
		PatientID  uuid.UUID `json:"patient_id"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	father, err := h.fatherService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrFatherNotFound {
			http.Error(w, "Padre no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	father.Update(
		req.Name,
		req.LastName,
		req.Email,
		req.Phone,
		req.DNI,
		req.Active,
		req.RoleID,
		req.LocalityID,
		req.PatientID,
	)

	if err := h.fatherService.Update(ctx, father); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(father)
}

// DeleteFather elimina un padre por su ID
func (h *FatherHandler) DeleteFather(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de padre no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.fatherService.Delete(ctx, id)
	if err != nil {
		if err == domain.ErrFatherNotFound {
			http.Error(w, "Padre no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdatePassword actualiza la contraseña de un padre
func (h *FatherHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de padre no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		PasswordHash string `json:"password_hash"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	if req.PasswordHash == "" {
		http.Error(w, "Contraseña no proporcionada", http.StatusBadRequest)
		return
	}

	err = h.fatherService.UpdatePassword(ctx, id, req.PasswordHash)
	if err != nil {
		if err == domain.ErrFatherNotFound {
			http.Error(w, "Padre no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateActive actualiza el estado activo de un padre
func (h *FatherHandler) UpdateActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de padre no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		Active bool `json:"active"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	err = h.fatherService.UpdateActive(ctx, id, req.Active)
	if err != nil {
		if err == domain.ErrFatherNotFound {
			http.Error(w, "Padre no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}