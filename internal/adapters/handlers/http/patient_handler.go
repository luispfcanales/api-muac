package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// PatientHandler maneja las peticiones HTTP relacionadas con pacientes
type PatientHandler struct {
	patientService ports.IPatientService
	fileService    ports.IFileService // Agregar servicio de archivos
}

// NewPatientHandler crea una nueva instancia de PatientHandler
func NewPatientHandler(patientService ports.IPatientService, fileService ports.IFileService) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
		fileService:    fileService,
	}
}

// RegisterRoutes registra las rutas del manejador
func (h *PatientHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/patients", h.GetAllPatients)
	mux.HandleFunc("POST /api/patients", h.CreatePatient)
	mux.HandleFunc("POST /api/patients/with-file", h.CreatePatientWithFile)
	mux.HandleFunc("GET /api/patients/{id}", h.GetPatientByID)
	mux.HandleFunc("PUT /api/patients/{id}", h.UpdatePatient)
	mux.HandleFunc("DELETE /api/patients/{id}", h.DeletePatient)
	mux.HandleFunc("GET /api/patients/dni/{dni}", h.GetPatientByDNI)
	mux.HandleFunc("GET /api/patients/father/{fatherId}", h.GetPatientsByFatherID)
	mux.HandleFunc("GET /api/patients/measurements/{id}", h.GetPatientMeasurements)
	mux.HandleFunc("POST /api/patients/measurements/{id}", h.AddPatientMeasurement)
	mux.HandleFunc("POST /api/patients/{id}/upload-dni", h.UploadPatientDNI)
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
func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Name         string `json:"name"`
		Lastname     string `json:"lastname"`
		Gender       string `json:"gender"`
		Age          int    `json:"age"`
		BirthDate    string `json:"birth_date"`
		ArmSize      string `json:"arm_size"`
		Weight       string `json:"weight"`
		Size         string `json:"size"`
		Description  string `json:"description"`
		ConsentGiven bool   `json:"consent_given"`
		CreatedBy    string `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.CreatedBy)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	patient := domain.NewPatient(
		req.Name,
		req.Lastname,
		req.Gender,
		req.BirthDate,
		req.ArmSize,
		req.Weight,
		req.Size,
		req.Description,
		req.Age,
		req.ConsentGiven,
		&id,
	)

	patient.Validate()

	if err := h.patientService.Create(ctx, patient); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(patient)
}

// CreatePatientWithFile crea un nuevo paciente con archivo DNI
func (h *PatientHandler) CreatePatientWithFile(w http.ResponseWriter, r *http.Request) {
	var patient *domain.Patient
	ctx := r.Context()

	// Parsear multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}

	// Verificar si se enviaron los datos como JSON en un campo del formulario
	if patientData := r.FormValue("patient"); patientData != "" {
		var req struct {
			Name         string `json:"name"`
			Lastname     string `json:"lastname"`
			Gender       string `json:"gender"`
			Age          int    `json:"age"`
			BirthDate    string `json:"birth_date"`
			ArmSize      string `json:"arm_size"`
			Weight       string `json:"weight"`
			Size         string `json:"size"`
			Description  string `json:"description"`
			ConsentGiven bool   `json:"consent_given"`
			CreatedBy    string `json:"created_by"`
		}

		if err := json.Unmarshal([]byte(patientData), &req); err != nil {
			http.Error(w, "Datos del paciente inválidos", http.StatusBadRequest)
			return
		}

		id, err := uuid.Parse(req.CreatedBy)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		patient = domain.NewPatient(
			req.Name,
			req.Lastname,
			req.Gender,
			req.BirthDate,
			req.ArmSize,
			req.Weight,
			req.Size,
			req.Description,
			req.Age,
			req.ConsentGiven,
			&id,
		)
	} else {
		// Obtener datos directamente de los campos del formulario
		createdBy := r.FormValue("created_by")
		id, err := uuid.Parse(createdBy)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		patient = domain.NewPatient(
			r.FormValue("name"),
			r.FormValue("lastname"),
			r.FormValue("gender"),
			r.FormValue("birth_date"),
			r.FormValue("arm_size"),
			r.FormValue("weight"),
			r.FormValue("size"),
			r.FormValue("description"),
			0, // Age se puede parsear si es necesario
			r.FormValue("consent_given") == "true",
			&id,
		)
	}

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
	}

	// Crear paciente
	if err := h.patientService.Create(ctx, patient); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(patient)
}

// UploadPatientDNI sube o actualiza el archivo DNI de un paciente existente
func (h *PatientHandler) UploadPatientDNI(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	patientID := r.PathValue("id")

	// Parsear ID del paciente
	id, err := uuid.Parse(patientID)
	if err != nil {
		http.Error(w, "ID de paciente inválido", http.StatusBadRequest)
		return
	}

	// Verificar que el paciente existe
	patient, err := h.patientService.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Paciente no encontrado", http.StatusNotFound)
		return
	}

	// Parsear multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}

	// Obtener archivo
	file, header, err := r.FormFile("dni_file")
	if err != nil {
		http.Error(w, "Error al obtener archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Subir archivo
	fileInfo, err := h.fileService.UploadFile(ctx, file, header, "patients/dni")
	if err != nil {
		http.Error(w, "Error al subir archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Actualizar paciente con nueva URL del DNI
	patient.UrlDNI = fileInfo.URL

	if err := h.patientService.Update(ctx, patient); err != nil {
		http.Error(w, "Error al actualizar paciente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder con información del archivo subido
	response := struct {
		Message  string          `json:"message"`
		FileInfo *ports.FileInfo `json:"file_info"`
		Patient  *domain.Patient `json:"patient"`
	}{
		Message:  "Archivo DNI subido exitosamente",
		FileInfo: fileInfo,
		Patient:  patient,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdatePatient godoc
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
func (h *PatientHandler) UpdatePatient(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		Name         string `json:"name"`
		Lastname     string `json:"lastname"`
		Gender       string `json:"gender"`
		Age          int    `json:"age"`
		BirthDate    string `json:"birth_date"`
		ArmSize      string `json:"arm_size"`
		Weight       string `json:"weight"`
		Size         string `json:"size"`
		ConsentGiven bool   `json:"consent_given"`
		Description  string `json:"description"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
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

	patient.Update(
		req.Name,
		req.Lastname,
		req.Gender,
		req.BirthDate,
		req.ArmSize,
		req.Weight,
		req.Size,
		req.Description,
		req.Age,
		req.ConsentGiven,
	)

	if err := h.patientService.Update(ctx, patient); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
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

// AddPatientMeasurement añade una nueva medición a un paciente
func (h *PatientHandler) AddPatientMeasurement(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		MuacValue        float64   `json:"muac_value"`
		Description      string    `json:"description"`
		Location         string    `json:"location"`
		UserID           uuid.UUID `json:"user_id"`
		TagID            uuid.UUID `json:"tag_id"`
		RecommendationID uuid.UUID `json:"recommendation_id"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	measurement := domain.NewMeasurement(
		req.MuacValue,
		req.Description,
		req.Location,
		time.Now(),
		id,
		req.UserID,
		req.TagID,
		req.RecommendationID,
	)

	if err := h.patientService.AddMeasurement(ctx, id, measurement); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(measurement)
}
