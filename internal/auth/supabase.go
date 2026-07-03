package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type SupabaseAuth struct {
	url         string
	anonKey     string
	jwtSecret   string
	httpClient  *http.Client
}

func NewSupabaseAuth(url, anonKey, jwtSecret string) *SupabaseAuth {
	return &SupabaseAuth{
		url:         url,
		anonKey:     anonKey,
		jwtSecret:   jwtSecret,
		httpClient:  &http.Client{},
	}
}

func (s *SupabaseAuth) SignIn(ctx context.Context, email, password string) (*SupabaseClaims, string, error) {
	reqBody := map[string]string{
		"email":    email,
		"password": password,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", 
		fmt.Sprintf("%s/auth/v1/token?grant_type=password", s.url),
		strings.NewReader(string(reqJSON)))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", s.anonKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("supabase auth error: %s", string(body))
	}

	var response struct {
		AccessToken string `json:"access_token"`
		User        struct {
			ID       string                 `json:"id"`
			Email    string                 `json:"email"`
			AppMeta  map[string]interface{} `json:"app_metadata"`
			UserMeta map[string]interface{} `json:"user_metadata"`
			Role     string                 `json:"role"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, "", fmt.Errorf("failed to decode response: %w", err)
	}

	claims, err := s.ParseToken(response.AccessToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse token: %w", err)
	}

	return claims, response.AccessToken, nil
}

func (s *SupabaseAuth) SignUp(ctx context.Context, email, password string, metadata map[string]interface{}) (*SupabaseClaims, string, error) {
	reqBody := map[string]interface{}{
		"email":    email,
		"password": password,
		"data":     metadata,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", 
		fmt.Sprintf("%s/auth/v1/signup", s.url),
		strings.NewReader(string(reqJSON)))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", s.anonKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("supabase signup error: %s", string(body))
	}

	var response struct {
		AccessToken string `json:"access_token"`
		User        struct {
			ID       string                 `json:"id"`
			Email    string                 `json:"email"`
			AppMeta  map[string]interface{} `json:"app_metadata"`
			UserMeta map[string]interface{} `json:"user_metadata"`
			Role     string                 `json:"role"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, "", fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("SignUp - AccessToken: '%s'\n", response.AccessToken)
	fmt.Printf("SignUp - UserID: %s\n", response.User.ID)
	fmt.Printf("SignUp - Email: %s\n", response.User.Email)

	claims, err := s.ParseToken(response.AccessToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse token: %w", err)
	}

	return claims, response.AccessToken, nil
}

func (s *SupabaseAuth) ParseToken(tokenString string) (*SupabaseClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SupabaseClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(s.jwtSecret), nil
		}
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); ok {
			jwksJSON := os.Getenv("SUPABASE_JWT_PUBLIC_KEY")
			if jwksJSON == "" {
				return nil, fmt.Errorf("SUPABASE_JWT_PUBLIC_KEY environment variable not set")
			}

			var jwks struct {
				Keys []struct {
					Kid string `json:"kid"`
					Kty string `json:"kty"`
					Crv string `json:"crv"`
					X   string `json:"x"`
					Y   string `json:"y"`
				} `json:"keys"`
			}

			if err := json.Unmarshal([]byte(jwksJSON), &jwks); err != nil {
				return nil, fmt.Errorf("failed to parse JWKS: %w", err)
			}

			if len(jwks.Keys) == 0 {
				return nil, fmt.Errorf("no keys found in JWKS")
			}

			decodedX, err := base64.RawURLEncoding.DecodeString(jwks.Keys[0].X)
			if err != nil {
				return nil, fmt.Errorf("failed to decode X coordinate: %w", err)
			}
			decodedY, err := base64.RawURLEncoding.DecodeString(jwks.Keys[0].Y)
			if err != nil {
				return nil, fmt.Errorf("failed to decode Y coordinate: %w", err)
			}

			publicKey := &ecdsa.PublicKey{
				Curve: elliptic.P256(),
				X:     new(big.Int).SetBytes(decodedX),
				Y:     new(big.Int).SetBytes(decodedY),
			}

			return publicKey, nil
		}
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*SupabaseClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *SupabaseAuth) VerifyToken(tokenString string) (*SupabaseClaims, error) {
	return s.ParseToken(tokenString)
}

func (s *SupabaseAuth) GetUser(ctx context.Context, tokenString string) (*SupabaseClaims, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", 
		fmt.Sprintf("%s/auth/v1/user", s.url), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+tokenString)
	req.Header.Set("apikey", s.anonKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user: status %d", resp.StatusCode)
	}

	var userResponse struct {
		ID       string                 `json:"id"`
		Email    string                 `json:"email"`
		AppMeta  map[string]interface{} `json:"app_metadata"`
		UserMeta map[string]interface{} `json:"user_metadata"`
		Role     string                 `json:"role"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &SupabaseClaims{
		UserID:      userResponse.ID,
		Email:       userResponse.Email,
		AppMetadata: userResponse.AppMeta,
		UserMetadata: userResponse.UserMeta,
		Role:        userResponse.Role,
	}, nil
}
