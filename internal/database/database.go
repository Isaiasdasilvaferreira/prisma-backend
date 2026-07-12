package database

import (
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/config"
	"github.com/nedpals/supabase-go"
)

type Database struct {
	Supabase       *supabase.Client
	SupabaseAdmin  *supabase.Client
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	supabaseClient := supabase.CreateClient(cfg.SupabaseURL, cfg.SupabaseAnonKey)

	var supabaseAdmin *supabase.Client
	if cfg.SupabaseServiceRoleKey != "" {
		supabaseAdmin = supabase.CreateClient(cfg.SupabaseURL, cfg.SupabaseServiceRoleKey)
	}

	return &Database{
		Supabase:      supabaseClient,
		SupabaseAdmin: supabaseAdmin,
	}, nil
}
