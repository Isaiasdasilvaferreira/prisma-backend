package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/config"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/database"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/middleware"
	"github.com/Isaiasdasilvaferreira/prisma-backend/routes"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	authService := auth.NewAuthService(cfg.SupabaseURL, cfg.SupabaseAnonKey, cfg.SupabaseJWTSecret, db.Supabase)

	authRoutes := routes.NewAuthRoutes(cfg, authService, db.Supabase)

	mux := http.NewServeMux()

	authRoutes.RegisterRoutes(mux)

	handler := middleware.CORSMiddleware(mux)

	port := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal(err)
	}
}
