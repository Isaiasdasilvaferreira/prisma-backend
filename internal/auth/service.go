package auth

import (
	"context"
)

type AuthService interface {
	SignIn(ctx context.Context, email, password string) (*SupabaseClaims, string, error)
	SignUp(ctx context.Context, email, password string, metadata map[string]interface{}) (*SupabaseClaims, string, error)
	VerifyToken(tokenString string) (*SupabaseClaims, error)
	GetUser(ctx context.Context, tokenString string) (*SupabaseClaims, error)
}

type authService struct {
	supabaseAuth *SupabaseAuth
}

func NewAuthService(supabaseAuth *SupabaseAuth) AuthService {
	return &authService{
		supabaseAuth: supabaseAuth,
	}
}

func (s *authService) SignIn(ctx context.Context, email, password string) (*SupabaseClaims, string, error) {
	return s.supabaseAuth.SignIn(ctx, email, password)
}

func (s *authService) SignUp(ctx context.Context, email, password string, metadata map[string]interface{}) (*SupabaseClaims, string, error) {
	return s.supabaseAuth.SignUp(ctx, email, password, metadata)
}

func (s *authService) VerifyToken(tokenString string) (*SupabaseClaims, error) {
	return s.supabaseAuth.VerifyToken(tokenString)
}

func (s *authService) GetUser(ctx context.Context, tokenString string) (*SupabaseClaims, error) {
	return s.supabaseAuth.GetUser(ctx, tokenString)
}
