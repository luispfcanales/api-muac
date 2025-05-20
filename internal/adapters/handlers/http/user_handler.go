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

// GetUsers obtiene todos los usuarios
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserByID obtiene un usuario por su ID
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

// CreateUser crea un nuevo usuario
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userDTO struct {
		Name     string    `json:"name"`
		LastName string    `json:"lastname"`
		Username string    `json:"username"`
		Email    string    `json:"email"`
		DNI      string    `json:"dni"`
		Phone    string    `json:"phone"`
		Password string    `json:"password"`
		RoleID   uuid.UUID `json:"role_id"`
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
		userDTO.RoleID,
	)

	if err := h.userService.Create(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser actualiza un usuario existente
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

// DeleteUser elimina un usuario por su ID
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

// UpdatePassword actualiza la contraseña de un usuario
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

// UpdateRole actualiza el rol de un usuario
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
