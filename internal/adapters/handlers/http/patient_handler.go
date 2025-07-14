package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// PatientHandler maneja las peticiones HTTP relacionadas con pacientes
type PatientHandler struct {
	patientService     ports.IPatientService
	measurementService ports.IMeasurementService
	fileService        ports.IFileService // Agregar servicio de archivos
}

// NewPatientHandler crea una nueva instancia de PatientHandler
func NewPatientHandler(patientService ports.IPatientService, measurementService ports.IMeasurementService, fileService ports.IFileService) *PatientHandler {
	return &PatientHandler{
		patientService:     patientService,
		measurementService: measurementService,
		fileService:        fileService,
	}
}

// RegisterRoutes registra las rutas del manejador
func (h *PatientHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/patients", h.GetAllPatients)
	// mux.HandleFunc("POST /api/patients", h.CreatePatient)
	mux.HandleFunc("GET /api/patients/patients-in-risk", h.GetPatientsInRisk)
	mux.HandleFunc("POST /api/patients/with-file", h.CreatePatientWithFile)
	mux.HandleFunc("GET /api/patients/{id}", h.GetPatientByID)
	mux.HandleFunc("PUT /api/patients/{id}", h.UpdatePatientWithFile)
	mux.HandleFunc("DELETE /api/patients/{id}", h.DeletePatient)
	mux.HandleFunc("GET /api/patients/dni/{dni}", h.GetPatientByDNI)
	mux.HandleFunc("GET /api/patients/father/{fatherId}", h.GetPatientsByFatherID)
	mux.HandleFunc("GET /api/patients/measurements/{id}", h.GetPatientMeasurements)
	mux.HandleFunc("POST /api/patients/measurements/{id}", h.AddPatientMeasurement)
	// mux.HandleFunc("POST /api/patients/upload-dni/{id}", h.UploadPatientDNI)
}

// GetAllPatients godoc
// @Summary Obtener todos los pacientes
// @Description Obtiene una lista de todos los pacientes registrados en el sistema
// @Tags pacientes
// @Accept json
// @Produce json
// @Success 200 {array} domain.Patient
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/patients [get]
func (h *PatientHandler) GetAllPatients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	patients, err := h.patientService.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patients)
}

// GetPatientByID godoc
// @Summary Obtener un paciente por ID
// @Description Obtiene un paciente específico por su ID
// @Tags pacientes
// @Accept json
// @Produce json
// @Param id path string true "ID del paciente"
// @Success 200 {object} domain.Patient
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Paciente no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/patients/{id} [get]
func (h *PatientHandler) GetPatientByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de paciente no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	patient, err := h.patientService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrPatientNotFound {
			http.Error(w, "Paciente no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
}

// GetPatientByDNI godoc
// @Summary Obtener un paciente por DNI
// @Description Obtiene un paciente específico por su número de DNI
// @Tags pacientes
// @Accept json
// @Produce json
// @Param dni path string true "DNI del paciente"
// @Success 200 {object} domain.Patient
// @Failure 400 {object} map[string]string "DNI no proporcionado"
// @Failure 404 {object} map[string]string "Paciente no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/patients/dni/{dni} [get]
func (h *PatientHandler) GetPatientByDNI(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dni := r.PathValue("dni")
	if dni == "" {
		http.Error(w, "DNI no proporcionado", http.StatusBadRequest)
		return
	}

	patient, err := h.patientService.GetByDNI(ctx, dni)
	if err != nil {
		if err == domain.ErrPatientNotFound {
			http.Error(w, "Paciente no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
}

// CreatePatient godoc
// @Summary Crear un nuevo paciente
// @Description Crea un nuevo paciente con la información proporcionada
// @Tags pacientes
// @Accept json
// @Produce json
// @Param patient body object true "Datos del paciente"
// @Success 201 {object} domain.Patient
// @Failure 400 {object} map[string]string "Solicitud inválida"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/patients [post]
//
// CreatePatientWithFile crea un nuevo paciente con datos de formulario
// CreatePatientWithFile crea un nuevo paciente con datos de formulario
func (h *PatientHandler) CreatePatientWithFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parsear multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}

	// Validar y parsear created_by
	createdBy := r.FormValue("created_by")
	if createdBy == "" {
		http.Error(w, "created_by es requerido", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(createdBy)
	if err != nil {
		http.Error(w, "created_by debe ser un UUID válido", http.StatusBadRequest)
		return
	}

	// Validar y parsear age
	ageStr := r.FormValue("age")
	if ageStr == "" {
		http.Error(w, "age es requerido", http.StatusBadRequest)
		return
	}

	age, err := strconv.ParseFloat(ageStr, 64)
	if err != nil {
		http.Error(w, "Edad debe ser un número válido", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	name := r.FormValue("name")
	lastname := r.FormValue("lastname")
	dni := r.FormValue("dni")

	if name == "" || lastname == "" || dni == "" {
		http.Error(w, "name, lastname y dni son campos requeridos", http.StatusBadRequest)
		return
	}

	// Crear paciente con datos del formulario
	patient := domain.NewPatient(
		name,
		lastname,
		r.FormValue("gender"),
		r.FormValue("birth_date"),
		r.FormValue("arm_size"),
		r.FormValue("weight"),
		r.FormValue("size"),
		r.FormValue("description"),
		age,
		dni,
		r.FormValue("consent_given") == "true",
		&userID,
	)

	// Variable para rastrear el ID del archivo subido
	var uploadedFileID string

	// Procesar archivo DNI si se proporciona
	if file, header, err := r.FormFile("dni_file"); err == nil {
		defer file.Close()

		// Subir archivo DNI
		fileInfo, err := h.fileService.UploadFile(ctx, file, header, "patients/dni")
		if err != nil {
			http.Error(w, "Error al subir archivo DNI: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Asignar URL del DNI al paciente
		patient.UrlDNI = fileInfo.URL

		// Extraer ID del archivo para poder eliminarlo si hay error
		// URL esperada: http://localhost:8003/files/patients/dni/b8e52703-959a-487e-af75-74e6d210fb01.jpg
		filename := filepath.Base(fileInfo.URL)                               // Obtiene: b8e52703-959a-487e-af75-74e6d210fb01.jpg
		uploadedFileID = strings.TrimSuffix(filename, filepath.Ext(filename)) // Obtiene: b8e52703-959a-487e-af75-74e6d210fb01

		// Validar que el ID extraído es un UUID válido
		if _, err := uuid.Parse(uploadedFileID); err != nil {
			log.Printf("[ Error ]: ID de archivo inválido extraído de URL %s -> %s", fileInfo.URL, uploadedFileID)
			// Intentar eliminar el archivo con el ID inválido de todos modos
			h.fileService.DeleteFileIfExists(ctx, uploadedFileID)
			http.Error(w, "Error interno al procesar archivo", http.StatusInternalServerError)
			return
		}

		log.Printf("[ Info ]: Archivo subido exitosamente - ID: %s, URL: %s", uploadedFileID, fileInfo.URL)
	}

	// Validar el paciente
	if err := patient.Validate(); err != nil {
		// Si hay un archivo subido, eliminarlo
		if uploadedFileID != "" {
			if deleteErr := h.fileService.DeleteFileIfExists(ctx, uploadedFileID); deleteErr != nil {
				log.Printf("[ Error al eliminar archivo DNI tras validación fallida ]: %v", deleteErr)
			} else {
				log.Printf("[ Archivo DNI eliminado tras validación fallida ]: %s", uploadedFileID)
			}
		}
		http.Error(w, "Datos del paciente inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Crear paciente en la base de datos
	if err := h.patientService.Create(ctx, patient); err != nil {
		// Si hay un archivo subido y falla la creación del paciente, eliminarlo
		if uploadedFileID != "" {
			if deleteErr := h.fileService.DeleteFileIfExists(ctx, uploadedFileID); deleteErr != nil {
				log.Printf("[ Error al eliminar archivo DNI tras fallo en creación ]: %v", deleteErr)
			} else {
				log.Printf("[ Archivo DNI eliminado exitosamente tras fallo en creación ]: %s", uploadedFileID)
			}
		}

		// Determinar el tipo de error para dar mejor feedback
		errorMessage := err.Error()
		if strings.Contains(strings.ToLower(errorMessage), "duplicate") ||
			strings.Contains(strings.ToLower(errorMessage), "unique") ||
			strings.Contains(strings.ToLower(errorMessage), "dni") {
			http.Error(w, "El DNI ya está registrado en el sistema", http.StatusConflict)
			return
		}

		http.Error(w, "Error al crear paciente: "+errorMessage, http.StatusInternalServerError)
		return
	}

	// Obtener el paciente completo por ID (con todas las relaciones)
	createdPatient, err := h.patientService.GetByID(ctx, patient.ID)
	if err != nil {
		log.Printf("[ Warning ]: Paciente creado pero error al obtener datos completos: %v", err)
		// No eliminar archivo aquí porque el paciente se creó exitosamente
		http.Error(w, "Paciente creado pero error al obtener datos completos", http.StatusInternalServerError)
		return
	}

	// Respuesta exitosa
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Paciente creado exitosamente",
		"patient": createdPatient,
	})
}

// UpdatePatientWithFile godoc
// @Summary Actualizar un paciente
// @Description Actualiza un paciente existente con la información proporcionada
// @Tags pacientes
// @Accept json
// @Produce json
// @Param id path string true "ID del paciente"
// @Param patient body object true "Datos actualizados del paciente"
// @Success 200 {object} domain.Patient
// @Failure 400 {object} map[string]string "ID inválido o solicitud inválida"
// @Failure 404 {object} map[string]string "Paciente no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/patients/{id} [put]
// UpdatePatientWithFile actualiza un paciente existente con sus datos y opcionalmente su archivo DNI
func (h *PatientHandler) UpdatePatientWithFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	patientID := r.PathValue("id")

	// Parsear ID del paciente
	id, err := uuid.Parse(patientID)
	if err != nil {
		http.Error(w, "ID de paciente inválido", http.StatusBadRequest)
		return
	}

	// Verificar que el paciente existe
	existingPatient, err := h.patientService.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Paciente no encontrado", http.StatusNotFound)
		return
	}

	// Parsear multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}

	// Parsear y validar campos opcionales (solo actualizar si se proporcionan)
	updatedPatient := *existingPatient // Copia del paciente existente

	// Actualizar campos si se proporcionan
	if name := r.FormValue("name"); name != "" {
		updatedPatient.Name = name
	}
	if lastname := r.FormValue("lastname"); lastname != "" {
		updatedPatient.Lastname = lastname
	}
	if dni := r.FormValue("dni"); dni != "" {
		updatedPatient.DNI = dni
	}
	if gender := r.FormValue("gender"); gender != "" {
		updatedPatient.Gender = gender
	}
	if birthDate := r.FormValue("birth_date"); birthDate != "" {
		updatedPatient.BirthDate = birthDate
	}
	if armSize := r.FormValue("arm_size"); armSize != "" {
		updatedPatient.ArmSize = armSize
	}
	if weight := r.FormValue("weight"); weight != "" {
		updatedPatient.Weight = weight
	}
	if size := r.FormValue("size"); size != "" {
		updatedPatient.Size = size
	}
	if description := r.FormValue("description"); description != "" {
		updatedPatient.Description = description
	}

	// Actualizar age si se proporciona
	if ageStr := r.FormValue("age"); ageStr != "" {
		if age, err := strconv.ParseFloat(ageStr, 64); err == nil {
			updatedPatient.Age = age
		} else {
			http.Error(w, "Edad debe ser un número válido", http.StatusBadRequest)
			return
		}
	}

	// Actualizar consent_given si se proporciona
	if consentStr := r.FormValue("consent_given"); consentStr != "" {
		updatedPatient.ConsentGiven = consentStr == "true"
	}

	// Variable para rastrear el ID del nuevo archivo subido
	var newUploadedFileID string
	var oldFileIDToDelete string

	// Procesar archivo DNI si se proporciona
	if file, header, err := r.FormFile("dni_file"); err == nil {
		defer file.Close()

		// Si el paciente ya tenía un archivo DNI, extraer su ID para eliminarlo después
		if existingPatient.UrlDNI != "" {
			filename := filepath.Base(existingPatient.UrlDNI)
			oldFileIDToDelete = strings.TrimSuffix(filename, filepath.Ext(filename))
		}

		// Subir nuevo archivo DNI
		fileInfo, err := h.fileService.UploadFile(ctx, file, header, "patients/dni")
		if err != nil {
			http.Error(w, "Error al subir archivo DNI: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Asignar nueva URL del DNI al paciente
		updatedPatient.UrlDNI = fileInfo.URL

		// Extraer ID del nuevo archivo para poder eliminarlo si hay error
		filename := filepath.Base(fileInfo.URL)
		newUploadedFileID = strings.TrimSuffix(filename, filepath.Ext(filename))

		// Validar que el ID extraído es un UUID válido
		if _, err := uuid.Parse(newUploadedFileID); err != nil {
			log.Printf("[ Error ]: ID de archivo inválido extraído de URL %s -> %s", fileInfo.URL, newUploadedFileID)
			// Intentar eliminar el archivo con el ID inválido
			h.fileService.DeleteFileIfExists(ctx, newUploadedFileID)
			http.Error(w, "Error interno al procesar archivo", http.StatusInternalServerError)
			return
		}

		log.Printf("[ Info ]: Nuevo archivo subido exitosamente - ID: %s, URL: %s", newUploadedFileID, fileInfo.URL)
	}

	// Validar el paciente actualizado
	if err := updatedPatient.Validate(); err != nil {
		// Si hay un nuevo archivo subido, eliminarlo
		if newUploadedFileID != "" {
			if deleteErr := h.fileService.DeleteFileIfExists(ctx, newUploadedFileID); deleteErr != nil {
				log.Printf("[ Error al eliminar nuevo archivo DNI tras validación fallida ]: %v", deleteErr)
			} else {
				log.Printf("[ Nuevo archivo DNI eliminado tras validación fallida ]: %s", newUploadedFileID)
			}
		}
		http.Error(w, "Datos del paciente inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Actualizar paciente en la base de datos
	if err := h.patientService.Update(ctx, &updatedPatient); err != nil {
		// Si hay un nuevo archivo subido y falla la actualización, eliminarlo
		if newUploadedFileID != "" {
			if deleteErr := h.fileService.DeleteFileIfExists(ctx, newUploadedFileID); deleteErr != nil {
				log.Printf("[ Error al eliminar nuevo archivo DNI tras fallo en actualización ]: %v", deleteErr)
			} else {
				log.Printf("[ Nuevo archivo DNI eliminado exitosamente tras fallo en actualización ]: %s", newUploadedFileID)
			}
		}

		// Determinar el tipo de error para dar mejor feedback
		errorMessage := err.Error()
		if strings.Contains(strings.ToLower(errorMessage), "duplicate") ||
			strings.Contains(strings.ToLower(errorMessage), "unique") ||
			strings.Contains(strings.ToLower(errorMessage), "dni") {
			http.Error(w, "El DNI ya está registrado en el sistema", http.StatusConflict)
			return
		}

		http.Error(w, "Error al actualizar paciente: "+errorMessage, http.StatusInternalServerError)
		return
	}

	// Si la actualización fue exitosa y había un archivo anterior, eliminarlo
	if oldFileIDToDelete != "" && newUploadedFileID != "" {
		if deleteErr := h.fileService.DeleteFileIfExists(ctx, oldFileIDToDelete); deleteErr != nil {
			log.Printf("[ Warning ]: No se pudo eliminar archivo DNI anterior: %v", deleteErr)
		} else {
			log.Printf("[ Info ]: Archivo DNI anterior eliminado exitosamente: %s", oldFileIDToDelete)
		}
	}

	// Obtener el paciente actualizado completo (con todas las relaciones)
	finalPatient, err := h.patientService.GetByID(ctx, updatedPatient.ID)
	if err != nil {
		log.Printf("[ Warning ]: Paciente actualizado pero error al obtener datos completos: %v", err)
		http.Error(w, "Paciente actualizado pero error al obtener datos completos", http.StatusInternalServerError)
		return
	}

	// Respuesta exitosa
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Paciente actualizado exitosamente",
		"patient": finalPatient,
	})
}

// DeletePatient godoc
// @Summary Eliminar un paciente
// @Description Elimina un paciente por su ID
// @Tags pacientes
// @Accept json
// @Produce json
// @Param id path string true "ID del paciente"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido o no proporcionado"
// @Failure 404 {object} map[string]string "Paciente no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/patients/{id} [delete]
func (h *PatientHandler) DeletePatient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de paciente no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.patientService.Delete(ctx, id)
	if err != nil {
		if err == domain.ErrPatientNotFound {
			http.Error(w, "Paciente no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetPatientsByFatherID obtiene los pacientes asociados a un padre específico
func (h *PatientHandler) GetPatientsByFatherID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fatherIDStr := r.PathValue("fatherId")
	if fatherIDStr == "" {
		http.Error(w, "ID de padre no proporcionado", http.StatusBadRequest)
		return
	}

	fatherID, err := uuid.Parse(fatherIDStr)
	if err != nil {
		http.Error(w, "ID de padre inválido", http.StatusBadRequest)
		return
	}

	patients, err := h.patientService.GetByFatherID(ctx, fatherID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patients)
}

// GetPatientMeasurements obtiene las mediciones de un paciente específico
func (h *PatientHandler) GetPatientMeasurements(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de paciente no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	measurements, err := h.patientService.GetMeasurements(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurements)
}

// // AddPatientMeasurement añade una nueva medición a un paciente
// func (h *PatientHandler) AddPatientMeasurement(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	idStr := r.PathValue("id")
// 	if idStr == "" {
// 		http.Error(w, "ID de paciente no proporcionado", http.StatusBadRequest)
// 		return
// 	}

// 	id, err := uuid.Parse(idStr)
// 	if err != nil {
// 		http.Error(w, "ID inválido", http.StatusBadRequest)
// 		return
// 	}

// 	var req struct {
// 		MuacValue        float64   `json:"muac_value"`
// 		Description      string    `json:"description"`
// 		UserID           uuid.UUID `json:"user_id"`
// 		TagID            uuid.UUID `json:"tag_id"`
// 		RecommendationID uuid.UUID `json:"recommendation_id"`
// 	}

// 	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
// 		return
// 	}

// 	measurement := domain.NewMeasurement(
// 		req.MuacValue,
// 		req.Description,
// 		time.Now(),
// 		id,
// 		req.UserID,
// 		&req.TagID,
// 		&req.RecommendationID,
// 	)

// 	if err := h.patientService.AddMeasurement(ctx, id, measurement); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"message": "Medición agregada exitosamente",
// 	})
// }

// AddPatientMeasurement añade una nueva medición a un paciente con asignación automática de tag y recomendación
func (h *PatientHandler) AddPatientMeasurement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Obtener ID del paciente desde la URL
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "ID de paciente no proporcionado", http.StatusBadRequest)
		return
	}

	patientID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de paciente inválido", http.StatusBadRequest)
		return
	}

	// Estructura de request simplificada - solo necesitamos los datos básicos
	var req struct {
		MuacValue   float64   `json:"muac_value" validate:"required,gt=0"`
		Description string    `json:"description"`
		UserID      uuid.UUID `json:"user_id" validate:"required"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validaciones básicas
	if req.MuacValue <= 0 {
		http.Error(w, "El valor MUAC debe ser mayor a 0", http.StatusBadRequest)
		return
	}

	if req.MuacValue > 50 {
		http.Error(w, "El valor MUAC debe ser menor a 50 cm", http.StatusBadRequest)
		return
	}

	if req.UserID == uuid.Nil {
		http.Error(w, "ID de usuario es requerido", http.StatusBadRequest)
		return
	}

	// Verificar que el paciente existe
	patient, err := h.patientService.GetByID(ctx, patientID)
	if err != nil {
		if errors.Is(err, domain.ErrPatientNotFound) {
			http.Error(w, "Paciente no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al verificar paciente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Usar el servicio de mediciones con asignación automática
	measurement, err := h.measurementService.CreateWithAutoAssignment(
		ctx,
		req.MuacValue,
		req.Description,
		patientID,
		req.UserID,
	)

	if err != nil {
		// Manejar diferentes tipos de errores
		switch {
		case strings.Contains(err.Error(), "valor MUAC inválido"):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case strings.Contains(err.Error(), "usuario no encontrado"):
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		default:
			log.Printf("Error creando medición con auto-asignación: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	// Preparar respuesta con toda la información
	response := map[string]interface{}{
		"success": true,
		"message": "Medición agregada exitosamente con clasificación automática",
		"data": map[string]interface{}{
			"measurement": map[string]interface{}{
				"id":          measurement.ID,
				"muac_value":  measurement.MuacValue,
				"description": measurement.Description,
				"patient_id":  measurement.PatientID,
				"user_id":     measurement.UserID,
				"created_at":  measurement.CreatedAt,
			},
			"patient": map[string]interface{}{
				"id":       patient.ID,
				"name":     patient.Name,
				"lastname": patient.Lastname,
			},
			"classification": map[string]interface{}{
				"tag": map[string]interface{}{
					"id":          measurement.Tag.ID,
					"name":        measurement.Tag.Name,
					"description": measurement.Tag.Description,
					"color":       measurement.Tag.Color,
					"muac_code":   measurement.Tag.MuacCode,
					"priority":    measurement.Tag.Priority,
				},
				"recommendation": map[string]interface{}{
					"id":                    measurement.Recommendation.ID,
					"name":                  measurement.Recommendation.Name,
					"description":           measurement.Recommendation.Description,
					"recommendation_umbral": measurement.Recommendation.RecommendationUmbral,
					"priority":              measurement.Recommendation.Priority,
					"color_code":            measurement.Recommendation.ColorCode,
					"muac_code":             measurement.Recommendation.MuacCode,
				},
			},
			"muac_analysis": map[string]interface{}{
				"risk_level":     domain.GetMuacRiskLevel(req.MuacValue),
				"threshold_info": getMuacThresholdInfo(req.MuacValue),
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// GetPatientsInRisk obtiene pacientes en riesgo
func (h *PatientHandler) GetPatientsInRisk(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filters, err := h.parseFilters(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	patients, err := h.patientService.GetPatientsInRisk(ctx, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Pacientes en riesgo obtenidos exitosamente",
		"count":   len(patients),
		"data":    patients,
	})
}

// parseFilters parsea los query parameters a filtros
func (h *PatientHandler) parseFilters(r *http.Request) (*domain.ReportFilters, error) {
	filters := &domain.ReportFilters{}

	// Locality ID
	if localityIDStr := r.URL.Query().Get("locality_id"); localityIDStr != "" {
		localityID, err := uuid.Parse(localityIDStr)
		if err != nil {
			return nil, fmt.Errorf("locality_id inválido: %v", err)
		}
		filters.LocalityID = &localityID
	}

	// User ID
	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("user_id inválido: %v", err)
		}
		filters.UserID = &userID
	}

	// Days
	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return nil, fmt.Errorf("days debe ser un número válido: %v", err)
		}
		if days < 0 {
			return nil, fmt.Errorf("days no puede ser negativo")
		}
		if days > 365 {
			return nil, fmt.Errorf("days no puede ser mayor a 365")
		}
		filters.Days = days
	} else {
		filters.Days = 30 // Por defecto últimos 30 días
	}

	// Limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("limit debe ser un número válido: %v", err)
		}
		if limit < 0 {
			return nil, fmt.Errorf("limit no puede ser negativo")
		}
		if limit > 1000 {
			return nil, fmt.Errorf("limit no puede ser mayor a 1000")
		}
		filters.Limit = limit
	}

	return filters, nil
}

// getMuacThresholdInfo proporciona información contextual sobre los umbrales MUAC
func getMuacThresholdInfo(muacValue float64) map[string]interface{} {
	info := map[string]interface{}{
		"measured_value": muacValue,
		"thresholds": map[string]float64{
			"severe_malnutrition":   domain.MuacThresholdSevere,   // < 11.5 cm
			"moderate_malnutrition": domain.MuacThresholdModerate, // 11.5-12.4 cm
			"normal_nutrition":      domain.MuacThresholdNormal,   // >= 12.5 cm
		},
	}

	// Agregar contexto específico
	switch {
	case muacValue < domain.MuacThresholdSevere:
		info["status"] = "severe_acute_malnutrition"
		info["action_required"] = "urgent_medical_attention"
		info["priority"] = "critical"
	case muacValue < domain.MuacThresholdModerate:
		info["status"] = "moderate_acute_malnutrition"
		info["action_required"] = "nutritional_support"
		info["priority"] = "high"
	default:
		info["status"] = "adequate_nutritional_state"
		info["action_required"] = "maintain_current_care"
		info["priority"] = "normal"
	}

	return info
}
