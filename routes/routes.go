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
	"github.com/nedpals/supabase-go"
)

type AuthRoutes struct {
	authController       *auth.AuthController
	authMiddleware       *middleware.AuthMiddleware
	planController       *plans.PlanController
	scraperController    *scraper.ScraperController
	opportunityController *opportunity.Controller
}

func NewAuthRoutes(cfg *config.Config, authService *auth.SupabaseAuth, supabaseClient *supabase.Client) *AuthRoutes {
	authMiddleware := middleware.NewAuthMiddleware(authService)
	authController := auth.NewAuthController(authService)

	planRepo := user.NewPlanRepository(supabaseClient)
	userSvc := user.NewService(planRepo)
	oppRepo := opportunity.NewRepository(supabaseClient)

	planService := plans.NewPlanService(planRepo)
	planController := plans.NewPlanController(planService)

	oppService := opportunity.NewService(oppRepo, userSvc)
	opportunityController := opportunity.NewController(oppService)

	scraperController := scraper.NewScraperController(supabaseClient, userSvc, oppRepo)

	return &AuthRoutes{
		authController:       authController,
		authMiddleware:       authMiddleware,
		planController:       planController,
		scraperController:    scraperController,
		opportunityController: opportunityController,
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
			r.opportunityController.GetUserOpportunityByID(w, req)
		} else {
			r.opportunityController.GetUserOpportunities(w, req)
		}
	}))
}
