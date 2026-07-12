package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func GenerateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func ValidateCSRFToken(r *http.Request) bool {
	tokenFromHeader := r.Header.Get("X-CSRF-Token")
	if tokenFromHeader == "" {
		return false
	}

	cookie, err := r.Cookie("csrf_token")
	if err != nil {
		return false
	}

	return cookie.Value == tokenFromHeader
}
