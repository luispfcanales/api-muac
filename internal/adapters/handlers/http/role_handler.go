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
	roleService ports.RoleService
}

// NewRoleHandler crea una nueva instancia de RoleHandler
func NewRoleHandler(roleService ports.RoleService) *RoleHandler {
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
	mux.HandleFunc("/api/roles", h.handleRoles)
	mux.HandleFunc("/api/roles/", h.handleRoleByID)
}

// handleRoles maneja las peticiones GET y POST a /api/roles
func (h *RoleHandler) handleRoles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllRoles(w, r)
	case http.MethodPost:
		h.createRole(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// handleRoleByID maneja las peticiones GET, PUT y DELETE a /api/roles/{id}
func (h *RoleHandler) handleRoleByID(w http.ResponseWriter, r *http.Request) {
	// Extraer ID del path
	idStr := r.URL.Path[len("/api/roles/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getRoleByID(w, r, id)
	case http.MethodPut:
		h.updateRole(w, r, id)
	case http.MethodDelete:
		h.deleteRole(w, r, id)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// getAllRoles obtiene todos los roles
func (h *RoleHandler) getAllRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	roles, err := h.roleService.GetAllRoles(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

// getRoleByID obtiene un rol por su ID
func (h *RoleHandler) getRoleByID(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	ctx := r.Context()

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

// createRole crea un nuevo rol
func (h *RoleHandler) createRole(w http.ResponseWriter, r *http.Request) {
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

// updateRole actualiza un rol existente
func (h *RoleHandler) updateRole(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	ctx := r.Context()

	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

// deleteRole elimina un rol por su ID
func (h *RoleHandler) deleteRole(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	ctx := r.Context()

	err := h.roleService.DeleteRole(ctx, id)
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
