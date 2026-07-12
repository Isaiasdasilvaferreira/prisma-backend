package middleware

import (
	"context"
	"net/http"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
)

type AuthMiddleware struct {
	authService *auth.SupabaseAuth
}

func NewAuthMiddleware(authService *auth.SupabaseAuth) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Missing authentication token")
			return
		}

		tokenString := cookie.Value
		if tokenString == "" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

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
		cookie, err := r.Cookie("token")
		if err != nil {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Missing authentication token")
			return
		}

		tokenString := cookie.Value
		if tokenString == "" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

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
