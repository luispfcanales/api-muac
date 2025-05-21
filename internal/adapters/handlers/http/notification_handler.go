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

// GetNotifications godoc
// @Summary Obtener todas las notificaciones
// @Description Obtiene una lista de todas las notificaciones registradas en el sistema
// @Tags notificaciones
// @Accept json
// @Produce json
// @Success 200 {array} domain.Notification
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/notifications [get]
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	notifications, err := h.notificationService.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// GetNotificationByID godoc
// @Summary Obtener una notificación por ID
// @Description Obtiene una notificación específica por su ID
// @Tags notificaciones
// @Accept json
// @Produce json
// @Param id path string true "ID de la notificación"
// @Success 200 {object} domain.Notification
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Notificación no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/notifications/{id} [get]
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

// CreateNotification godoc
// @Summary Crear una nueva notificación
// @Description Crea una nueva notificación con la información proporcionada
// @Tags notificaciones
// @Accept json
// @Produce json
// @Param notification body object true "Datos de la notificación"
// @Success 201 {object} domain.Notification
// @Failure 400 {object} map[string]string "Solicitud inválida"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/notifications [post]
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

// UpdateNotification godoc
// @Summary Actualizar una notificación
// @Description Actualiza una notificación existente con la información proporcionada
// @Tags notificaciones
// @Accept json
// @Produce json
// @Param id path string true "ID de la notificación"
// @Param notification body object true "Datos actualizados de la notificación"
// @Success 200 {object} domain.Notification
// @Failure 400 {object} map[string]string "ID inválido o solicitud inválida"
// @Failure 404 {object} map[string]string "Notificación no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/notifications/{id} [put]
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

// DeleteNotification godoc
// @Summary Eliminar una notificación
// @Description Elimina una notificación por su ID
// @Tags notificaciones
// @Accept json
// @Produce json
// @Param id path string true "ID de la notificación"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Notificación no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/notifications/{id} [delete]
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

// SetVisibility godoc
// @Summary Actualizar visibilidad de una notificación
// @Description Actualiza el estado de visibilidad de una notificación específica
// @Tags notificaciones
// @Accept json
// @Produce json
// @Param id path string true "ID de la notificación"
// @Param visibility body object true "Estado de visibilidad" 
// @Success 200 {object} domain.Notification
// @Failure 400 {object} map[string]string "ID inválido o solicitud inválida"
// @Failure 404 {object} map[string]string "Notificación no encontrada"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/notifications/{id}/visible [put]
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
