package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"my_backend/internal/database"
	"my_backend/internal/handler"
	"my_backend/internal/repository"
	"my_backend/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (ignore error if not found - production uses env vars)
	_ = godotenv.Load()

	// 1. Connect to Database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// 2. Run Migrations
	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 3. Init Dependencies
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my_secret_key" // fallback for development
	}

	userRepo := repository.NewPostgresUserRepository()
	authService := service.NewAuthService(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)

	// Seed Admin User (ignore error if already exists)
	_, err := authService.Register(context.Background(), "admin", "12345678@X")
	if err != nil {
		fmt.Printf("Admin user already exists or error: %v\n", err)
	} else {
		fmt.Println("Admin user created: admin / 12345678@X")
	}

	// 2. Setup Router
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 3. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("Server is running on http://localhost%s\n", addr)

	// Add CORS middleware
	handler := enableCORS(mux)

	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
