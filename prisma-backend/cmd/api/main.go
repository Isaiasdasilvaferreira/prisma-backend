package main

import (
	"fmt"
	"log"
	"net/http"

	"prisma-backend/internal/config"
	"prisma-backend/internal/middleware"
	"prisma-backend/routes"
)

func main() {
	cfg := config.LoadConfig()

	authRoutes := routes.NewAuthRoutes(cfg)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth/login", authRoutes.LoginHandler)
	mux.HandleFunc("/api/auth/signup", authRoutes.SignupHandler)

	authMiddleware := middleware.NewAuthMiddleware(authRoutes.AuthService())
	
	mux.HandleFunc("/api/auth/me", authMiddleware.Authenticate(authRoutes.MeHandler))

	handler := middleware.CORSMiddleware(mux)

	port := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal(err)
	}
}