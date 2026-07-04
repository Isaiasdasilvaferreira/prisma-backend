package database

import (
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/config"
	"github.com/nedpals/supabase-go"
)

type Database struct {
	Supabase *supabase.Client
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	supabaseClient := supabase.CreateClient(cfg.SupabaseURL, cfg.SupabaseAnonKey)

	return &Database{
		Supabase: supabaseClient,
	}, nil
}
