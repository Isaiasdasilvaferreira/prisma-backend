package auth

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/user"
)

type AuthService interface {
	SignIn(ctx context.Context, email, password string) (*SupabaseClaims, string, error)
	SignUp(ctx context.Context, email, password string, metadata map[string]interface{}) (*SupabaseClaims, string, error)
	VerifyToken(ctx context.Context, tokenString string) (*SupabaseClaims, error)
	GetUser(ctx context.Context, tokenString string) (*SupabaseClaims, error)
}

type authService struct {
	supabase *SupabaseAuth
	planRepo user.PlanRepository
}

func NewAuthService(supabaseURL, anonKey, jwtSecret string, db *sql.DB) AuthService {
	return &authService{
		supabase: NewSupabaseAuth(supabaseURL, anonKey, jwtSecret),
		planRepo: user.NewPlanRepository(db),
	}
}

func (s *authService) SignIn(ctx context.Context, email, password string) (*SupabaseClaims, string, error) {
	return s.supabase.SignIn(ctx, email, password)
}

func (s *authService) SignUp(ctx context.Context, email, password string, metadata map[string]interface{}) (*SupabaseClaims, string, error) {
	claims, token, err := s.supabase.SignUp(ctx, email, password, metadata)
	if err != nil {
		return nil, "", err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, "", fmt.Errorf("invalid user ID format: %w", err)
	}

	_, err = s.planRepo.CreateUserPlan(ctx, userID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user plan: %w", err)
	}

	return claims, token, nil
}

func (s *authService) VerifyToken(ctx context.Context, tokenString string) (*SupabaseClaims, error) {
	return s.supabase.VerifyToken(tokenString)
}

func (s *authService) GetUser(ctx context.Context, tokenString string) (*SupabaseClaims, error) {
	return s.supabase.GetUser(ctx, tokenString)
}
