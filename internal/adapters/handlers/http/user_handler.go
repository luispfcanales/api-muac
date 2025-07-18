package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler maneja las peticiones HTTP relacionadas con usuarios
type UserHandler struct {
	userService ports.IUserService
	// excelService ports.IFileService
}

// NewUserHandler crea una nueva instancia de UserHandler
func NewUserHandler(userService ports.IUserService, excelService ports.IFileService) *UserHandler {
	return &UserHandler{
		userService: userService,
		// excelService: excelService,
	}
}

// RegisterRoutes registra las rutas del handler en el router
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	// mux.HandleFunc("GET /api/users/reporte/excel", h.GetApoderados)
	mux.HandleFunc("GET /api/users", h.GetUsers)
	mux.HandleFunc("POST /api/users/login", h.Login)
	mux.HandleFunc("POST /api/users", h.CreateUser)
	mux.HandleFunc("GET /api/users/{id}", h.GetUserByID)
	mux.HandleFunc("PUT /api/users/{id}", h.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", h.DeleteUser)
	mux.HandleFunc("PUT /api/users/{id}/password", h.UpdatePassword)
	mux.HandleFunc("PUT /api/users/{id}/role", h.UpdateRole)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		UsernameOrEmail string `json:"username_or_email"`
		Password        string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Error en los datos de entrada", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByUsernameOrEmail(
		r.Context(),
		loginRequest.UsernameOrEmail,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Usuario o contraseñas incorrectos", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password))
	if err != nil {
		http.Error(w, "Usuario o contraseña incorrectos", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
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
	// Extraer locality_id del query parameter
	var localityID *uuid.UUID
	if localityIDStr := r.URL.Query().Get("locality_id"); localityIDStr != "" {
		parsedID, err := uuid.Parse(localityIDStr)
		if err != nil {
			http.Error(w, "locality_id inválido: "+err.Error(), http.StatusBadRequest)
			return
		}
		localityID = &parsedID
	}

	users, err := h.userService.GetAll(r.Context(), localityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// func (h *UserHandler) GetApoderados(w http.ResponseWriter, r *http.Request) {
// 	// Extraer locality_id del query parameter (opcional)
// 	var localityID *uuid.UUID
// 	if localityIDStr := r.URL.Query().Get("locality_id"); localityIDStr != "" {
// 		parsedID, err := uuid.Parse(localityIDStr)
// 		if err != nil {
// 			http.Error(w, "locality_id inválido: "+err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		localityID = &parsedID
// 	}

// 	// Obtener usuarios
// 	users, err := h.userService.GetApoderados(r.Context(), localityID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Generar Excel
// 	excelData, err := h.excelService.GenerateApoderadosReport(r.Context(), users)
// 	if err != nil {
// 		http.Error(w, "Error generando reporte Excel: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(excelData)

// 	// // Configurar headers para descarga
// 	// filename := fmt.Sprintf("reporte_apoderados_%s.xlsx", time.Now().Format("2006-01-02_15-04-05"))
// 	// w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
// 	// w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
// 	// w.Header().Set("Content-Length", strconv.Itoa(len(excelData)))

// 	// // Escribir el archivo
// 	// w.Write(excelData)
// }

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

	userCreated, err := h.userService.GetByID(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userCreated)
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
		Name       string     `json:"name"`
		LastName   string     `json:"lastname"`
		Username   string     `json:"username"`
		Email      string     `json:"email"`
		DNI        string     `json:"dni"`
		Phone      string     `json:"phone"`
		Password   string     `json:"password,omitempty"`
		RoleID     uuid.UUID  `json:"role_id"`
		LocalityID *uuid.UUID `json:"locality_id,omitempty"`
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

	// Hashear la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
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

	user.Update(
		userDTO.Name,
		userDTO.LastName,
		userDTO.Username,
		userDTO.Email,
		userDTO.Phone,
		userDTO.DNI,
		passwordHash,
		userDTO.RoleID,
		userDTO.LocalityID,
	)

	if err := h.userService.Update(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userUpdated, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userUpdated)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Contraseña actualizada"})
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
