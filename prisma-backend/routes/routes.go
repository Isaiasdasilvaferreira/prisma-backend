package routes

import (
	"encoding/json"
	"net/http"

	"prisma-backend/internal/auth"
	"prisma-backend/internal/config"
	"prisma-backend/internal/middleware"
	"prisma-backend/internal/utils"
)

type AuthRoutes struct {
	authService    auth.AuthService
	authMiddleware *middleware.AuthMiddleware
}

func NewAuthRoutes(cfg *config.Config) *AuthRoutes {
	authService := auth.NewAuthService(cfg.SupabaseURL, cfg.SupabaseAnonKey, cfg.SupabaseJWTSecret)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	return &AuthRoutes{
		authService:    authService,
		authMiddleware: authMiddleware,
	}
}

func (r *AuthRoutes) AuthService() auth.AuthService {
	return r.authService
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Email    string                 `json:"email"`
	Password string                 `json:"password"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (r *AuthRoutes) LoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	claims, token, err := r.authService.SignIn(req.Context(), loginReq.Email, loginReq.Password)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	response := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":    claims.UserID,
			"email": claims.Email,
			"role":  claims.Role,
		},
	}

	utils.SuccessResponse(w, http.StatusOK, response)
}

func (r *AuthRoutes) SignupHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var signupReq SignupRequest
	if err := json.NewDecoder(req.Body).Decode(&signupReq); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	claims, token, err := r.authService.SignUp(req.Context(), signupReq.Email, signupReq.Password, signupReq.Metadata)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":    claims.UserID,
			"email": claims.Email,
			"role":  claims.Role,
		},
	}

	utils.SuccessResponse(w, http.StatusCreated, response)
}

func (r *AuthRoutes) MeHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	claims, ok := auth.GetUserFromContext(req.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	response := map[string]interface{}{
		"id":    claims.UserID,
		"email": claims.Email,
		"role":  claims.Role,
	}

	utils.SuccessResponse(w, http.StatusOK, response)
}
