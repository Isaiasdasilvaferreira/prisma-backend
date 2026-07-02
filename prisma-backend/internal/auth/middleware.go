package auth

import (
	"context"
)

func GetUserFromContext(ctx context.Context) (*SupabaseClaims, bool) {
	claims, ok := ctx.Value("user_claims").(*SupabaseClaims)
	return claims, ok
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	claims, ok := GetUserFromContext(ctx)
	if !ok {
		return "", false
	}
	return claims.GetUserID(), true
}
