package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// ApplyMiddlewares aplica todos los middlewares necesarios
func ApplyMiddlewares(handler http.Handler) http.Handler {
	// Middleware de logging
	handler = LoggingMiddleware(handler)

	// Middleware CORS
	handler = CorsMiddleware(handler)

	// Middleware de recuperación de pánico
	handler = RecoveryMiddleware(handler)

	return handler
}

// LoggingMiddleware registra información sobre cada solicitud
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Llamar al siguiente handler
		next.ServeHTTP(w, r)

		// Registrar la solicitud después de procesarla
		duration := time.Since(start)

		// Log de la solicitud
		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			duration,
		)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configurar cabeceras CORS
		w.Header().Set("Access-Control-Allow-Origin", "*") // o "*" para desarrollo
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 horas

		// Manejar solicitudes OPTIONS (preflight)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Procesar la solicitud
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware recupera de pánicos y devuelve un error 500
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Registrar el pánico
				log.Printf("Pánico recuperado: %v", err)
				log.Printf("Stack trace: %s", debug.Stack())

				// Devolver error 500
				http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
