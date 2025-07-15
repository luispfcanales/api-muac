package config

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDatabase inserta datos iniciales basados en estándares OMS/UNICEF/Sphere Handbook
func SeedDatabase(db *gorm.DB) error {
	log.Println("🌱 Iniciando siembra de datos para Sistema MUAC (OMS/UNICEF/Sphere)")

	// Verificar si ya existen datos
	var roleCount int64
	if err := db.Model(&domain.Role{}).Count(&roleCount).Error; err != nil {
		return fmt.Errorf("error verificando roles existentes: %w", err)
	}

	if roleCount > 0 {
		log.Println("📋 Roles existentes detectados, verificando datos complementarios...")
		return seedAdditionalData(db)
	}

	// Iniciar transacción
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("error iniciando transacción: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Transacción revertida debido a panic: %v", r)
		}
	}()

	// Sembrar datos completos
	if err := seedRoles(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("error sembrando roles: %w", err)
	}

	if err := seedTags(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("error sembrando tags: %w", err)
	}

	if err := seedRecommendations(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("error sembrando recomendaciones: %w", err)
	}

	if err := seedAdminUser(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("error creando usuario admin: %w", err)
	}

	if err := seedFAQs(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("error creando FAQs: %w", err)
	}

	// Confirmar transacción
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error confirmando transacción: %w", err)
	}

	logSeedingSummary(db)
	return nil
}

// ============= FUNCIONES DE SIEMBRA ESPECÍFICAS =============

// seedRoles crea los roles del sistema MUAC
func seedRoles(tx *gorm.DB) error {
	log.Println("👥 Creando roles del sistema...")

	roles := []domain.Role{
		{
			ID:          uuid.New(),
			Name:        "ADMINISTRADOR",
			Description: "Acceso completo al sistema MUAC - Gestión de usuarios, configuración y reportes",
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "SUPERVISOR",
			Description: "Supervisión de apoderados, análisis de mediciones y reportes regionales",
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "APODERADO",
			Description: "Registro de mediciones MUAC de pacientes asignados en campo",
			CreatedAt:   time.Now(),
		},
	}

	if err := tx.Create(&roles).Error; err != nil {
		return fmt.Errorf("error creando roles: %w", err)
	}

	log.Printf("✅ %d roles creados exitosamente", len(roles))
	return nil
}

// seedTags crea los tags MUAC según estándares oficiales
func seedTags(tx *gorm.DB) error {
	log.Println("🏷️  Creando tags de clasificación MUAC...")

	tags := []domain.Tag{
		{
			ID:          uuid.New(),
			Name:        "MUAC-R1",
			Description: fmt.Sprintf("Alerta Roja - Desnutrición aguda severa (SAM) - < %.1f cm", domain.MuacThresholdSevere),
			Color:       domain.ColorRed,
			Active:      true,
			MuacCode:    domain.MuacCodeRed,
			Priority:    domain.PriorityExtreme,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "MUAC-Y1",
			Description: fmt.Sprintf("Alerta Amarilla - Desnutrición aguda moderada (MAM) - %.1f-%.1f cm", domain.MuacThresholdSevere, domain.MuacThresholdModerate),
			Color:       domain.ColorYellow,
			Active:      true,
			MuacCode:    domain.MuacCodeYellow,
			Priority:    domain.PriorityHigh,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "MUAC-G1",
			Description: fmt.Sprintf("Zona Verde - Estado nutricional adecuado - ≥ %.1f cm", domain.MuacThresholdNormal),
			Color:       domain.ColorGreen,
			Active:      true,
			MuacCode:    domain.MuacCodeGreen,
			Priority:    domain.PriorityLow,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "SEGUIMIENTO",
			Description: "Paciente en seguimiento post-intervención nutricional",
			Color:       domain.ColorBlue,
			Active:      true,
			MuacCode:    domain.MuacCodeFollow,
			Priority:    domain.PriorityMedium,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := tx.Create(&tags).Error; err != nil {
		return fmt.Errorf("error creando tags: %w", err)
	}

	log.Printf("✅ %d tags MUAC oficiales creados", len(tags))
	return nil
}

// seedRecommendations crea las recomendaciones nutricionales contextualizadas
func seedRecommendations(tx *gorm.DB) error {
	log.Println("💡 Creando recomendaciones nutricionales para comunidades amazónicas...")

	// Valores según estándares OMS/UNICEF
	valorSevere := domain.MuacThresholdSevere
	valorModerate := domain.MuacThresholdModerate
	valorNormal := domain.MuacThresholdNormal

	recommendations := []domain.Recommendation{
		{
			ID:   uuid.New(),
			Name: "🚨 ALERTA ROJA - Acción Urgente Requerida",
			Description: "⚠️ Esta medición indica DESNUTRICIÓN AGUDA SEVERA (SAM). Tu niño o niña necesita atención médica URGENTE. No es tu culpa, pero sí es momento de actuar rápido.\n\n" +
				"ACCIONES INMEDIATAS:\n" +
				"1. 🏥 Acude HOY MISMO al establecimiento de salud más cercano\n" +
				"2. 🚫 No retrases la consulta, incluso si el niño parece estar bien\n" +
				"3. 💧 Mientras te trasladas: mantén hidratado con agua hervida, mates suaves\n" +
				"4. 🍌 Ofrece alimentos fáciles: plátano sancochado, puré de yuca, mazamorra\n" +
				"5. 📞 Si no puedes movilizarte: contacta al agente comunitario de salud\n" +
				"6. 🔄 Repite medición solo DESPUÉS de consulta médica\n\n" +
				"⚠️ Este resultado no sustituye diagnóstico médico. Es una herramienta de alerta familiar.",
			RecommendationUmbral: fmt.Sprintf("< %.1f cm", valorSevere),
			MinValue:             nil,
			MaxValue:             &valorSevere,
			Priority:             domain.PriorityUrgent,
			Active:               true,
			ColorCode:            domain.ColorRed,
			MuacCode:             domain.MuacCodeRed,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		},
		{
			ID:   uuid.New(),
			Name: "🟡 ALERTA AMARILLA - Zona de Riesgo Nutricional",
			Description: "🟡 Tu niño o niña está en RIESGO NUTRICIONAL (MAM). No es emergencia, pero es una señal importante. Es momento de fortalecer su alimentación.\n\n" +
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
				"💪 Con amor, buena comida y atención, tu niño/a puede recuperarse.",
			RecommendationUmbral: fmt.Sprintf("%.1f - %.1f cm", valorSevere, valorModerate),
			MinValue:             &valorSevere,
			MaxValue:             &valorModerate,
			Priority:             domain.PriorityAttention,
			Active:               true,
			ColorCode:            domain.ColorYellow,
			MuacCode:             domain.MuacCodeYellow,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		},
		{
			ID:   uuid.New(),
			Name: "✅ ZONA VERDE - Estado Nutricional Adecuado",
			Description: "✅ ¡Excelente! Tu niño o niña tiene BUEN ESTADO NUTRICIONAL. Sigue alimentándolo con cariño y atención para que continúe creciendo fuerte y sano.\n\n" +
				"ACCIONES PARA MANTENER LA SALUD:\n" +
				"1. 🥗 Mantén alimentación balanceada con productos locales:\n" +
				"   • Frutas amazónicas: copoazú, piña, camu camu\n" +
				"   • Proteínas: pescado, huevos, frejoles, hígado\n" +
				"   • Energía: yuca, plátano, arroz, maíz\n" +
				"   • Hierro/Vitamina A: sangrecita, zanahoria, sacha culantro\n" +
				"2. 📅 Controles CRED según edad (cada 2-3 meses)\n" +
				"3. 📏 Medición MUAC mensual o si baja el apetito\n" +
				"4. 🤝 Comparte esta herramienta con otras familias\n\n" +
				"🎉 ¡Felicitaciones por cuidar tan bien a tu niño/a!",
			RecommendationUmbral: fmt.Sprintf("≥ %.1f cm", valorNormal),
			MinValue:             &valorNormal,
			MaxValue:             nil,
			Priority:             domain.PriorityNormal,
			Active:               true,
			ColorCode:            domain.ColorGreen,
			MuacCode:             domain.MuacCodeGreen,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		},
		{
			ID:   uuid.New(),
			Name: "📋 Seguimiento Post-Intervención",
			Description: "📋 Paciente en proceso de RECUPERACIÓN NUTRICIONAL. Mantener cuidados especiales y seguimiento médico.\n\n" +
				"PROTOCOLO DE SEGUIMIENTO:\n" +
				"1. 💊 Continuar plan alimentario establecido por el centro de salud\n" +
				"2. 📅 Controles semanales obligatorios\n" +
				"3. ⚖️ Monitoreo de peso y talla regularmente\n" +
				"4. 👨‍👩‍👧‍👦 Apoyo psicosocial a la familia\n" +
				"5. 📱 Registro diario de alimentos consumidos\n" +
				"6. 🚨 Alerta inmediata si empeoran síntomas\n\n" +
				"⏰ El seguimiento constante es clave para la recuperación completa.",
			RecommendationUmbral: "Todas las mediciones",
			MinValue:             nil,
			MaxValue:             nil,
			Priority:             domain.PriorityAttention,
			Active:               true,
			ColorCode:            domain.ColorBlue,
			MuacCode:             domain.MuacCodeFollow,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		},
	}

	if err := tx.Create(&recommendations).Error; err != nil {
		return fmt.Errorf("error creando recomendaciones: %w", err)
	}

	log.Printf("✅ %d recomendaciones contextualizadas creadas", len(recommendations))
	return nil
}

// seedAdminUser crea el usuario administrador inicial
func seedAdminUser(tx *gorm.DB) error {
	log.Println("👤 Creando usuario administrador inicial...")

	// Obtener rol de administrador
	var adminRole domain.Role
	if err := tx.Where("name = ?", "ADMINISTRADOR").First(&adminRole).Error; err != nil {
		return fmt.Errorf("rol administrador no encontrado: %w", err)
	}

	// Hashear contraseña
	password := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hasheando contraseña: %w", err)
	}

	adminUser := domain.User{
		ID:           uuid.New(),
		Name:         "ADMINISTRADOR",
		LastName:     "Sistema MUAC",
		Username:     "admin",
		Email:        "admin@muac.org",
		DNI:          "00000000",
		Phone:        "999000000",
		PasswordHash: string(hashedPassword),
		Active:       true,
		RoleID:       adminRole.ID,
		CreatedAt:    time.Now(),
	}

	if err := tx.Create(&adminUser).Error; err != nil {
		return fmt.Errorf("error creando usuario admin: %w", err)
	}

	log.Println("✅ Usuario administrador creado exitosamente")
	return nil
}

// seedFAQs crea las preguntas frecuentes iniciales del sistema
func seedFAQs(tx *gorm.DB) error {
	log.Println("❓ Creando preguntas frecuentes (FAQs)...")

	faqs := []domain.FAQ{
		// SOBRE EL USO DE LA CINTA Y EL APP
		{
			Question: "¿Qué significa la medida que ingreso en la app?",
			Answer:   "La medida que ingresas es la circunferencia del brazo de tu niño o niña, en centímetros. Esta medida nos ayuda a saber si está en buen estado nutricional o si necesita atención médica.",
			Category: domain.FAQCategoryTapeAndApp,
		},
		{
			Question: "¿Cómo sé si usé bien la cinta MUAC?",
			Answer:   "La cinta debe colocarse a la mitad del brazo izquierdo del niño, entre el hombro y el codo. No debe estar ni muy floja ni muy apretada. El número que se muestra en la ventana es el que debes ingresar en la app. Puedes repetir la medición si no estás seguro.",
			Category: domain.FAQCategoryTapeAndApp,
		},
		{
			Question: "¿Puedo usar la cinta con ropa puesta?",
			Answer:   "No. Para que la medición sea correcta, el brazo del niño/a debe estar sin ropa (manga arremangada o brazo desnudo). La ropa puede alterar el resultado.",
			Category: domain.FAQCategoryTapeAndApp,
		},

		// SOBRE EL FUNCIONAMIENTO DEL APLICATIVO
		{
			Question: "¿Necesito tener internet para usar el app?",
			Answer:   "No. El aplicativo está diseñado para funcionar sin conexión a internet. Solo necesitas tener cargado el teléfono. Algunas funciones como el mapa de centros de salud pueden necesitar GPS o conexión si no están precargadas.",
			Category: domain.FAQCategoryAppInfo,
		},
		{
			Question: "¿Puedo usar el app con cualquier niño/a?",
			Answer:   "Sí, siempre que tenga entre 6 y 59 meses (de 0.5 a 5 años de edad). No se recomienda para bebés menores de 6 meses ni para niños mayores de 5 años.",
			Category: domain.FAQCategoryAppInfo,
		},

		{
			Question: "¿Qué pasa si me equivoco al ingresar el número?",
			Answer:   "Puedes volver atrás y corregir la medición. El app solo guarda la última medición ingresada, así que puedes repetirla si es necesario.",
			Category: domain.FAQCategoryAppInfo,
		},

		// SOBRE LOS RESULTADOS Y LO QUE DEBO HACER
		{
			Question: "¿Qué significa si aparece Alerta Roja?",
			Answer:   "Significa que tu niño/a podría estar con desnutrición severa. Es muy importante que lo lleves al centro de salud lo antes posible, aunque parezca que está bien. La desnutrición no siempre se nota de inmediato.",
			Category: domain.FAQCategoryResults,
		},
		{
			Question: "¿Y si me sale Alerta Amarilla?",
			Answer:   "Significa que hay riesgo de desnutrición. No es una emergencia, pero sí una señal de cuidado. Revisa su alimentación, y llévalo al centro de salud para un chequeo. Puedes volver a medir en 7 días.",
			Category: domain.FAQCategoryResults,
		},
		{
			Question: "¿Y si sale Zona Verde? ¿Todo está bien?",
			Answer:   "Sí, es una buena señal. Pero igual debes seguir con sus controles en el centro de salud y alimentarlo bien. Puedes repetir la medición una vez al mes o si lo ves enfermo o sin apetito.",
			Category: domain.FAQCategoryResults,
		},

		// SOBRE LOS CENTROS DE SALUD Y EL APOYO LOCAL
		{
			Question: "¿Cómo encuentro el centro de salud más cercano?",
			Answer:   "El app puede mostrarte un listado o un mapa, usando GPS si está disponible. Si no tienes conexión, verás una lista precargada con los puestos más cercanos según tu comunidad.",
			Category: domain.FAQCategoryHealthCenters,
		},
		{
			Question: "¿Qué hago si no puedo ir al centro de salud?",
			Answer:   "Busca apoyo del teniente gobernador, promotor de salud o el centro poblado. Ellos pueden ayudarte a comunicarte o trasladarte.",
			Category: domain.FAQCategoryHealthCenters,
		},
		{
			Question: "¿Puedo usar esta app para otros niños de la comunidad?",
			Answer:   "Sí. Puedes usar la cinta y la app con cualquier niño de entre 6 y 59 meses. Solo asegúrate de no confundir las mediciones si lo haces con varios.",
			Category: domain.FAQCategoryHealthCenters,
		},

		// SOBRE PRIVACIDAD Y SEGURIDAD
		{
			Question: "¿El app guarda información personal del niño/a?",
			Answer:   "No. El aplicativo no registra nombres, fotos ni datos personales. Solo guarda las mediciones y los resultados para que puedas consultarlos tú mismo.",
			Category: domain.FAQCategoryPrivacy,
		},
		{
			Question: "¿Quién puede ver los datos que ingreso?",
			Answer:   "Solo tú. Nadie más tiene acceso a tu teléfono ni a lo que registres. Si en el futuro deseas compartir la información con el centro de salud, puedes mostrarla desde tu pantalla.",
			Category: domain.FAQCategoryPrivacy,
		},

		// OTRAS PREGUNTAS
		{
			Question: "¿Este app reemplaza al personal de salud?",
			Answer:   "No. El app es una herramienta de apoyo para el cuidado en casa, pero no reemplaza al centro de salud ni a los profesionales. Siempre debes acudir si tienes dudas o si el niño/a está enfermo.",
			Category: domain.FAQCategoryOther,
		},
		{
			Question: "¿Puedo hacer una nueva medición el mismo día?",
			Answer:   "Sí. Si crees que te equivocaste o si el niño/a comió y crees que cambió, puedes repetir la medición. Lo importante es hacerlo siempre en el mismo brazo y bien colocado.",
			Category: domain.FAQCategoryOther,
		},
	}

	// Crear FAQs con IDs generados
	for i := range faqs {
		faqs[i].ID = uuid.New()
		faqs[i].CreatedAt = time.Now()
	}

	if err := tx.Create(&faqs).Error; err != nil {
		return fmt.Errorf("error creando FAQs: %w", err)
	}

	log.Printf("✅ %d preguntas frecuentes creadas en %d categorías", len(faqs), len(domain.ValidFAQCategories))
	return nil
}

// ============= FUNCIONES DE DATOS ADICIONALES =============

// seedAdditionalData agrega datos faltantes si los roles ya existen
func seedAdditionalData(db *gorm.DB) error {
	log.Println("🔍 Verificando y completando datos del sistema...")

	if err := checkAndCreateTags(db); err != nil {
		return fmt.Errorf("error verificando tags: %w", err)
	}

	if err := checkAndCreateRecommendations(db); err != nil {
		return fmt.Errorf("error verificando recomendaciones: %w", err)
	}

	if err := checkAndCreateFAQs(db); err != nil {
		return fmt.Errorf("error verificando FAQs: %w", err)
	}

	if err := updateExistingData(db); err != nil {
		return fmt.Errorf("error actualizando datos existentes: %w", err)
	}

	log.Println("✅ Verificación de datos completada")
	return nil
}

func checkAndCreateFAQs(db *gorm.DB) error {
	var faqCount int64
	if err := db.Model(&domain.FAQ{}).Count(&faqCount).Error; err != nil {
		return err
	}

	if faqCount == 0 {
		log.Println("❓ No se encontraron FAQs, creando preguntas frecuentes...")
		return seedFAQs(db)
	}

	log.Println("✅ FAQs verificadas - OK")
	return nil
}

// checkAndCreateTags verifica y crea tags faltantes
func checkAndCreateTags(db *gorm.DB) error {
	var tagCount int64
	if err := db.Model(&domain.Tag{}).Count(&tagCount).Error; err != nil {
		return err
	}

	if tagCount == 0 {
		log.Println("🏷️  Creando tags MUAC faltantes...")
		return seedTags(db)
	}

	// Verificar si tienen los campos nuevos
	var tagsWithoutMuacCode int64
	db.Model(&domain.Tag{}).Where("muac_code IS NULL OR muac_code = ''").Count(&tagsWithoutMuacCode)

	if tagsWithoutMuacCode > 0 {
		log.Printf("🔧 Actualizando %d tags con códigos MUAC...", tagsWithoutMuacCode)
		return updateTagsWithMuacCodes(db)
	}

	log.Println("✅ Tags MUAC verificados - OK")
	return nil
}

// checkAndCreateRecommendations verifica y crea recomendaciones faltantes
func checkAndCreateRecommendations(db *gorm.DB) error {
	var recCount int64
	if err := db.Model(&domain.Recommendation{}).Count(&recCount).Error; err != nil {
		return err
	}

	if recCount == 0 {
		log.Println("💡 Creando recomendaciones nutricionales faltantes...")
		return seedRecommendations(db)
	}

	// Verificar si tienen los campos nuevos
	var recsWithoutMuacCode int64
	db.Model(&domain.Recommendation{}).Where("muac_code IS NULL OR muac_code = ''").Count(&recsWithoutMuacCode)

	if recsWithoutMuacCode > 0 {
		log.Printf("🔧 Actualizando %d recomendaciones con códigos MUAC...", recsWithoutMuacCode)
		return updateRecommendationsWithMuacCodes(db)
	}

	log.Println("✅ Recomendaciones verificadas - OK")
	return nil
}

// updateExistingData actualiza datos existentes con campos nuevos
func updateExistingData(db *gorm.DB) error {
	// Activar tags que puedan estar inactivos
	if err := db.Model(&domain.Tag{}).Where("active IS NULL").Update("active", true).Error; err != nil {
		log.Printf("Warning: Error activando tags: %v", err)
	}

	// Activar recomendaciones que puedan estar inactivas
	if err := db.Model(&domain.Recommendation{}).Where("active IS NULL").Update("active", true).Error; err != nil {
		log.Printf("Warning: Error activando recomendaciones: %v", err)
	}

	return nil
}

// updateTagsWithMuacCodes actualiza tags existentes con códigos MUAC
func updateTagsWithMuacCodes(db *gorm.DB) error {
	updates := map[string]map[string]interface{}{
		"MUAC-R1": {
			"muac_code": domain.MuacCodeRed,
			"color":     domain.ColorRed,
			"priority":  domain.PriorityExtreme,
		},
		"MUAC-Y1": {
			"muac_code": domain.MuacCodeYellow,
			"color":     domain.ColorYellow,
			"priority":  domain.PriorityHigh,
		},
		"MUAC-G1": {
			"muac_code": domain.MuacCodeGreen,
			"color":     domain.ColorGreen,
			"priority":  domain.PriorityLow,
		},
		"SEGUIMIENTO": {
			"muac_code": domain.MuacCodeFollow,
			"color":     domain.ColorBlue,
			"priority":  domain.PriorityMedium,
		},
	}

	for name, fields := range updates {
		if err := db.Model(&domain.Tag{}).Where("name = ?", name).Updates(fields).Error; err != nil {
			log.Printf("Warning: Error actualizando tag %s: %v", name, err)
		}
	}

	return nil
}

// updateRecommendationsWithMuacCodes actualiza recomendaciones existentes
func updateRecommendationsWithMuacCodes(db *gorm.DB) error {
	// Buscar recomendaciones por patrones en el nombre
	updates := []struct {
		pattern string
		fields  map[string]interface{}
	}{
		{
			pattern: "%ALERTA ROJA%",
			fields: map[string]interface{}{
				"muac_code":  domain.MuacCodeRed,
				"color_code": domain.ColorRed,
				"priority":   domain.PriorityUrgent,
			},
		},
		{
			pattern: "%ALERTA AMARILLA%",
			fields: map[string]interface{}{
				"muac_code":  domain.MuacCodeYellow,
				"color_code": domain.ColorYellow,
				"priority":   domain.PriorityAttention,
			},
		},
		{
			pattern: "%ZONA VERDE%",
			fields: map[string]interface{}{
				"muac_code":  domain.MuacCodeGreen,
				"color_code": domain.ColorGreen,
				"priority":   domain.PriorityNormal,
			},
		},
		{
			pattern: "%Seguimiento%",
			fields: map[string]interface{}{
				"muac_code":  domain.MuacCodeFollow,
				"color_code": domain.ColorBlue,
				"priority":   domain.PriorityAttention,
			},
		},
	}

	for _, update := range updates {
		if err := db.Model(&domain.Recommendation{}).Where("name LIKE ?", update.pattern).Updates(update.fields).Error; err != nil {
			log.Printf("Warning: Error actualizando recomendaciones con patrón %s: %v", update.pattern, err)
		}
	}

	return nil
}

// ============= FUNCIONES DE LOGGING =============

// logSeedingSummary muestra un resumen de la siembra
func logSeedingSummary(db *gorm.DB) {
	var counts struct {
		Users           int64
		Roles           int64
		Tags            int64
		Recommendations int64
		Patients        int64
		Measurements    int64
		FAQs            int64
	}

	db.Model(&domain.User{}).Count(&counts.Users)
	db.Model(&domain.Role{}).Count(&counts.Roles)
	db.Model(&domain.Tag{}).Count(&counts.Tags)
	db.Model(&domain.Recommendation{}).Count(&counts.Recommendations)
	db.Model(&domain.Patient{}).Count(&counts.Patients)
	db.Model(&domain.Measurement{}).Count(&counts.Measurements)
	db.Model(&domain.FAQ{}).Count(&counts.FAQs)

	log.Println("")
	log.Println("🎉 ============= SISTEMA MUAC INICIALIZADO =============")
	log.Println("📊 Resumen de datos:")
	log.Printf("   👥 Usuarios: %d", counts.Users)
	log.Printf("   🔐 Roles: %d", counts.Roles)
	log.Printf("   🏷️  Tags MUAC: %d", counts.Tags)
	log.Printf("   💡 Recomendaciones: %d", counts.Recommendations)
	log.Printf("   ❓ Preguntas frecuentes: %d", counts.FAQs)
	log.Printf("   🧒 Pacientes: %d", counts.Patients)
	log.Printf("   📏 Mediciones: %d", counts.Measurements)
	log.Println("")
	log.Println("🌍 Clasificación MUAC según estándares OMS/UNICEF/Sphere:")
	log.Printf("   🔴 MUAC-R1: < %.1f cm  (Desnutrición Aguda Severa - SAM)", domain.MuacThresholdSevere)
	log.Printf("   🟡 MUAC-Y1: %.1f-%.1f cm (Desnutrición Aguda Moderada - MAM)", domain.MuacThresholdSevere, domain.MuacThresholdModerate)
	log.Printf("   🟢 MUAC-G1: ≥ %.1f cm  (Estado Nutricional Adecuado)", domain.MuacThresholdNormal)
	log.Println("   🔵 SEGUIMIENTO: Pacientes en recuperación")
	log.Println("")
	log.Println("🔑 Credenciales administrador:")
	log.Println("   📧 Email: admin@muac.org")
	log.Println("   🔒 Password: admin123")
	log.Println("")
	log.Println("✅ Sistema listo para registro de mediciones con")
	log.Println("   clasificación automática y recomendaciones contextualizadas")
	log.Println("=========================================================")
}

// ============= FUNCIONES DE UTILIDAD =============

// GetSeedingStatus retorna el estado actual de la siembra
func GetSeedingStatus(db *gorm.DB) map[string]interface{} {
	var counts struct {
		Users           int64
		Roles           int64
		Tags            int64
		Recommendations int64
	}

	db.Model(&domain.User{}).Count(&counts.Users)
	db.Model(&domain.Role{}).Count(&counts.Roles)
	db.Model(&domain.Tag{}).Count(&counts.Tags)
	db.Model(&domain.Recommendation{}).Count(&counts.Recommendations)

	return map[string]interface{}{
		"seeded":          counts.Roles > 0,
		"users":           counts.Users,
		"roles":           counts.Roles,
		"tags":            counts.Tags,
		"recommendations": counts.Recommendations,
		"muac_ready":      counts.Tags >= 3 && counts.Recommendations >= 3,
		"has_admin":       counts.Users > 0,
	}
}

// ValidateSeedData valida que los datos sembrados sean correctos
func ValidateSeedData(db *gorm.DB) error {
	// Verificar roles esenciales
	requiredRoles := []string{"ADMINISTRADOR", "SUPERVISOR", "APODERADO"}
	for _, roleName := range requiredRoles {
		var count int64
		if err := db.Model(&domain.Role{}).Where("name = ?", roleName).Count(&count).Error; err != nil {
			return fmt.Errorf("error verificando rol %s: %w", roleName, err)
		}
		if count == 0 {
			return fmt.Errorf("rol requerido '%s' no encontrado", roleName)
		}
	}

	// Verificar tags MUAC esenciales
	requiredTags := []string{domain.MuacCodeRed, domain.MuacCodeYellow, domain.MuacCodeGreen}
	for _, tagCode := range requiredTags {
		var count int64
		if err := db.Model(&domain.Tag{}).Where("muac_code = ?", tagCode).Count(&count).Error; err != nil {
			return fmt.Errorf("error verificando tag %s: %w", tagCode, err)
		}
		if count == 0 {
			return fmt.Errorf("tag MUAC requerido '%s' no encontrado", tagCode)
		}
	}

	// Verificar recomendaciones MUAC esenciales
	for _, muacCode := range requiredTags {
		var count int64
		if err := db.Model(&domain.Recommendation{}).Where("muac_code = ?", muacCode).Count(&count).Error; err != nil {
			return fmt.Errorf("error verificando recomendación %s: %w", muacCode, err)
		}
		if count == 0 {
			return fmt.Errorf("recomendación MUAC requerida '%s' no encontrada", muacCode)
		}
	}

	// Verificar que existe al menos un admin
	var adminCount int64
	if err := db.Model(&domain.User{}).
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.name = ?", "ADMINISTRADOR").
		Count(&adminCount).Error; err != nil {
		return fmt.Errorf("error verificando usuarios admin: %w", err)
	}
	if adminCount == 0 {
		return fmt.Errorf("no se encontró ningún usuario administrador")
	}

	return nil
}

// CleanSeedData limpia todos los datos sembrados (útil para testing)
func CleanSeedData(db *gorm.DB) error {
	log.Println("🧹 Limpiando datos sembrados...")

	// Orden inverso por dependencias
	tables := []string{
		"measurements",
		"patients",
		"users",
		"recommendations",
		"tags",
		"roles",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			log.Printf("Warning: Error limpiando tabla %s: %v", table, err)
		}
	}

	log.Println("✅ Datos limpiados")
	return nil
}
