package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// FAQHandler maneja las peticiones HTTP relacionadas con preguntas frecuentes
type FAQHandler struct {
	faqService ports.IFAQService
}

// NewFAQHandler crea una nueva instancia de FAQHandler
func NewFAQHandler(faqService ports.IFAQService) *FAQHandler {
	return &FAQHandler{
		faqService: faqService,
	}
}

// RegisterRoutes registra las rutas del manejador
func (h *FAQHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/faqs", h.GetAllFAQs)
	mux.HandleFunc("POST /api/faqs", h.CreateFAQ)
	mux.HandleFunc("GET /api/faqs/{id}", h.GetFAQByID)
	mux.HandleFunc("PUT /api/faqs/{id}", h.UpdateFAQ)
	mux.HandleFunc("DELETE /api/faqs/{id}", h.DeleteFAQ)
}

// GetAllFAQs godoc
// @Summary Obtener todas las preguntas frecuentes
// @Description Obtiene una lista de todas las preguntas frecuentes registradas
// @Tags faqs
// @Accept json
// @Produce json
// @Success 200 {array} domain.FAQ
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/faqs [get]
func (h *FAQHandler) GetAllFAQs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	faqs, err := h.faqService.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faqs)
}

// GetFAQByID godoc
// @Summary Obtener una pregunta frecuente por ID
// @Description Obtiene una pregunta frecuente específica por su ID
// @Tags faqs
// @Accept json
// @Produce json
// @Param id path string true "ID de la pregunta frecuente"
// @Success 200 {object} domain.FAQ
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Pregunta frecuente no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/faqs/{id} [get]
func (h *FAQHandler) GetFAQByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de FAQ no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	faq, err := h.faqService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrFAQNotFound {
			http.Error(w, "Pregunta frecuente no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faq)
}

// CreateFAQ godoc
// @Summary Crear una nueva pregunta frecuente
// @Description Crea una nueva pregunta frecuente con la información proporcionada
// @Tags faqs
// @Accept json
// @Produce json
// @Param faq body object true "Datos de la pregunta frecuente"
// @Success 201 {object} domain.FAQ
// @Failure 400 {object} map[string]string "Solicitud inválida"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/faqs [post]
func (h *FAQHandler) CreateFAQ(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	faq := domain.NewFAQ(req.Question, req.Answer)

	if err := h.faqService.Create(ctx, faq); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(faq)
}

// UpdateFAQ godoc
// @Summary Actualizar una pregunta frecuente
// @Description Actualiza una pregunta frecuente existente con la información proporcionada
// @Tags faqs
// @Accept json
// @Produce json
// @Param id path string true "ID de la pregunta frecuente"
// @Param faq body object true "Datos actualizados de la pregunta frecuente"
// @Success 200 {object} domain.FAQ
// @Failure 400 {object} map[string]string "ID inválido o solicitud inválida"
// @Failure 404 {object} map[string]string "Pregunta frecuente no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/faqs/{id} [put]
func (h *FAQHandler) UpdateFAQ(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de FAQ no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	faq, err := h.faqService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrFAQNotFound {
			http.Error(w, "Pregunta frecuente no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	faq.Update(req.Question, req.Answer)

	if err := h.faqService.Update(ctx, faq); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faq)
}

// DeleteFAQ godoc
// @Summary Eliminar una pregunta frecuente
// @Description Elimina una pregunta frecuente por su ID
// @Tags faqs
// @Accept json
// @Produce json
// @Param id path string true "ID de la pregunta frecuente"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Pregunta frecuente no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/faqs/{id} [delete]
func (h *FAQHandler) DeleteFAQ(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de FAQ no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.faqService.Delete(ctx, id)
	if err != nil {
		if err == domain.ErrFAQNotFound {
			http.Error(w, "Pregunta frecuente no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
