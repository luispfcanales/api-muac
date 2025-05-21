package main

import (
	"log"
	stdhttp "net/http"
	"reflect"

	"github.com/luispfcanales/api-muac/internal/adapters/handlers/http"
	"github.com/luispfcanales/api-muac/internal/adapters/repositories/postgres"
	"github.com/luispfcanales/api-muac/internal/core/domain"
	"github.com/luispfcanales/api-muac/internal/core/services"
	"github.com/luispfcanales/api-muac/internal/infrastructure/config"
	"github.com/luispfcanales/api-muac/internal/infrastructure/server"
)

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
		&domain.Father{},
		&domain.Recommendation{},
		&domain.Measurement{},
		&domain.Notification{},
		&domain.FAQ{},
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
	//fatherRepo := postgres.NewFatherRepository(db)
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
