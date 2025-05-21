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

// GetAllFathers godoc
// @Summary Obtener todos los padres
// @Description Obtiene una lista de todos los padres registrados en el sistema
// @Tags padres
// @Accept json
// @Produce json
// @Success 200 {array} domain.Father
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/fathers [get]
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

// GetFatherByID godoc
// @Summary Obtener un padre por ID
// @Description Obtiene un padre específico por su ID
// @Tags padres
// @Accept json
// @Produce json
// @Param id path string true "ID del padre"
// @Success 200 {object} domain.Father
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Padre no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/fathers/{id} [get]
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

// GetFatherByEmail godoc
// @Summary Obtener un padre por email
// @Description Obtiene un padre específico por su dirección de email
// @Tags padres
// @Accept json
// @Produce json
// @Param email path string true "Email del padre"
// @Success 200 {object} domain.Father
// @Failure 400 {object} map[string]string "Email no proporcionado"
// @Failure 404 {object} map[string]string "Padre no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/fathers/email/{email} [get]
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

// GetFatherByDNI godoc
// @Summary Obtener un padre por DNI
// @Description Obtiene un padre específico por su número de DNI
// @Tags padres
// @Accept json
// @Produce json
// @Param dni path string true "DNI del padre"
// @Success 200 {object} domain.Father
// @Failure 400 {object} map[string]string "DNI no proporcionado o inválido"
// @Failure 404 {object} map[string]string "Padre no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/fathers/dni/{dni} [get]
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

// GetFathersByPatientID godoc
// @Summary Obtener padres por ID del paciente
// @Description Obtiene todos los padres asociados a un paciente específico
// @Tags padres
// @Accept json
// @Produce json
// @Param patientId path string true "ID del paciente"
// @Success 200 {array} domain.Father
// @Failure 400 {object} map[string]string "ID de paciente inválido o no proporcionado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/fathers/patient/{patientId} [get]
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

// GetFathersByLocalityID godoc
// @Summary Obtener padres por ID de localidad
// @Description Obtiene todos los padres asociados a una localidad específica
// @Tags padres
// @Accept json
// @Produce json
// @Param localityId path string true "ID de la localidad"
// @Success 200 {array} domain.Father
// @Failure 400 {object} map[string]string "ID de localidad inválido o no proporcionado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/fathers/locality/{localityId} [get]
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