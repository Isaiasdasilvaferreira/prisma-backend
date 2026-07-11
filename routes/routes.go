package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/config"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/middleware"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/opportunity"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/plans"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/scraper"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/user"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
	"github.com/nedpals/supabase-go"
)

type AuthRoutes struct {
	authService          *auth.SupabaseAuth
	authMiddleware       *middleware.AuthMiddleware
	planController       *plans.PlanController
	scraperController    *scraper.ScraperController
	opportunityController *opportunity.Controller
}

func NewAuthRoutes(cfg *config.Config, authService *auth.SupabaseAuth, supabaseClient *supabase.Client) *AuthRoutes {
	authMiddleware := middleware.NewAuthMiddleware(authService)

	planRepo := user.NewPlanRepository(supabaseClient)
	userSvc := user.NewService(planRepo)
	oppRepo := opportunity.NewRepository(supabaseClient)
	
	planService := plans.NewPlanService(planRepo)
	planController := plans.NewPlanController(planService)
	
	oppService := opportunity.NewService(oppRepo, userSvc)
	opportunityController := opportunity.NewController(oppService)
	
	scraperController := scraper.NewScraperController(supabaseClient, userSvc, oppRepo)

	return &AuthRoutes{
		authService:          authService,
		authMiddleware:       authMiddleware,
		planController:       planController,
		scraperController:    scraperController,
		opportunityController: opportunityController,
	}
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

func (r *AuthRoutes) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/login", r.LoginHandler)
	mux.HandleFunc("/api/auth/signup", r.SignupHandler)
	mux.HandleFunc("/api/auth/me", r.authMiddleware.Authenticate(r.MeHandler))

	mux.HandleFunc("/api/plans/can-scrape", r.authMiddleware.Authenticate(r.planController.CanScrape))
	mux.HandleFunc("/api/plans/upgrade", r.authMiddleware.Authenticate(r.planController.UpgradeToProfessional))
	mux.HandleFunc("/api/user/plan", r.authMiddleware.Authenticate(r.planController.GetUserPlan))

	mux.HandleFunc("/api/scrape/ashby", r.authMiddleware.Authenticate(r.scraperController.ScrapeAshby))
	mux.HandleFunc("/api/scrape/greenhouse", r.authMiddleware.Authenticate(r.scraperController.ScrapeGreenhouse))
	mux.HandleFunc("/api/scrape/lever", r.authMiddleware.Authenticate(r.scraperController.ScrapeLever))
	mux.HandleFunc("/api/scrape/all", r.authMiddleware.Authenticate(r.scraperController.ScrapeAll))
	mux.HandleFunc("/api/scraping/trigger", r.authMiddleware.AuthenticateAdmin(r.scraperController.TriggerScraping))

	mux.HandleFunc("/api/opportunities", r.authMiddleware.Authenticate(r.opportunityController.GetUserOpportunities))
	mux.HandleFunc("/api/opportunities/", r.authMiddleware.Authenticate(func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		if strings.Contains(path, "/source/") {
			r.opportunityController.GetOpportunitiesBySource(w, req)
		} else if strings.HasPrefix(path, "/api/opportunities/") && path != "/api/opportunities" {
			r.opportunityController.GetUserOpportunityByID(w, req)
		} else {
			r.opportunityController.GetUserOpportunities(w, req)
		}
	}))
	mux.HandleFunc("/api/opportunities/stats", r.authMiddleware.Authenticate(r.opportunityController.GetOpportunitiesStats))
}
