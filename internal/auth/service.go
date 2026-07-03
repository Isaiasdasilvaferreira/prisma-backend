package auth

import (
	"context"
)

type AuthService interface {
	SignIn(ctx context.Context, email, password string) (*SupabaseClaims, string, error)
	SignUp(ctx context.Context, email, password string, metadata map[string]interface{}) (*SupabaseClaims, string, error)
	VerifyToken(ctx context.Context, tokenString string) (*SupabaseClaims, error)
	GetUser(ctx context.Context, tokenString string) (*SupabaseClaims, error)
}

type authService struct {
	supabase *SupabaseAuth
}

func NewAuthService(supabaseURL, anonKey, jwtSecret string) AuthService {
	return &authService{
		supabase: NewSupabaseAuth(supabaseURL, anonKey, jwtSecret),
	}
}

func (s *authService) SignIn(ctx context.Context, email, password string) (*SupabaseClaims, string, error) {
	return s.supabase.SignIn(ctx, email, password)
}

func (s *authService) SignUp(ctx context.Context, email, password string, metadata map[string]interface{}) (*SupabaseClaims, string, error) {
	return s.supabase.SignUp(ctx, email, password, metadata)
}

func (s *authService) VerifyToken(ctx context.Context, tokenString string) (*SupabaseClaims, error) {
	return s.supabase.VerifyToken(tokenString)
}

func (s *authService) GetUser(ctx context.Context, tokenString string) (*SupabaseClaims, error) {
	return s.supabase.GetUser(ctx, tokenString)
}
