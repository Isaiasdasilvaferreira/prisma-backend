package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
)

type AuthMiddleware struct {
	authService auth.AuthService
}

func NewAuthMiddleware(authService auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]

		claims, err := m.authService.VerifyToken(r.Context(), tokenString)
		if err != nil {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), "user_claims", claims)
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_role", claims.Role)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

func (m *AuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				claims, err := m.authService.VerifyToken(r.Context(), tokenString)
				if err == nil {
					ctx := context.WithValue(r.Context(), "user_claims", claims)
					ctx = context.WithValue(ctx, "user_id", claims.UserID)
					ctx = context.WithValue(ctx, "user_role", claims.Role)
					r = r.WithContext(ctx)
				}
			}
		}

		next(w, r)
	}
}

func (m *AuthMiddleware) RequireRole(role string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("user_claims").(*auth.SupabaseClaims)
			if !ok {
				utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
				return
			}

			if claims.Role != role {
				utils.ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			next(w, r)
		}
	}
}

func (m *AuthMiddleware) AuthenticateAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]

		claims, err := m.authService.VerifyToken(r.Context(), tokenString)
		if err != nil {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		if claims.Role != "admin" {
			utils.ErrorResponse(w, http.StatusForbidden, "Admin access required")
			return
		}

		ctx := context.WithValue(r.Context(), "user_claims", claims)
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_role", claims.Role)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
