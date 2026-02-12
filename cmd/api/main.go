package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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
		if os.Getenv("ENV") == "production" {
			log.Fatal("JWT_SECRET environment variable is required in production")
		}
		log.Println("⚠️  WARNING: Using default JWT_SECRET for development only")
		jwtSecret = "my_secret_key"
	}

	userRepo := repository.NewPostgresUserRepository()
	authService := service.NewAuthService(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)

	// Init Lixi Dependencies
	lixiRepo := repository.NewPostgresLixiRepository()
	greetingRepo := repository.NewPostgresLixiGreetingRepository()
	lixiService := service.NewLixiService(lixiRepo, greetingRepo)
	lixiHandler := handler.NewLixiHandler(lixiService)

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
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to My Backend API"))
	})

	// Lixi Routes - Public
	mux.HandleFunc("GET /api/lixi/active", lixiHandler.GetActive)
	mux.HandleFunc("POST /api/lixi/greeting", lixiHandler.SubmitGreeting)

	// Lixi Routes - Admin
	mux.HandleFunc("GET /api/admin/lixi", lixiHandler.GetAll)
	mux.HandleFunc("POST /api/admin/lixi", lixiHandler.Create)
	mux.HandleFunc("PUT /api/admin/lixi/{id}", lixiHandler.Update)
	mux.HandleFunc("DELETE /api/admin/lixi/{id}", lixiHandler.Delete)
	mux.HandleFunc("POST /api/admin/lixi/{id}/activate", lixiHandler.Activate)
	mux.HandleFunc("GET /api/admin/lixi/greetings", lixiHandler.GetAllGreetings)

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
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		if os.Getenv("ENV") == "production" {
			log.Fatal("ALLOWED_ORIGINS must be set in production")
		}
		allowedOrigins = "http://localhost:3000,http://localhost:5000"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		originAllowed := false
		for _, allowed := range strings.Split(allowedOrigins, ",") {
			if origin == strings.TrimSpace(allowed) {
				originAllowed = true
				break
			}
		}

		// Set CORS headers if origin is allowed
		if originAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			if originAllowed {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
