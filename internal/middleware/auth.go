package middleware

import (
	"context"
	"log"
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
			log.Printf("❌ Cookie 'token' não encontrado: %v", err)
			utils.ErrorResponse(w, http.StatusUnauthorized, "Missing authentication token")
			return
		}

		tokenString := cookie.Value
		if tokenString == "" {
			log.Printf("❌ Token vazio no cookie")
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		log.Printf("✅ Token encontrado no cookie: %s...", tokenString[:20])

		claims, err := m.authService.VerifyToken(tokenString)
		if err != nil {
			log.Printf("❌ Erro ao verificar token: %v", err)
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		log.Printf("✅ Token válido para usuário: %s", claims.UserID)

		ctx := context.WithValue(r.Context(), auth.UserClaimsKey, claims)
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
			log.Printf("❌ Cookie 'token' não encontrado (admin): %v", err)
			utils.ErrorResponse(w, http.StatusUnauthorized, "Missing authentication token")
			return
		}

		tokenString := cookie.Value
		if tokenString == "" {
			log.Printf("❌ Token vazio no cookie (admin)")
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, err := m.authService.VerifyToken(tokenString)
		if err != nil {
			log.Printf("❌ Erro ao verificar token (admin): %v", err)
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		if claims.Role != "admin" {
			log.Printf("❌ Usuário não é admin: %s", claims.Role)
			utils.ErrorResponse(w, http.StatusForbidden, "Admin access required")
			return
		}

		ctx := context.WithValue(r.Context(), auth.UserClaimsKey, claims)
		ctx = context.WithValue(ctx, auth.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, auth.UserRoleKey, claims.Role)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
