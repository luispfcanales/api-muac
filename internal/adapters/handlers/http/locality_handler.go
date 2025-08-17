package http

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	mux.HandleFunc("GET /api/localities/nearby", h.GetNearbyLocalities)
}

// GetAllLocalities godoc
// @Summary Obtener todas las localidades
// @Description Obtiene una lista de todas las localidades registradas en el sistema
// @Tags localidades
// @Accept json
// @Produce json
// @Success 200 {array} domain.Locality
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/localities [get]
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

// CreateLocality godoc
// @Summary Crear una nueva localidad
// @Description Crea una nueva localidad con la información proporcionada
// @Tags localidades
// @Accept json
// @Produce json
// @Param locality body object true "Datos de la localidad"
// @Success 201 {object} domain.Locality
// @Failure 400 {object} map[string]string "Solicitud inválida"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/localities [post]
func (h *LocalityHandler) CreateLocality(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Name        string `json:"name"`
		Latitude    string `json:"latitude"`
		Longitude   string `json:"longitude"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	locality := domain.NewLocality(req.Name, req.Latitude, req.Longitude, req.Description)

	if err := h.localityService.Create(ctx, locality); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(locality)
}

// GetLocalityByID godoc
// @Summary Obtener una localidad por ID
// @Description Obtiene una localidad específica por su ID
// @Tags localidades
// @Accept json
// @Produce json
// @Param id path string true "ID de la localidad"
// @Success 200 {object} domain.Locality
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Localidad no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/localities/{id} [get]
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

// UpdateLocality godoc
// @Summary Actualizar una localidad
// @Description Actualiza una localidad existente con la información proporcionada
// @Tags localidades
// @Accept json
// @Produce json
// @Param id path string true "ID de la localidad"
// @Param locality body object true "Datos actualizados de la localidad"
// @Success 200 {object} domain.Locality
// @Failure 400 {object} map[string]string "ID inválido o solicitud inválida"
// @Failure 404 {object} map[string]string "Localidad no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/localities/{id} [put]
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

// DeleteLocality godoc
// @Summary Eliminar una localidad
// @Description Elimina una localidad por su ID
// @Tags localidades
// @Accept json
// @Produce json
// @Param id path string true "ID de la localidad"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Localidad no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/localities/{id} [delete]
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

// GetLocalityByName godoc
// @Summary Obtener una localidad por nombre
// @Description Obtiene una localidad específica por su nombre
// @Tags localidades
// @Accept json
// @Produce json
// @Param name path string true "Nombre de la localidad"
// @Success 200 {object} domain.Locality
// @Failure 400 {object} map[string]string "Nombre no proporcionado"
// @Failure 404 {object} map[string]string "Localidad no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/localities/name/{name} [get]
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

func (h *LocalityHandler) GetNearbyLocalities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parsear parámetros
	var request struct {
		Latitude  string  `json:"latitude"`
		Longitude string  `json:"longitude"`
		RadiusKm  float64 `json:"radius_km"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validar parámetros
	if request.RadiusKm <= 0 {
		request.RadiusKm = 10 // Radio por defecto de 10km
	}
	//parseamos las coordenadas
	lat, err := strconv.ParseFloat(request.Latitude, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}
	lng, err := strconv.ParseFloat(request.Longitude, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	// Obtener localidades cercanas
	localities, err := h.localityService.FindNearbyLocalities(
		ctx,
		lat,
		lng,
		request.RadiusKm,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if localities == nil {
		localities = []domain.Locality{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(localities)
}
