package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// SignIn realiza login com email e senha
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

	// Parsear o token JWT para obter os claims
	claims, err := s.ParseToken(response.AccessToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse token: %w", err)
	}

	return claims, response.AccessToken, nil
}

// SignUp realiza cadastro de usuário
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

	claims, err := s.ParseToken(response.AccessToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse token: %w", err)
	}

	return claims, response.AccessToken, nil
}

// ParseToken valida e parseia o JWT token
func (s *SupabaseAuth) ParseToken(tokenString string) (*SupabaseClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SupabaseClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
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