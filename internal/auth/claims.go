package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type SupabaseClaims struct {
	jwt.RegisteredClaims
	Email       string                 `json:"email"`
	UserID      string                 `json:"sub"`
	AppMetadata map[string]interface{} `json:"app_metadata"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
	Role        string                 `json:"role"`
	Audience    string                 `json:"aud"`
}

func (c *SupabaseClaims) IsAuthenticated() bool {
	return c.UserID != ""
}

func (c *SupabaseClaims) GetUserID() string {
	return c.UserID
}