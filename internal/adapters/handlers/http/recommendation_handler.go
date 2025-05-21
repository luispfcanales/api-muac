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

// GetAllRecommendations obtiene todas las recomendaciones
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

// GetRecommendationByID obtiene una recomendación por su ID
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

// GetRecommendationByName obtiene una recomendación por su nombre
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

// GetRecommendationsByUmbral obtiene recomendaciones por su umbral
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

// CreateRecommendation crea una nueva recomendación
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

// UpdateRecommendation actualiza una recomendación existente
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

// DeleteRecommendation elimina una recomendación por su ID
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