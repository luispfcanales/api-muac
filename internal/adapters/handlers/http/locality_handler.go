package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// LocalityHandler maneja las peticiones HTTP relacionadas con localidades
type LocalityHandler struct {
	localityService ports.ILocalityService
}

// NewLocalityHandler crea una nueva instancia de LocalityHandler
func NewLocalityHandler(localityService ports.ILocalityService) *LocalityHandler {
	return &LocalityHandler{
		localityService: localityService,
	}
}

// RegisterRoutes registra las rutas del manejador
func (h *LocalityHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/localities", h.GetAllLocalities)
	mux.HandleFunc("POST /api/localities", h.CreateLocality)
	mux.HandleFunc("GET /api/localities/{id}", h.GetLocalityByID)
	mux.HandleFunc("PUT /api/localities/{id}", h.UpdateLocality)
	mux.HandleFunc("DELETE /api/localities/{id}", h.DeleteLocality)
	mux.HandleFunc("GET /api/localities/name/{name}", h.GetLocalityByName)
}

// GetAllLocalities obtiene todas las localidades
func (h *LocalityHandler) GetAllLocalities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	localities, err := h.localityService.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(localities)
}

// GetLocalityByID obtiene una localidad por su ID
func (h *LocalityHandler) GetLocalityByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de localidad no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	locality, err := h.localityService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrLocalityNotFound {
			http.Error(w, "Localidad no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locality)
}

// GetLocalityByName obtiene una localidad por su nombre
func (h *LocalityHandler) GetLocalityByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "Nombre de localidad no proporcionado", http.StatusBadRequest)
		return
	}

	locality, err := h.localityService.GetByName(ctx, name)
	if err != nil {
		if err == domain.ErrLocalityNotFound {
			http.Error(w, "Localidad no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locality)
}

// CreateLocality crea una nueva localidad
func (h *LocalityHandler) CreateLocality(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Name        string `json:"name"`
		Location    string `json:"location"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	locality := domain.NewLocality(req.Name, req.Location, req.Description)

	if err := h.localityService.Create(ctx, locality); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(locality)
}

// UpdateLocality actualiza una localidad existente
func (h *LocalityHandler) UpdateLocality(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de localidad no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Location    string `json:"location"`
		Description string `json:"description"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	locality, err := h.localityService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrLocalityNotFound {
			http.Error(w, "Localidad no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	locality.Update(req.Name, req.Location, req.Description)

	if err := h.localityService.Update(ctx, locality); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locality)
}

// DeleteLocality elimina una localidad por su ID
func (h *LocalityHandler) DeleteLocality(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de localidad no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.localityService.Delete(ctx, id)
	if err != nil {
		if err == domain.ErrLocalityNotFound {
			http.Error(w, "Localidad no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}