package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/config"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/middleware"
	"github.com/Isaiasdasilvaferreira/prisma-backend/routes"
)

func main() {
	cfg := config.LoadConfig()

	authRoutes := routes.NewAuthRoutes(cfg)

	mux := http.NewServeMux()

	authRoutes.RegisterRoutes(mux)

	handler := middleware.CORSMiddleware(mux)

	port := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal(err)
	}
}
