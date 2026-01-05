package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"my_backend/internal/handler"
	"my_backend/internal/repository"
	"my_backend/internal/service"
)

func main() {
	// 1. Init Dependencies
	// In a real app, you would load these from config/env
	jwtSecret := "my_secret_key"

	userRepo := repository.NewMemoryUserRepository()
	authService := service.NewAuthService(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)

	// Seed Admin User
	_, err := authService.Register(context.Background(), "admin", "12345678@X")
	if err != nil {
		fmt.Printf("Error seeding admin user: %v\n", err)
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
