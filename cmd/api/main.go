package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/config"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/database"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/middleware"
	"github.com/Isaiasdasilvaferreira/prisma-backend/routes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Error().Err(err).Msg("⚠️ Erro ao criar pasta logs")
	}

	logFiles := []string{"logs/error.txt", "logs/info.txt", "logs/data.txt"}
	for _, file := range logFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			if f, err := os.Create(file); err != nil {
				log.Error().Err(err).Msgf("⚠️ Erro ao criar %s", file)
			} else {
				f.Close()
				log.Info().Msgf("✅ Arquivo criado: %s", file)
			}
		}
	}

	cfg := config.LoadConfig()

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	supabaseAuth := auth.NewSupabaseAuth(cfg.SupabaseURL, cfg.SupabaseAnonKey, cfg.SupabaseJWTSecret)

	authRoutes := routes.NewAuthRoutes(cfg, supabaseAuth, db.Supabase, db.SupabaseAdmin)

	mux := http.NewServeMux()

	authRoutes.RegisterRoutes(mux)

	handler := middleware.CORSMiddleware(mux)

	port := fmt.Sprintf(":%s", cfg.Port)
	log.Info().Msgf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
