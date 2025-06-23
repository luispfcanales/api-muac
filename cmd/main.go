package main

import (
	"fmt"
	"log"
	stdhttp "net/http"
	"reflect"

	"github.com/luispfcanales/api-muac/docs"
	_ "github.com/luispfcanales/api-muac/docs" // Importa los docs generados
	"github.com/luispfcanales/api-muac/internal/adapters/handlers/http"
	"github.com/luispfcanales/api-muac/internal/adapters/repositories/postgres"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/services"
	"github.com/luispfcanales/api-muac/internal/infrastructure/config"
	"github.com/luispfcanales/api-muac/internal/infrastructure/server"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title API MUAC
// @version 1.0
// @description API para el sistema de medición MUAC
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email tu.email@ejemplo.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	// Cargar configuración
	cfg := config.LoadConfig()

	// Conectar a la base de datos con GORM
	db, err := config.NewGormDBConnection(cfg)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// Lista de modelos a migrar
	modelos := []interface{}{
		&domain.Role{},
		&domain.Locality{},
		&domain.Patient{},
		&domain.Tag{},
		&domain.User{},
		&domain.Recommendation{},
		&domain.Measurement{},
		&domain.Notification{},
		&domain.FAQ{},
	}

	log.Println("Eliminando TODAS las tablas de la base de datos...")

	// Obtener todas las tablas
	var tables []string
	db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables)

	// Eliminar cada tabla
	for _, table := range tables {
		log.Printf("Eliminando tabla: %s", table)
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			log.Printf("Advertencia: Error al eliminar tabla %s: %v", table, err)
		}
	}
	// Migrar cada modelo y registrar en el log
	log.Println("Iniciando migración de modelos...")
	for _, modelo := range modelos {
		nombreModelo := reflect.TypeOf(modelo).Elem().Name()
		log.Printf("Migrando modelo: %s", nombreModelo)
		if err := db.AutoMigrate(modelo); err != nil {
			log.Fatalf("Error al migrar modelo %s: %v", nombreModelo, err)
		}
		log.Printf("Modelo %s migrado exitosamente", nombreModelo)
	}
	log.Println("Migración completada exitosamente")

	// Sembrar datos iniciales
	if err := config.SeedDatabase(db); err != nil {
		log.Fatalf("Error al sembrar datos iniciales: %v", err)
	}
	// Crear repositorios
	roleRepo := postgres.NewRoleRepository(db)
	userRepo := postgres.NewUserRepository(db)
	notificationRepo := postgres.NewNotificationRepository(db)
	faqRepo := postgres.NewFAQRepository(db)
	localityRepo := postgres.NewLocalityRepository(db)
	recommendationRepo := postgres.NewRecommendationRepository(db)
	tagRepo := postgres.NewTagRepository(db)
	measurementRepo := postgres.NewMeasurementRepository(db)
	patientRepo := postgres.NewPatientRepository(db)

	// Crear servicios
	roleService := services.NewRoleService(roleRepo)
	userService := services.NewUserService(userRepo, roleRepo)
	notificationService := services.NewNotificationService(notificationRepo)
	faqService := services.NewFAQService(faqRepo)
	localityService := services.NewLocalityService(localityRepo)
	recommendationService := services.NewRecommendationService(recommendationRepo)
	tagService := services.NewTagService(tagRepo)
	measurementService := services.NewMeasurementService(measurementRepo, tagRepo, recommendationRepo)
	patientService := services.NewPatientService(patientRepo, measurementRepo)

	// Crear manejadores HTTP
	roleHandler := http.NewRoleHandler(roleService)
	userHandler := http.NewUserHandler(userService)
	notificationHandler := http.NewNotificationHandler(notificationService)
	faqHandler := http.NewFAQHandler(faqService)
	localityHandler := http.NewLocalityHandler(localityService)
	recommendationHandler := http.NewRecommendationHandler(recommendationService)
	tagHandler := http.NewTagHandler(tagService)
	measurementHandler := http.NewMeasurementHandler(measurementService)
	patientHandler := http.NewPatientHandler(patientService)

	// Configurar rutas
	mux := stdhttp.NewServeMux()

	// Servir el archivo swagger.json directamente
	mux.HandleFunc("GET /swagger/doc.json", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Usar el JSON ya procesado en lugar de la plantilla
		w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
	})

	// Agregar documentación Swagger - Modificar esta parte
	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	roleHandler.RegisterRoutes(mux)
	userHandler.RegisterRoutes(mux)
	notificationHandler.RegisterRoutes(mux)
	faqHandler.RegisterRoutes(mux)
	localityHandler.RegisterRoutes(mux)
	recommendationHandler.RegisterRoutes(mux)
	tagHandler.RegisterRoutes(mux)
	measurementHandler.RegisterRoutes(mux)
	patientHandler.RegisterRoutes(mux)

	// Crear y iniciar servidor
	srv := server.NewServer(cfg, mux)
	if err := srv.Start(); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
