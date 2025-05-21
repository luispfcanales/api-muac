package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// TagHandler maneja las peticiones HTTP relacionadas con etiquetas
type TagHandler struct {
	tagService ports.ITagService
}

// NewTagHandler crea una nueva instancia de TagHandler
func NewTagHandler(tagService ports.ITagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// RegisterRoutes registra las rutas del manejador
func (h *TagHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/tags", h.GetAllTags)
	mux.HandleFunc("POST /api/tags", h.CreateTag)
	mux.HandleFunc("GET /api/tags/{id}", h.GetTagByID)
	mux.HandleFunc("PUT /api/tags/{id}", h.UpdateTag)
	mux.HandleFunc("DELETE /api/tags/{id}", h.DeleteTag)
	mux.HandleFunc("GET /api/tags/name/{name}", h.GetTagByName)
}

// GetAllTags obtiene todas las etiquetas
func (h *TagHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tags, err := h.tagService.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

// GetTagByID obtiene una etiqueta por su ID
func (h *TagHandler) GetTagByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de etiqueta no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	tag, err := h.tagService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrTagNotFound {
			http.Error(w, "Etiqueta no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tag)
}

// GetTagByName obtiene una etiqueta por su nombre
func (h *TagHandler) GetTagByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "Nombre de etiqueta no proporcionado", http.StatusBadRequest)
		return
	}

	tag, err := h.tagService.GetByName(ctx, name)
	if err != nil {
		if err == domain.ErrTagNotFound {
			http.Error(w, "Etiqueta no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tag)
}

// CreateTag crea una nueva etiqueta
func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	tag := domain.NewTag(req.Name, req.Description)

	if err := h.tagService.Create(ctx, tag); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tag)
}

// UpdateTag actualiza una etiqueta existente
func (h *TagHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de etiqueta no proporcionado", http.StatusBadRequest)
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
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	tag, err := h.tagService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrTagNotFound {
			http.Error(w, "Etiqueta no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tag.Update(req.Name, req.Description)

	if err := h.tagService.Update(ctx, tag); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tag)
}

// DeleteTag elimina una etiqueta por su ID
func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de etiqueta no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.tagService.Delete(ctx, id)
	if err != nil {
		if err == domain.ErrTagNotFound {
			http.Error(w, "Etiqueta no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}