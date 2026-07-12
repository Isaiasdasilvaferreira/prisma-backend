package middleware

import (
	"net/http"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
)

func CSRFMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" || r.Method == "PATCH" {
			if !auth.ValidateCSRFToken(r) {
				utils.ErrorResponse(w, http.StatusForbidden, "Invalid CSRF token")
				return
			}
		}
		next(w, r)
	}
}
