// services/measurement_service.go (ACTUALIZADO CON MANEJO DE DUPLICADOS)
package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/ports"
)

// measurementService implementa la lógica de negocio para mediciones
type measurementService struct {
	measurementRepo ports.IMeasurementRepository
	tagRepo         ports.ITagRepository
	recommendRepo   ports.IRecommendationRepository
}

// NewMeasurementService crea una nueva instancia de MeasurementService
func NewMeasurementService(
	measurementRepo ports.IMeasurementRepository,
	tagRepo ports.ITagRepository,
	recommendRepo ports.IRecommendationRepository,
) ports.IMeasurementService {
	return &measurementService{
		measurementRepo: measurementRepo,
		tagRepo:         tagRepo,
		recommendRepo:   recommendRepo,
	}
}

// Create crea una nueva medición (método original - MANTIENE COMPATIBILIDAD)
func (s *measurementService) Create(ctx context.Context, measurement *domain.Measurement) error {
	if err := measurement.Validate(); err != nil {
		return err
	}
	return s.measurementRepo.Create(ctx, measurement)
}

// CreateWithAutoAssignment crea una nueva medición con asignación automática de tag y recomendación (ACTUALIZADO)
func (s *measurementService) CreateWithAutoAssignment(ctx context.Context, muacValue float64, description string, patientID, userID uuid.UUID) (*domain.Measurement, error) {
	// Validar valor MUAC
	if !domain.IsValidMuacValue(muacValue) {
		return nil, fmt.Errorf("valor MUAC inválido: %.2f", muacValue)
	}

	// Clasificar el valor MUAC
	muacCode, colorCode, priority := domain.ClassifyMuacValue(muacValue)

	// Obtener o crear tag apropiado
	tag, err := s.getOrCreateMuacTag(ctx, muacCode, colorCode, priority)
	if err != nil {
		return nil, fmt.Errorf("error al obtener tag MUAC: %w", err)
	}

	// Obtener recomendación apropiada
	recommendation, err := s.getOrCreateMuacRecommendation(ctx, muacValue, muacCode)
	if err != nil {
		return nil, fmt.Errorf("error al obtener recomendación MUAC: %w", err)
	}

	// Crear la medición con IDs asignados
	measurement := &domain.Measurement{
		ID:               uuid.New(),
		MuacValue:        muacValue,
		Description:      description,
		PatientID:        patientID,
		UserID:           userID,
		TagID:            &tag.ID,
		RecommendationID: &recommendation.ID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Validar y crear
	if err := measurement.Validate(); err != nil {
		return nil, err
	}

	if err := s.measurementRepo.Create(ctx, measurement); err != nil {
		return nil, err
	}

	// Cargar relaciones para retornar
	measurement.Tag = tag
	measurement.Recommendation = recommendation

	return measurement, nil
}

// getOrCreateMuacTag obtiene o crea el tag apropiado para el código MUAC (MÉTODO CORREGIDO)
func (s *measurementService) getOrCreateMuacTag(ctx context.Context, muacCode, colorCode string, priority int) (*domain.Tag, error) {
	// PASO 1: Intentar obtener tag existente por código MUAC si el repo lo soporta
	if tagRepoExtended, ok := s.tagRepo.(interface {
		GetByMuacCode(ctx context.Context, muacCode string) (*domain.Tag, error)
	}); ok {
		tag, err := tagRepoExtended.GetByMuacCode(ctx, muacCode)
		if err == nil && tag != nil {
			return tag, nil
		}
	}

	// PASO 2: Buscar por nombre del tag usando GetAll (FALLBACK PRINCIPAL)
	allTags, err := s.tagRepo.GetAll(ctx)
	if err == nil {
		// Buscar tag existente por código MUAC
		for _, tag := range allTags {
			if tag.MuacCode == muacCode && tag.Active {
				return tag, nil
			}
		}

		// Si no se encuentra por muac_code, buscar por nombre generado
		expectedName := s.getMuacTagName(muacCode)
		for _, tag := range allTags {
			if tag.Name == expectedName && tag.Active {
				// Actualizar el tag existente con el muac_code si no lo tiene
				if tag.MuacCode == "" {
					tag.MuacCode = muacCode
					tag.Color = colorCode
					tag.Priority = priority
					tag.UpdatedAt = time.Now()

					// Actualizar en la base de datos
					if updateErr := s.tagRepo.Update(ctx, tag); updateErr != nil {
						log.Printf("Warning: No se pudo actualizar tag existente: %v", updateErr)
					}
				}
				return tag, nil
			}
		}

		// Buscar por nombre similar (manejo de emojis/caracteres especiales)
		for _, tag := range allTags {
			if s.isTagNameSimilar(tag.Name, expectedName) && tag.Active {
				// Actualizar con muac_code si no lo tiene
				if tag.MuacCode == "" {
					tag.MuacCode = muacCode
					tag.Color = colorCode
					tag.Priority = priority
					tag.UpdatedAt = time.Now()

					if updateErr := s.tagRepo.Update(ctx, tag); updateErr != nil {
						log.Printf("Warning: No se pudo actualizar tag similar: %v", updateErr)
					}
				}
				return tag, nil
			}
		}
	}

	// PASO 3: Si no existe, crear uno nuevo
	name := s.getMuacTagName(muacCode)
	description := s.getMuacTagDescription(muacCode)

	newTag := domain.NewMuacTag(name, description, colorCode, muacCode, priority)

	if err := s.tagRepo.Create(ctx, newTag); err != nil {
		// Si hay error de duplicado, intentar buscar nuevamente
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			log.Printf("Tag duplicado detectado, buscando tag existente para: %s", muacCode)

			// Reintentiar búsqueda por si acaso otro proceso lo creó
			if allTags, retryErr := s.tagRepo.GetAll(ctx); retryErr == nil {
				for _, tag := range allTags {
					if tag.MuacCode == muacCode && tag.Active {
						return tag, nil
					}
					if tag.Name == name && tag.Active {
						return tag, nil
					}
					if s.isTagNameSimilar(tag.Name, name) && tag.Active {
						return tag, nil
					}
				}
			}
		}
		return nil, fmt.Errorf("error al crear tag MUAC: %w", err)
	}

	return newTag, nil
}

// getOrCreateMuacRecommendation obtiene o crea la recomendación apropiada (MÉTODO CORREGIDO)
func (s *measurementService) getOrCreateMuacRecommendation(ctx context.Context, muacValue float64, muacCode string) (*domain.Recommendation, error) {
	// PASO 1: Intentar obtener recomendaciones activas si el repo lo soporta
	if recommendRepoExtended, ok := s.recommendRepo.(interface {
		GetActiveRecommendations(ctx context.Context) ([]*domain.Recommendation, error)
	}); ok {
		recommendations, err := recommendRepoExtended.GetActiveRecommendations(ctx)
		if err == nil {
			// Buscar recomendación aplicable por muac_code primero
			for _, rec := range recommendations {
				if rec.MuacCode == muacCode && rec.IsApplicableForMuac(muacValue) {
					return rec, nil
				}
			}
			// Si no encuentra por muac_code, buscar por rango
			for _, rec := range recommendations {
				if rec.IsApplicableForMuac(muacValue) {
					return rec, nil
				}
			}
		}
	}

	// PASO 2: Usar GetAll como fallback
	allRecommendations, err := s.recommendRepo.GetAll(ctx)
	if err == nil {
		// Filtrar solo las activas
		activeRecs := domain.FilterActiveRecommendations(allRecommendations)

		// Buscar por código MUAC específico
		for _, rec := range activeRecs {
			if rec.MuacCode == muacCode && rec.IsApplicableForMuac(muacValue) {
				return rec, nil
			}
		}

		// Buscar por rango aplicable
		for _, rec := range activeRecs {
			if rec.IsApplicableForMuac(muacValue) {
				return rec, nil
			}
		}

		// Buscar por nombre si las anteriores fallan
		expectedName := s.getExpectedRecommendationName(muacCode)
		for _, rec := range activeRecs {
			if strings.Contains(rec.Name, expectedName) {
				// Actualizar muac_code si no lo tiene
				if rec.MuacCode == "" {
					rec.MuacCode = muacCode
					rec.UpdatedAt = time.Now()
					if updateErr := s.recommendRepo.Update(ctx, rec); updateErr != nil {
						log.Printf("Warning: No se pudo actualizar recomendación: %v", updateErr)
					}
				}
				return rec, nil
			}
		}
	}

	// PASO 3: Si no hay recomendaciones aplicables, crear una por defecto
	return s.createDefaultRecommendation(ctx, muacCode, muacValue)
}

// createDefaultRecommendation crea una recomendación por defecto completa y contextualizada (MEJORADO)
func (s *measurementService) createDefaultRecommendation(ctx context.Context, muacCode string, muacValue float64) (*domain.Recommendation, error) {
	var name, description string
	var minValue, maxValue *float64
	var priority int
	var colorCode string

	switch muacCode {
	case domain.MuacCodeRed:
		name = "🚨 ALERTA ROJA - Acción Urgente Requerida"
		description = "⚠️ Esta medición indica DESNUTRICIÓN AGUDA SEVERA (SAM). Tu niño o niña necesita atención médica URGENTE. No es tu culpa, pero sí es momento de actuar rápido.\n\n" +
			"ACCIONES INMEDIATAS:\n" +
			"1. 🏥 Acude HOY MISMO al establecimiento de salud más cercano\n" +
			"2. 🚫 No retrases la consulta, incluso si el niño parece estar bien\n" +
			"3. 💧 Mientras te trasladas: mantén hidratado con agua hervida, mates suaves\n" +
			"4. 🍌 Ofrece alimentos fáciles: plátano sancochado, puré de yuca, mazamorra\n" +
			"5. 📞 Si no puedes movilizarte: contacta al agente comunitario de salud\n" +
			"6. 🔄 Repite medición solo DESPUÉS de consulta médica\n\n" +
			"⚠️ Este resultado no sustituye diagnóstico médico. Es una herramienta de alerta familiar."
		severeThreshold := domain.MuacThresholdSevere
		maxValue = &severeThreshold
		priority = domain.PriorityUrgent
		colorCode = domain.ColorRed

	case domain.MuacCodeYellow:
		name = "🟡 ALERTA AMARILLA - Zona de Riesgo Nutricional"
		description = "🟡 Tu niño o niña está en RIESGO NUTRICIONAL (MAM). No es emergencia, pero es una señal importante. Es momento de fortalecer su alimentación.\n\n" +
			"ACCIONES RECOMENDADAS:\n" +
			"1. 🏥 Solicita evaluación en centro de salud en los próximos 5 días\n" +
			"2. 🍳 Mejora alimentación con productos locales:\n" +
			"   • Proteínas: huevos, pescado regional, sangrecita\n" +
			"   • Frutas amazónicas: camu camu, aguaje, cocona\n" +
			"   • Energía: plátano, quinua, lenteja, maní, maíz tierno\n" +
			"3. 🍽️ Aumenta frecuencia a 4-5 comidas diarias\n" +
			"4. 🚫 Evita ultraprocesados (galletas, gaseosas, embutidos)\n" +
			"5. 📅 Nuevo control MUAC en 7 días\n" +
			"6. 🌡️ Si hay fiebre, diarrea o pérdida de apetito: acude antes\n\n" +
			"💪 Con amor, buena comida y atención, tu niño/a puede recuperarse."
		severeThreshold := domain.MuacThresholdSevere
		moderateThreshold := domain.MuacThresholdModerate
		minValue = &severeThreshold
		maxValue = &moderateThreshold
		priority = domain.PriorityAttention
		colorCode = domain.ColorYellow

	case domain.MuacCodeGreen:
		name = "✅ ZONA VERDE - Estado Nutricional Adecuado"
		description = "✅ ¡Excelente! Tu niño o niña tiene BUEN ESTADO NUTRICIONAL. Sigue alimentándolo con cariño y atención para que continúe creciendo fuerte y sano.\n\n" +
			"ACCIONES PARA MANTENER LA SALUD:\n" +
			"1. 🥗 Mantén alimentación balanceada con productos locales:\n" +
			"   • Frutas amazónicas: copoazú, piña, camu camu\n" +
			"   • Proteínas: pescado, huevos, frejoles, hígado\n" +
			"   • Energía: yuca, plátano, arroz, maíz\n" +
			"   • Hierro/Vitamina A: sangrecita, zanahoria, sacha culantro\n" +
			"2. 📅 Controles CRED según edad (cada 2-3 meses)\n" +
			"3. 📏 Medición MUAC mensual o si baja el apetito\n" +
			"4. 🤝 Comparte esta herramienta con otras familias\n\n" +
			"🎉 ¡Felicitaciones por cuidar tan bien a tu niño/a!"
		normalThreshold := domain.MuacThresholdNormal
		minValue = &normalThreshold
		priority = domain.PriorityNormal
		colorCode = domain.ColorGreen

	case domain.MuacCodeFollow:
		name = "📋 Seguimiento Post-Intervención Nutricional"
		description = "📋 Tu niño o niña está en proceso de RECUPERACIÓN NUTRICIONAL. Mantener cuidados especiales y seguimiento médico es fundamental.\n\n" +
			"PROTOCOLO DE SEGUIMIENTO:\n" +
			"1. 💊 Continuar plan alimentario establecido por el centro de salud\n" +
			"2. 📅 Controles semanales obligatorios - NO faltar\n" +
			"3. ⚖️ Monitoreo de peso y talla regularmente\n" +
			"4. 🍳 Alimentación especial reforzada:\n" +
			"   • Comidas pequeñas y frecuentes (cada 2-3 horas)\n" +
			"   • Proteínas en cada comida: huevo, pescado, sangrecita\n" +
			"   • Aceites vegetales para agregar energía\n" +
			"   • Frutas ricas en vitaminas: aguaje, camu camu\n" +
			"5. 👨‍👩‍👧‍👦 Apoyo familiar: todos participan en la recuperación\n" +
			"6. 📱 Registro diario de alimentos consumidos\n" +
			"7. 🚨 Alerta inmediata si empeoran síntomas\n\n" +
			"⏰ La constancia en el seguimiento es clave para la recuperación completa."
		priority = domain.PriorityAttention
		colorCode = domain.ColorBlue

	default:
		name = "📋 Seguimiento General"
		description = "📋 Medición registrada en el sistema MUAC. Continúa con el protocolo de seguimiento nutricional establecido.\n\n" +
			"RECOMENDACIONES GENERALES:\n" +
			"1. 🍽️ Mantén alimentación balanceada y variada\n" +
			"2. 📏 Mediciones regulares según protocolo\n" +
			"3. 🏥 Consultas médicas programadas\n" +
			"4. 📊 Seguimiento del crecimiento y desarrollo\n\n" +
			"💡 Para recomendaciones específicas, consulta con personal de salud."
		priority = domain.PriorityNormal
		colorCode = domain.ColorGray
	}

	// Verificar si ya existe una recomendación similar
	if existingRec, exists := s.recommendationExists(ctx, name, muacCode); exists {
		return existingRec, nil
	}

	// Crear nueva recomendación
	recommendation := domain.NewMuacRecommendation(
		name,
		description,
		minValue,
		maxValue,
		priority,
		colorCode,
		muacCode,
	)

	if err := s.recommendRepo.Create(ctx, recommendation); err != nil {
		// Si hay error de duplicado, intentar buscar la existente
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			log.Printf("Recomendación duplicada detectada, buscando existente para: %s", muacCode)

			if existingRec, exists := s.recommendationExists(ctx, name, muacCode); exists {
				return existingRec, nil
			}
		}
		return nil, fmt.Errorf("error al crear recomendación por defecto: %w", err)
	}

	return recommendation, nil
}

// ============= MÉTODOS HELPER PRIVADOS =============

// isTagNameSimilar verifica si dos nombres de tags son similares (maneja emojis)
func (s *measurementService) isTagNameSimilar(name1, name2 string) bool {
	// Remover emojis y espacios extra para comparar
	clean1 := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(name1, "🚨", ""), "🟡", ""), "✅", ""))
	clean2 := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(name2, "🚨", ""), "🟡", ""), "✅", ""))

	return strings.Contains(clean1, "ALERTA ROJA") && strings.Contains(clean2, "ALERTA ROJA") ||
		strings.Contains(clean1, "ALERTA AMARILLA") && strings.Contains(clean2, "ALERTA AMARILLA") ||
		strings.Contains(clean1, "ZONA VERDE") && strings.Contains(clean2, "ZONA VERDE") ||
		strings.Contains(clean1, "SEGUIMIENTO") && strings.Contains(clean2, "SEGUIMIENTO")
}

// getExpectedRecommendationName retorna el nombre esperado para buscar recomendaciones existentes
func (s *measurementService) getExpectedRecommendationName(muacCode string) string {
	switch muacCode {
	case domain.MuacCodeRed:
		return "ALERTA ROJA"
	case domain.MuacCodeYellow:
		return "ALERTA AMARILLA"
	case domain.MuacCodeGreen:
		return "ZONA VERDE"
	case domain.MuacCodeFollow:
		return "Seguimiento"
	default:
		return "General"
	}
}

// recommendationExists verifica si una recomendación ya existe antes de crear
func (s *measurementService) recommendationExists(ctx context.Context, name string, muacCode string) (*domain.Recommendation, bool) {
	allRecommendations, err := s.recommendRepo.GetAll(ctx)
	if err != nil {
		return nil, false
	}

	for _, rec := range allRecommendations {
		if rec.Active && (rec.Name == name || rec.MuacCode == muacCode) {
			return rec, true
		}
		// Buscar por nombre similar
		if rec.Active && strings.Contains(rec.Name, s.getExpectedRecommendationName(muacCode)) {
			return rec, true
		}
	}
	return nil, false
}

// getMuacTagName retorna el nombre del tag según código MUAC (MÉTODO PRIVADO)
func (s *measurementService) getMuacTagName(muacCode string) string {
	switch muacCode {
	case domain.MuacCodeRed:
		return "🚨 ALERTA ROJA"
	case domain.MuacCodeYellow:
		return "🟡 ALERTA AMARILLA"
	case domain.MuacCodeGreen:
		return "✅ ZONA VERDE"
	case domain.MuacCodeFollow:
		return "📋 SEGUIMIENTO"
	default:
		return "⚪ SIN CLASIFICAR"
	}
}

// getMuacTagDescription retorna la descripción del tag según código MUAC (MÉTODO PRIVADO)
func (s *measurementService) getMuacTagDescription(muacCode string) string {
	switch muacCode {
	case domain.MuacCodeRed:
		return fmt.Sprintf("Desnutrición aguda severa (SAM) - < %.1f cm - Requiere atención urgente", domain.MuacThresholdSevere)
	case domain.MuacCodeYellow:
		return fmt.Sprintf("Desnutrición aguda moderada (MAM) - %.1f-%.1f cm - Requiere seguimiento", domain.MuacThresholdSevere, domain.MuacThresholdModerate)
	case domain.MuacCodeGreen:
		return fmt.Sprintf("Estado nutricional adecuado - ≥ %.1f cm - Mantener cuidados", domain.MuacThresholdNormal)
	case domain.MuacCodeFollow:
		return "Paciente en seguimiento post-intervención nutricional"
	default:
		return "Medición sin clasificación MUAC específica"
	}
}

// ============= TODOS LOS MÉTODOS ORIGINALES (SIN CAMBIOS) =============

// GetByID obtiene una medición por su ID
func (s *measurementService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Measurement, error) {
	return s.measurementRepo.GetByID(ctx, id)
}

// GetByPatientID obtiene mediciones por ID de paciente
func (s *measurementService) GetByPatientID(ctx context.Context, patientID uuid.UUID) ([]*domain.Measurement, error) {
	return s.measurementRepo.GetByPatientID(ctx, patientID)
}

// GetByUserID obtiene mediciones por ID de usuario
func (s *measurementService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Measurement, error) {
	return s.measurementRepo.GetByUserID(ctx, userID)
}

// GetByTagID obtiene mediciones por ID de etiqueta
func (s *measurementService) GetByTagID(ctx context.Context, tagID uuid.UUID) ([]*domain.Measurement, error) {
	return s.measurementRepo.GetByTagID(ctx, tagID)
}

// GetByRecommendationID obtiene mediciones por ID de recomendación
func (s *measurementService) GetByRecommendationID(ctx context.Context, recommendationID uuid.UUID) ([]*domain.Measurement, error) {
	return s.measurementRepo.GetByRecommendationID(ctx, recommendationID)
}

// GetByDateRange obtiene mediciones dentro de un rango de fechas
func (s *measurementService) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*domain.Measurement, error) {
	return s.measurementRepo.GetByDateRange(ctx, startDate, endDate)
}

// GetAll obtiene todas las mediciones
func (s *measurementService) GetAll(ctx context.Context) ([]*domain.Measurement, error) {
	return s.measurementRepo.GetAll(ctx)
}

// Update actualiza una medición existente
func (s *measurementService) Update(ctx context.Context, measurement *domain.Measurement) error {
	if err := measurement.Validate(); err != nil {
		return err
	}
	return s.measurementRepo.Update(ctx, measurement)
}

// Delete elimina una medición por su ID
func (s *measurementService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.measurementRepo.Delete(ctx, id)
}

// AssignTag asigna una etiqueta a una medición
func (s *measurementService) AssignTag(ctx context.Context, measurementID, tagID uuid.UUID) error {
	// Verificar que la medición existe
	measurement, err := s.measurementRepo.GetByID(ctx, measurementID)
	if err != nil {
		return err
	}

	// Verificar que la etiqueta existe
	if tagID != uuid.Nil {
		_, err = s.tagRepo.GetByID(ctx, tagID)
		if err != nil {
			return err
		}
	}

	// Asignar la etiqueta
	measurement.SetTag(&tagID)

	return s.measurementRepo.Update(ctx, measurement)
}

// AssignRecommendation asigna una recomendación a una medición
func (s *measurementService) AssignRecommendation(ctx context.Context, measurementID, recommendationID uuid.UUID) error {
	// Verificar que la medición existe
	measurement, err := s.measurementRepo.GetByID(ctx, measurementID)
	if err != nil {
		return err
	}

	// Verificar que la recomendación existe
	if recommendationID != uuid.Nil {
		_, err = s.recommendRepo.GetByID(ctx, recommendationID)
		if err != nil {
			return err
		}
	}

	// Asignar la recomendación
	measurement.SetRecommendation(&recommendationID)
	return s.measurementRepo.Update(ctx, measurement)
}
