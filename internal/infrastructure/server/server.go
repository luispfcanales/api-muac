package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luispfcanales/api-muac/internal/infrastructure/config"
	"github.com/luispfcanales/api-muac/internal/infrastructure/server/middleware"
)

// Server representa el servidor HTTP
type Server struct {
	server *http.Server
	config *config.Config
}

// NewServer crea una nueva instancia del servidor
func NewServer(config *config.Config, handler http.Handler) *Server {

	handler = middleware.ApplyMiddlewares(handler)

	return &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", config.ServerPort),
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		config: config,
	}
}

// Start inicia el servidor HTTP
func (s *Server) Start() error {
	// Canal para capturar señales del sistema operativo
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Canal para errores del servidor
	errCh := make(chan error, 1)

	// Iniciar el servidor en una goroutine
	go func() {
		log.Printf("Servidor iniciado en http://localhost:%d", s.config.ServerPort)
		errCh <- s.server.ListenAndServe()
	}()

	// Esperar a que ocurra un error o se reciba una señal de parada
	select {
	case <-stop:
		log.Println("Apagando servidor...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("error al apagar el servidor: %w", err)
		}
		log.Println("Servidor apagado correctamente")
		return nil
	case err := <-errCh:
		if err != http.ErrServerClosed {
			return fmt.Errorf("error en el servidor: %w", err)
		}
		return nil
	}
}
