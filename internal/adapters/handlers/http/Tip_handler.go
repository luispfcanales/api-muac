package http

import (
	"encoding/json"
	"net/http"

	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// TipRecipesHandler maneja las solicitudes HTTP relacionadas con las recetas de tips
type TipHandler struct {
	TipRecipeService ports.ITipService
	RecipeService    ports.IRecipeService
}

// NewTipRecipesHandler crea una nueva instancia de TipRecipesHandler
func NewTipHandler(
	tipSrv ports.ITipService,
	recipeSrv ports.IRecipeService,
) *TipHandler {
	return &TipHandler{
		TipRecipeService: tipSrv,
		RecipeService:    recipeSrv,
	}
}

// RegisterRoutes registra las rutas HTTP para las recetas de tips
func (h *TipHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/tip-recipes", h.GetAllTipRecipes)
}

// GetAllTipRecipes godoc
// @Summary Obtener todas las recetas de tips
// @Description Obtiene una lista de todas las recetas de tips registradas
// @Tags tip-recipes
// @Accept json
// @Produce json
// @Success 200 {array} domain.Tip
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/tip-recipes [get]
func (h *TipHandler) GetAllTipRecipes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request struct {
		MUACCode string  `json:"muac_code"`
		Age      float64 `json:"age"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tips, err := h.TipRecipeService.List(ctx, request.MUACCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	recipes, err := h.RecipeService.ListRecipesByAge(ctx, request.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tips":    tips,
		"recipes": recipes,
	})
}
