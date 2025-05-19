package main

import (
	"log"
	stdhttp "net/http"

	"github.com/luispfcanales/api-muac/internal/adapters/handlers/http"
	"github.com/luispfcanales/api-muac/internal/adapters/repositories/postgres"
	"github.com/luispfcanales/api-muac/internal/core/services"
	"github.com/luispfcanales/api-muac/internal/infrastructure/config"
	"github.com/luispfcanales/api-muac/internal/infrastructure/server"
)

func main() {
	// Cargar configuración
	cfg := config.LoadConfig()

	// Conectar a la base de datos
	db, err := config.NewDBConnection(cfg)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Crear repositorios
	roleRepo := postgres.NewRoleRepository(db)

	// Crear servicios
	roleService := services.NewRoleService(roleRepo)

	// Crear manejadores HTTP
	roleHandler := http.NewRoleHandler(roleService)

	// Configurar rutas
	mux := stdhttp.NewServeMux()
	roleHandler.RegisterRoutes(mux)

	// Crear y iniciar servidor
	srv := server.NewServer(cfg, mux)
	if err := srv.Start(); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
