package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// RoleHandler maneja las peticiones HTTP relacionadas con roles
type RoleHandler struct {
	roleService ports.IRoleService
}

// NewRoleHandler crea una nueva instancia de RoleHandler
func NewRoleHandler(roleService ports.IRoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// CreateRoleRequest representa la solicitud para crear un rol
type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateRoleRequest representa la solicitud para actualizar un rol
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// RegisterRoutes registra las rutas del manejador
func (h *RoleHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/roles", h.GetAllRoles)
	mux.HandleFunc("POST /api/roles", h.CreateRole)
	mux.HandleFunc("GET /api/roles/{id}", h.GetRoleByID)
	mux.HandleFunc("PUT /api/roles/{id}", h.UpdateRole)
	mux.HandleFunc("DELETE /api/roles/{id}", h.DeleteRole)
}

// GetAllRoles godoc
// @Summary Obtener todos los roles
// @Description Obtiene una lista de todos los roles registrados en el sistema
// @Tags roles
// @Accept json
// @Produce json
// @Success 200 {array} domain.Role
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/roles [get]
func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	roles, err := h.roleService.GetAllRoles(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

// GetRoleByID godoc
// @Summary Obtener un rol por ID
// @Description Obtiene un rol específico por su ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "ID del rol"
// @Success 200 {object} domain.Role
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Rol no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/roles/{id} [get]
func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de rol no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	role, err := h.roleService.GetRoleByID(ctx, id)
	if err != nil {
		if err == domain.ErrRoleNotFound {
			http.Error(w, "Rol no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

// CreateRole godoc
// @Summary Crear un nuevo rol
// @Description Crea un nuevo rol con la información proporcionada
// @Tags roles
// @Accept json
// @Produce json
// @Param role body CreateRoleRequest true "Datos del rol"
// @Success 201 {object} domain.Role
// @Failure 400 {object} map[string]string "Solicitud inválida o nombre de rol vacío"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/roles [post]
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	role, err := h.roleService.CreateRole(ctx, req.Name, req.Description)
	if err != nil {
		if err == domain.ErrEmptyRoleName {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

// UpdateRole godoc
// @Summary Actualizar un rol
// @Description Actualiza un rol existente con la información proporcionada
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "ID del rol"
// @Param role body UpdateRoleRequest true "Datos actualizados del rol"
// @Success 200 {object} domain.Role
// @Failure 400 {object} map[string]string "ID inválido, solicitud inválida o nombre de rol vacío"
// @Failure 404 {object} map[string]string "Rol no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/roles/{id} [put]
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de rol no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req UpdateRoleRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	role, err := h.roleService.UpdateRole(ctx, id, req.Name, req.Description)
	if err != nil {
		if err == domain.ErrRoleNotFound {
			http.Error(w, "Rol no encontrado", http.StatusNotFound)
			return
		}
		if err == domain.ErrEmptyRoleName {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

// DeleteRole godoc
// @Summary Eliminar un rol
// @Description Elimina un rol por su ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "ID del rol"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Rol no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de rol no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.roleService.DeleteRole(ctx, id)
	if err != nil {
		if err == domain.ErrRoleNotFound {
			http.Error(w, "Rol no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
