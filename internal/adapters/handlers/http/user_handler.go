package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler maneja las peticiones HTTP relacionadas con usuarios
type UserHandler struct {
	userService ports.IUserService
}

// NewUserHandler crea una nueva instancia de UserHandler
func NewUserHandler(userService ports.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRoutes registra las rutas del handler en el router
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/users", h.GetUsers)
	mux.HandleFunc("POST /api/users", h.CreateUser)
	mux.HandleFunc("GET /api/users/{id}", h.GetUserByID)
	mux.HandleFunc("PUT /api/users/{id}", h.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", h.DeleteUser)
	mux.HandleFunc("PUT /api/users/{id}/password", h.UpdatePassword)
	mux.HandleFunc("PUT /api/users/{id}/role", h.UpdateRole)
}

// GetUsers godoc
// @Summary Obtener todos los usuarios
// @Description Obtiene una lista de todos los usuarios registrados en el sistema
// @Tags usuarios
// @Accept json
// @Produce json
// @Success 200 {array} domain.User
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/users [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserByID godoc
// @Summary Obtener un usuario por ID
// @Description Obtiene un usuario específico por su ID
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 200 {object} domain.User
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Usuario no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser godoc
// @Summary Crear un nuevo usuario
// @Description Crea un nuevo usuario con la información proporcionada
// @Tags usuarios
// @Accept json
// @Produce json
// @Param user body object true "Datos del usuario"
// @Success 201 {object} domain.User
// @Failure 400 {object} map[string]string "Solicitud inválida"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userDTO struct {
		Name       string     `json:"name"`
		LastName   string     `json:"lastname"`
		Username   string     `json:"username"`
		Email      string     `json:"email"`
		DNI        string     `json:"dni"`
		Phone      string     `json:"phone"`
		Password   string     `json:"password"`
		LocalityID *uuid.UUID `json:"locality_id,omitempty"`

		RoleID uuid.UUID `json:"role_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la petición", http.StatusBadRequest)
		return
	}

	// Hashear la contraseña usando bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al hashear la contraseña", http.StatusInternalServerError)
		return
	}
	passwordHash := string(hashedPassword)

	user := domain.NewUser(
		userDTO.Name,
		userDTO.LastName,
		userDTO.Username,
		userDTO.DNI,
		userDTO.Phone,
		userDTO.Email,
		passwordHash,
		// [],
		userDTO.RoleID,
		userDTO.LocalityID,
	)

	if err := h.userService.Create(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Usuario creado exitosamente",
	})
}

// UpdateUser godoc
// @Summary Actualizar un usuario
// @Description Actualiza un usuario existente con la información proporcionada
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Param user body object true "Datos actualizados del usuario"
// @Success 200 {object} domain.User
// @Failure 400 {object} map[string]string "ID inválido o solicitud inválida"
// @Failure 404 {object} map[string]string "Usuario no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	var userDTO struct {
		Name     string    `json:"name"`
		LastName string    `json:"lastname"`
		Username string    `json:"username"`
		Email    string    `json:"email"`
		DNI      string    `json:"dni"`
		Phone    string    `json:"phone"`
		RoleID   uuid.UUID `json:"role_id"`
	}

	if err = json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la petición", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Update(
		userDTO.Name,
		userDTO.LastName,
		userDTO.Username,
		userDTO.Email,
		userDTO.Phone,
		userDTO.DNI,
		userDTO.RoleID,
	)

	if err := h.userService.Update(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeleteUser godoc
// @Summary Eliminar un usuario
// @Description Elimina un usuario por su ID
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Usuario no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	if err := h.userService.Delete(r.Context(), id); err != nil {
		if err == domain.ErrUserNotFound {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdatePassword godoc
// @Summary Actualizar contraseña de un usuario
// @Description Actualiza la contraseña de un usuario específico
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Param password body object true "Nueva contraseña"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido o contraseña no proporcionada"
// @Failure 404 {object} map[string]string "Usuario no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/users/{id}/password [put]
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	var passwordDTO struct {
		Password string `json:"password"`
	}

	if err = json.NewDecoder(r.Body).Decode(&passwordDTO); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la petición", http.StatusBadRequest)
		return
	}

	// Hashear la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al hashear la contraseña", http.StatusInternalServerError)
		return
	}
	passwordHash := string(hashedPassword)

	if err := h.userService.UpdatePassword(r.Context(), id, passwordHash); err != nil {
		if err == domain.ErrUserNotFound {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateRole godoc
// @Summary Actualizar rol de un usuario
// @Description Actualiza el rol de un usuario específico
// @Tags usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Param role body object true "ID del nuevo rol"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido o rol no proporcionado"
// @Failure 404 {object} map[string]string "Usuario no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/users/{id}/role [put]
func (h *UserHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	var roleDTO struct {
		RoleID uuid.UUID `json:"role_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&roleDTO); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la petición", http.StatusBadRequest)
		return
	}

	if err := h.userService.UpdateRole(r.Context(), id, roleDTO.RoleID); err != nil {
		if err == domain.ErrUserNotFound {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
