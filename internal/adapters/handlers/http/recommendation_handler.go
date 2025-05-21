package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// RecommendationHandler maneja las peticiones HTTP relacionadas con recomendaciones
type RecommendationHandler struct {
	recommendationService ports.IRecommendationService
}

// NewRecommendationHandler crea una nueva instancia de RecommendationHandler
func NewRecommendationHandler(recommendationService ports.IRecommendationService) *RecommendationHandler {
	return &RecommendationHandler{
		recommendationService: recommendationService,
	}
}

// RegisterRoutes registra las rutas del manejador
func (h *RecommendationHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/recommendations", h.GetAllRecommendations)
	mux.HandleFunc("POST /api/recommendations", h.CreateRecommendation)
	mux.HandleFunc("GET /api/recommendations/{id}", h.GetRecommendationByID)
	mux.HandleFunc("PUT /api/recommendations/{id}", h.UpdateRecommendation)
	mux.HandleFunc("DELETE /api/recommendations/{id}", h.DeleteRecommendation)
	mux.HandleFunc("GET /api/recommendations/name/{name}", h.GetRecommendationByName)
	mux.HandleFunc("GET /api/recommendations/umbral/{umbral}", h.GetRecommendationsByUmbral)
}

// GetAllRecommendations godoc
// @Summary Obtener todas las recomendaciones
// @Description Obtiene una lista de todas las recomendaciones registradas en el sistema
// @Tags recomendaciones
// @Accept json
// @Produce json
// @Success 200 {array} domain.Recommendation
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/recommendations [get]
func (h *RecommendationHandler) GetAllRecommendations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	recommendations, err := h.recommendationService.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}

// CreateRecommendation godoc
// @Summary Crear una nueva recomendación
// @Description Crea una nueva recomendación con la información proporcionada
// @Tags recomendaciones
// @Accept json
// @Produce json
// @Param recommendation body object true "Datos de la recomendación"
// @Success 201 {object} domain.Recommendation
// @Failure 400 {object} map[string]string "Solicitud inválida"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/recommendations [post]
func (h *RecommendationHandler) CreateRecommendation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Umbral      string `json:"recommendation_umbral"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	recommendation := domain.NewRecommendation(req.Name, req.Description, req.Umbral)

	if err := h.recommendationService.Create(ctx, recommendation); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(recommendation)
}

// GetRecommendationByID godoc
// @Summary Obtener una recomendación por ID
// @Description Obtiene una recomendación específica por su ID
// @Tags recomendaciones
// @Accept json
// @Produce json
// @Param id path string true "ID de la recomendación"
// @Success 200 {object} domain.Recommendation
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Recomendación no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/recommendations/{id} [get]
func (h *RecommendationHandler) GetRecommendationByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de recomendación no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	recommendation, err := h.recommendationService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrRecommendationNotFound {
			http.Error(w, "Recomendación no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendation)
}

// UpdateRecommendation godoc
// @Summary Actualizar una recomendación
// @Description Actualiza una recomendación existente con la información proporcionada
// @Tags recomendaciones
// @Accept json
// @Produce json
// @Param id path string true "ID de la recomendación"
// @Param recommendation body object true "Datos actualizados de la recomendación"
// @Success 200 {object} domain.Recommendation
// @Failure 400 {object} map[string]string "ID inválido o solicitud inválida"
// @Failure 404 {object} map[string]string "Recomendación no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/recommendations/{id} [put]
func (h *RecommendationHandler) UpdateRecommendation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de recomendación no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Umbral      string `json:"recommendation_umbral"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	recommendation, err := h.recommendationService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrRecommendationNotFound {
			http.Error(w, "Recomendación no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	recommendation.Update(req.Name, req.Description, req.Umbral)

	if err := h.recommendationService.Update(ctx, recommendation); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendation)
}

// DeleteRecommendation godoc
// @Summary Eliminar una recomendación
// @Description Elimina una recomendación por su ID
// @Tags recomendaciones
// @Accept json
// @Produce json
// @Param id path string true "ID de la recomendación"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Recomendación no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/recommendations/{id} [delete]
func (h *RecommendationHandler) DeleteRecommendation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de recomendación no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.recommendationService.Delete(ctx, id)
	if err != nil {
		if err == domain.ErrRecommendationNotFound {
			http.Error(w, "Recomendación no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetRecommendationByName godoc
// @Summary Obtener una recomendación por nombre
// @Description Obtiene una recomendación específica por su nombre
// @Tags recomendaciones
// @Accept json
// @Produce json
// @Param name path string true "Nombre de la recomendación"
// @Success 200 {object} domain.Recommendation
// @Failure 400 {object} map[string]string "Nombre no proporcionado"
// @Failure 404 {object} map[string]string "Recomendación no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/recommendations/name/{name} [get]
func (h *RecommendationHandler) GetRecommendationByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "Nombre de recomendación no proporcionado", http.StatusBadRequest)
		return
	}

	recommendation, err := h.recommendationService.GetByName(ctx, name)
	if err != nil {
		if err == domain.ErrRecommendationNotFound {
			http.Error(w, "Recomendación no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendation)
}

// GetRecommendationsByUmbral godoc
// @Summary Obtener recomendaciones por umbral
// @Description Obtiene todas las recomendaciones que coinciden con un umbral específico
// @Tags recomendaciones
// @Accept json
// @Produce json
// @Param umbral path string true "Umbral de la recomendación"
// @Success 200 {array} domain.Recommendation
// @Failure 400 {object} map[string]string "Umbral no proporcionado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/recommendations/umbral/{umbral} [get]
func (h *RecommendationHandler) GetRecommendationsByUmbral(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	umbral := r.PathValue("umbral")
	if umbral == "" {
		http.Error(w, "Umbral de recomendación no proporcionado", http.StatusBadRequest)
		return
	}

	recommendations, err := h.recommendationService.GetByUmbral(ctx, umbral)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}