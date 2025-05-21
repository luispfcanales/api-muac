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

// GetAllTags godoc
// @Summary Obtener todas las etiquetas
// @Description Obtiene una lista de todas las etiquetas registradas en el sistema
// @Tags etiquetas
// @Accept json
// @Produce json
// @Success 200 {array} domain.Tag
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/tags [get]
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

// GetTagByID godoc
// @Summary Obtener una etiqueta por ID
// @Description Obtiene una etiqueta específica por su ID
// @Tags etiquetas
// @Accept json
// @Produce json
// @Param id path string true "ID de la etiqueta"
// @Success 200 {object} domain.Tag
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Etiqueta no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/tags/{id} [get]
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

// GetTagByName godoc
// @Summary Obtener una etiqueta por nombre
// @Description Obtiene una etiqueta específica por su nombre
// @Tags etiquetas
// @Accept json
// @Produce json
// @Param name path string true "Nombre de la etiqueta"
// @Success 200 {object} domain.Tag
// @Failure 400 {object} map[string]string "Nombre no proporcionado"
// @Failure 404 {object} map[string]string "Etiqueta no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/tags/name/{name} [get]
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

// CreateTag godoc
// @Summary Crear una nueva etiqueta
// @Description Crea una nueva etiqueta con la información proporcionada
// @Tags etiquetas
// @Accept json
// @Produce json
// @Param tag body object true "Datos de la etiqueta"
// @Success 201 {object} domain.Tag
// @Failure 400 {object} map[string]string "Solicitud inválida"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/tags [post]
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

// UpdateTag godoc
// @Summary Actualizar una etiqueta
// @Description Actualiza una etiqueta existente con la información proporcionada
// @Tags etiquetas
// @Accept json
// @Produce json
// @Param id path string true "ID de la etiqueta"
// @Param tag body object true "Datos actualizados de la etiqueta"
// @Success 200 {object} domain.Tag
// @Failure 400 {object} map[string]string "ID inválido o solicitud inválida"
// @Failure 404 {object} map[string]string "Etiqueta no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/tags/{id} [put]
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

// DeleteTag godoc
// @Summary Eliminar una etiqueta
// @Description Elimina una etiqueta por su ID
// @Tags etiquetas
// @Accept json
// @Produce json
// @Param id path string true "ID de la etiqueta"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Etiqueta no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/tags/{id} [delete]
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