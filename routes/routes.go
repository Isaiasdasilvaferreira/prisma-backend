package routes

import (
	"net/http"
	"strings"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/config"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/middleware"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/opportunity"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/plans"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/scraper"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/user"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/user_opportunity"
	"github.com/nedpals/supabase-go"
)

type AuthRoutes struct {
	authController            *auth.AuthController
	authMiddleware            *middleware.AuthMiddleware
	planController            *plans.PlanController
	scraperController         *scraper.ScraperController
	opportunityController     *opportunity.Controller
	userOpportunityController *user_opportunity.Controller
}

func NewAuthRoutes(cfg *config.Config, authService *auth.SupabaseAuth, supabaseClient *supabase.Client, supabaseAdmin *supabase.Client) *AuthRoutes {
	authMiddleware := middleware.NewAuthMiddleware(authService)
	authController := auth.NewAuthController(authService)

	planRepo := user.NewPlanRepository(supabaseClient, supabaseAdmin)
	userSvc := user.NewService(planRepo)
	oppRepo := opportunity.NewRepository(supabaseClient, supabaseAdmin)

	planService := plans.NewPlanService(planRepo)
	planController := plans.NewPlanController(planService)

	oppService := opportunity.NewService(oppRepo, userSvc)
	opportunityController := opportunity.NewController(oppService)

	scraperController := scraper.NewScraperController(supabaseClient, userSvc, oppRepo)

	userOpportunityRepo := user_opportunity.NewRepository(supabaseClient, supabaseAdmin)
	userOpportunityService := user_opportunity.NewService(userOpportunityRepo)
	userOpportunityController := user_opportunity.NewController(userOpportunityService)

	return &AuthRoutes{
		authController:            authController,
		authMiddleware:            authMiddleware,
		planController:            planController,
		scraperController:         scraperController,
		opportunityController:     opportunityController,
		userOpportunityController: userOpportunityController,
	}
}

func (r *AuthRoutes) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/login", r.authController.Login)
	mux.HandleFunc("/api/auth/signup", r.authController.Signup)
	mux.HandleFunc("/api/auth/logout", r.authMiddleware.Authenticate(middleware.CSRF(r.authController.Logout)))
	mux.HandleFunc("/api/auth/me", r.authMiddleware.Authenticate(r.authController.Me))

	mux.HandleFunc("/api/plans/can-scrape", r.authMiddleware.Authenticate(r.planController.CanScrape))
	mux.HandleFunc("/api/plans/upgrade", r.authMiddleware.Authenticate(middleware.CSRF(r.planController.UpgradeToProfessional)))
	mux.HandleFunc("/api/user/plan", r.authMiddleware.Authenticate(r.planController.GetUserPlan))

	mux.HandleFunc("/api/scrape/ashby", r.authMiddleware.Authenticate(r.scraperController.ScrapeAshby))
	mux.HandleFunc("/api/scrape/greenhouse", r.authMiddleware.Authenticate(r.scraperController.ScrapeGreenhouse))
	mux.HandleFunc("/api/scrape/lever", r.authMiddleware.Authenticate(r.scraperController.ScrapeLever))
	mux.HandleFunc("/api/scrape/all", r.authMiddleware.Authenticate(r.scraperController.ScrapeAll))
	mux.HandleFunc("/api/scraping/trigger", r.authMiddleware.AuthenticateAdmin(middleware.CSRF(r.scraperController.TriggerScraping)))

	mux.HandleFunc("/api/opportunities", r.authMiddleware.Authenticate(r.opportunityController.GetUserOpportunities))
	mux.HandleFunc("/api/opportunities/stats", r.authMiddleware.Authenticate(r.opportunityController.GetOpportunitiesStats))
	mux.HandleFunc("/api/opportunities/", r.authMiddleware.Authenticate(func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		if strings.Contains(path, "/source/") {
			r.opportunityController.GetOpportunitiesBySource(w, req)
		} else if path != "/api/opportunities/" && path != "/api/opportunities" {
			r.opportunityController.GetUserOpportunityByExternalID(w, req)
		} else {
			r.opportunityController.GetUserOpportunities(w, req)
		}
	}))

	mux.HandleFunc("/api/user-opportunities", r.authMiddleware.Authenticate(func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			r.userOpportunityController.CreateUserOpportunity(w, req)
		case http.MethodGet:
			r.userOpportunityController.GetAllUserOpportunities(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/user-opportunities/", r.authMiddleware.Authenticate(func(w http.ResponseWriter, req *http.Request) {
		path := strings.TrimPrefix(req.URL.Path, "/api/user-opportunities/")
		parts := strings.Split(path, "/")

		if len(parts) == 0 || parts[0] == "" {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		if len(parts) == 1 {
			switch req.Method {
			case http.MethodGet:
				r.userOpportunityController.GetUserOpportunity(w, req)
			case http.MethodPut:
				r.userOpportunityController.UpdateUserOpportunity(w, req)
			case http.MethodDelete:
				r.userOpportunityController.DeleteUserOpportunity(w, req)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if len(parts) == 2 {
			switch parts[1] {
			case "approve":
				if req.Method == http.MethodPatch || req.Method == http.MethodPost {
					r.userOpportunityController.ApproveUserOpportunity(w, req)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			case "reject":
				if req.Method == http.MethodPatch || req.Method == http.MethodPost {
					r.userOpportunityController.RejectUserOpportunity(w, req)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			case "apply":
				if req.Method == http.MethodPost {
					r.userOpportunityController.ApplyToOpportunity(w, req)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			default:
				http.Error(w, "Invalid endpoint", http.StatusNotFound)
			}
			return
		}

		http.Error(w, "Invalid endpoint", http.StatusNotFound)
	}))

	mux.HandleFunc("/api/user-applications", r.authMiddleware.Authenticate(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			r.userOpportunityController.GetUserApplications(w, req)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/logs/error", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "logs/error.txt")
	})
	mux.HandleFunc("/api/logs/info", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "logs/info.txt")
	})
	mux.HandleFunc("/api/logs/data", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "logs/data.txt")
	})
}
