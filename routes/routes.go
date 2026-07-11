package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/config"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/middleware"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/opportunity"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/plans"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/scraper"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/user"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
	"github.com/gorilla/mux"
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
	router := mux.NewRouter()
	
	router.HandleFunc("/api/auth/login", r.LoginHandler).Methods("POST")
	router.HandleFunc("/api/auth/signup", r.SignupHandler).Methods("POST")
	router.HandleFunc("/api/auth/me", r.authMiddleware.Authenticate(r.MeHandler)).Methods("GET")

	router.HandleFunc("/api/plans/can-scrape", r.authMiddleware.Authenticate(r.planController.CanScrape)).Methods("GET")
	router.HandleFunc("/api/plans/upgrade", r.authMiddleware.Authenticate(r.planController.UpgradeToProfessional)).Methods("POST")
	router.HandleFunc("/api/user/plan", r.authMiddleware.Authenticate(r.planController.GetUserPlan)).Methods("GET")

	router.HandleFunc("/api/scrape/ashby", r.authMiddleware.Authenticate(r.scraperController.ScrapeAshby)).Methods("GET")
	router.HandleFunc("/api/scrape/greenhouse", r.authMiddleware.Authenticate(r.scraperController.ScrapeGreenhouse)).Methods("GET")
	router.HandleFunc("/api/scrape/lever", r.authMiddleware.Authenticate(r.scraperController.ScrapeLever)).Methods("GET")
	router.HandleFunc("/api/scrape/all", r.authMiddleware.Authenticate(r.scraperController.ScrapeAll)).Methods("GET")
	router.HandleFunc("/api/scraping/trigger", r.authMiddleware.AuthenticateAdmin(r.scraperController.TriggerScraping)).Methods("POST")

	router.HandleFunc("/api/opportunities", r.authMiddleware.Authenticate(r.opportunityController.GetUserOpportunities)).Methods("GET")
	router.HandleFunc("/api/opportunities/{id}", r.authMiddleware.Authenticate(r.opportunityController.GetUserOpportunityByID)).Methods("GET")
	router.HandleFunc("/api/opportunities/source/{source}", r.authMiddleware.Authenticate(r.opportunityController.GetOpportunitiesBySource)).Methods("GET")
	router.HandleFunc("/api/opportunities/stats", r.authMiddleware.Authenticate(r.opportunityController.GetOpportunitiesStats)).Methods("GET")

	mux.Handle("/", router)
}
