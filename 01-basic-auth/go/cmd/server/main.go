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
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// healthHandler — простой публичный endpoint для проверки, что сервер жив.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
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
