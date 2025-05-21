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

// GetAllFAQs obtiene todas las preguntas frecuentes
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

// GetFAQByID obtiene una pregunta frecuente por su ID
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

// CreateFAQ crea una nueva pregunta frecuente
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

// UpdateFAQ actualiza una pregunta frecuente existente
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

// DeleteFAQ elimina una pregunta frecuente por su ID
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