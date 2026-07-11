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

		claims, err := m.authService.VerifyToken(tokenString)
		if err != nil {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, auth.UserClaimsKey, claims)
		ctx = context.WithValue(ctx, auth.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, auth.UserRoleKey, claims.Role)
		r = r.WithContext(ctx)

		next(w, r)
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

		claims, err := m.authService.VerifyToken(tokenString)
		if err != nil {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		if claims.Role != "admin" {
			utils.ErrorResponse(w, http.StatusForbidden, "Admin access required")
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, auth.UserClaimsKey, claims)
		ctx = context.WithValue(ctx, auth.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, auth.UserRoleKey, claims.Role)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
