package main

// Точка входа HTTP-сервера (Go).
//
// Здесь мы:
// - собираем зависимости (Repository → UseCase → Delivery)
// - подключаем middleware Basic Auth к защищённым роутам
// - поднимаем HTTP-сервер и делаем graceful shutdown
import (
	"basic-auth/internal/delivery"
	"basic-auth/internal/delivery/middleware"
	"basic-auth/internal/repository/memory"
	"basic-auth/internal/usecase"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

// healthHandler — простой публичный endpoint для проверки, что сервер жив.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// swaggerUIHandler отдаёт простую Swagger UI страницу (без зависимостей),
// которая берёт OpenAPI JSON из /openapi.json.
func swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = () => {
        SwaggerUIBundle({ url: '/openapi.json', dom_id: '#swagger-ui' });
      };
    </script>
  </body>
</html>`)
}

// openAPIHandler отдаёт OpenAPI спецификацию (openapi.json) из директории проекта.
func openAPIHandler(w http.ResponseWriter, r *http.Request) {
	// Файл лежит рядом с go.mod в корне реализации: 01-basic-auth/go/openapi.json
	exe, err := os.Executable()
	if err != nil {
		http.Error(w, "failed to locate executable", http.StatusInternalServerError)
		return
	}
	// При go run executable лежит во временной папке, поэтому идём от текущей директории процесса.
	wd, err := os.Getwd()
	if err != nil {
		wd = filepath.Dir(exe)
	}
	specPath := filepath.Join(wd, "openapi.json")
	data, err := os.ReadFile(specPath)
	if err != nil {
		http.Error(w, "openapi.json not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func main() {
	// 1) Data layer (in-memory хранилище пользователей)
	userRepository := memory.NewUserRepository()

	// 2) Business logic (use cases авторизации)
	authUsecase := usecase.NewAuthUsecase(userRepository)

	// 3) HTTP handlers (Delivery слой)
	authHandler := delivery.NewAuthHandler(authUsecase)

	// 4) Middleware, которое будет проверять заголовок Authorization: Basic ...
	basicAuthMiddleware := middleware.BasicAuth(authUsecase)

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/swagger", swaggerUIHandler)
	mux.HandleFunc("/openapi.json", openAPIHandler)
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.RegisterHandler)

	// Protected routes (require Basic Auth header)
	mux.Handle("GET /api/v1/auth/me", basicAuthMiddleware(http.HandlerFunc(authHandler.MeHandler)))
	mux.Handle("DELETE /api/v1/auth/me", basicAuthMiddleware(http.HandlerFunc(authHandler.DeleteUserHandler)))

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      mux,
	}

	go func() {
		log.Println("Basic Auth server running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
