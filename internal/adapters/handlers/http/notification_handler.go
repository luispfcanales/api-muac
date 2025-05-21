package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// NotificationHandler maneja las solicitudes HTTP relacionadas con notificaciones
type NotificationHandler struct {
	notificationService ports.INotificationService
}

// NewNotificationHandler crea una nueva instancia de NotificationHandler
func NewNotificationHandler(notificationService ports.INotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// RegisterRoutes registra las rutas del handler en el router
func (h *NotificationHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/notifications", h.GetNotifications)
	mux.HandleFunc("GET /api/notifications/{id}", h.GetNotificationByID)
	mux.HandleFunc("POST /api/notifications", h.CreateNotification)
	mux.HandleFunc("PUT /api/notifications/{id}", h.UpdateNotification)
	mux.HandleFunc("DELETE /api/notifications/{id}", h.DeleteNotification)
	mux.HandleFunc("PUT /api/notifications/{id}/visible", h.SetVisibility)
}

// GetNotifications obtiene todas las notificaciones
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	notifications, err := h.notificationService.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// GetNotificationByID obtiene una notificación por su ID
func (h *NotificationHandler) GetNotificationByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de notificación no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de notificación inválido", http.StatusBadRequest)
		return
	}

	notification, err := h.notificationService.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotificationNotFound {
			http.Error(w, "Notificación no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notification)
}

// CreateNotification crea una nueva notificación
func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var notificationDTO struct {
		Title   string `json:"title"`
		Body    string `json:"body"`
		Visible bool   `json:"visible"`
	}

	if err := json.NewDecoder(r.Body).Decode(&notificationDTO); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la petición", http.StatusBadRequest)
		return
	}

	notification := domain.NewNotification(
		notificationDTO.Title,
		notificationDTO.Body,
		notificationDTO.Visible,
	)

	if err := notification.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.notificationService.Create(r.Context(), notification); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notification)
}

// UpdateNotification actualiza una notificación existente
func (h *NotificationHandler) UpdateNotification(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de notificación no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de notificación inválido", http.StatusBadRequest)
		return
	}

	var notificationDTO struct {
		Title   string `json:"title"`
		Body    string `json:"body"`
		Visible bool   `json:"visible"`
	}

	if err = json.NewDecoder(r.Body).Decode(&notificationDTO); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la petición", http.StatusBadRequest)
		return
	}

	notification, err := h.notificationService.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotificationNotFound {
			http.Error(w, "Notificación no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notification.Update(
		notificationDTO.Title,
		notificationDTO.Body,
		notificationDTO.Visible,
	)

	if err := notification.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.notificationService.Update(r.Context(), notification); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notification)
}

// DeleteNotification elimina una notificación
func (h *NotificationHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de notificación no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de notificación inválido", http.StatusBadRequest)
		return
	}

	if err := h.notificationService.Delete(r.Context(), id); err != nil {
		if err == domain.ErrNotificationNotFound {
			http.Error(w, "Notificación no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SetVisibility actualiza la visibilidad de una notificación
func (h *NotificationHandler) SetVisibility(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de notificación no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de notificación inválido", http.StatusBadRequest)
		return
	}

	var visibilityDTO struct {
		Visible bool `json:"visible"`
	}

	if err = json.NewDecoder(r.Body).Decode(&visibilityDTO); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la petición", http.StatusBadRequest)
		return
	}

	notification, err := h.notificationService.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotificationNotFound {
			http.Error(w, "Notificación no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notification.SetVisible(visibilityDTO.Visible)

	if err := h.notificationService.Update(r.Context(), notification); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notification)
}
