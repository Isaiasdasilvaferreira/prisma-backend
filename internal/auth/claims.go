package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type SupabaseClaims struct {
	jwt.RegisteredClaims
	Email        string                 `json:"email"`
	UserID       string                 `json:"sub"`
	AppMetadata  map[string]interface{} `json:"app_metadata"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
	Role         string                 `json:"role"`
	Audience     string                 `json:"aud"`
}

func (c *SupabaseClaims) IsAuthenticated() bool {
	return c.UserID != ""
}

func (c *SupabaseClaims) GetUserID() string {
	return c.UserID
}

type contextKey string

const UserIDKey contextKey = "user_id"
const UserClaimsKey contextKey = "user_claims"
const UserRoleKey contextKey = "user_role"

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

func GetUserFromContext(ctx context.Context) (*SupabaseClaims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*SupabaseClaims)
	return claims, ok
}
