package scraper

import (
	"net/http"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/opportunity"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/user"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type ScraperController struct {
	service *ScraperService
}

func NewScraperController(supabase *supabase.Client, userSvc user.Service, oppRepo opportunity.Repository) *ScraperController {
	return &ScraperController{
		service: NewScraperService(supabase, userSvc, oppRepo),
	}
}

func (c *ScraperController) TriggerScraping(w http.ResponseWriter, r *http.Request) {
	userRole := r.Context().Value(auth.UserRoleKey)
	if userRole != "admin" {
		utils.ErrorResponse(w, http.StatusForbidden, "Apenas administradores")
		return
	}

	go func() {
		ctx := r.Context()
		c.service.RunScraping(ctx)
	}()

	utils.SuccessResponse(w, http.StatusAccepted, map[string]string{
		"message": "Scraping iniciado",
	})
}

func (c *ScraperController) ScrapeAshby(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	plan, err := c.service.userSvc.GetUserPlan(r.Context(), userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	limit := 10
	if plan != nil && plan.PlanType == user.PlanProfessional {
		limit = 30
	}

	opps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceAshby, limit)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, opps)
}

func (c *ScraperController) ScrapeGreenhouse(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	plan, err := c.service.userSvc.GetUserPlan(r.Context(), userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	limit := 10
	if plan != nil && plan.PlanType == user.PlanProfessional {
		limit = 30
	}

	opps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceGreenhouse, limit)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, opps)
}

func (c *ScraperController) ScrapeLever(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	plan, err := c.service.userSvc.GetUserPlan(r.Context(), userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	limit := 10
	if plan != nil && plan.PlanType == user.PlanProfessional {
		limit = 30
	}

	opps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceLever, limit)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, opps)
}

func (c *ScraperController) ScrapeAll(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	plan, err := c.service.userSvc.GetUserPlan(r.Context(), userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	limit := 10
	if plan != nil && plan.PlanType == user.PlanProfessional {
		limit = 30
	}

	type ScrapeResult struct {
		Source string                        `json:"source"`
		Count  int                           `json:"count"`
		Data   []opportunity.Opportunity     `json:"data"`
		Error  string                        `json:"error,omitempty"`
	}

	results := make([]ScrapeResult, 0)

	ashbyOpps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceAshby, limit)
	if err != nil {
		results = append(results, ScrapeResult{Source: "ashby", Error: err.Error()})
	} else {
		results = append(results, ScrapeResult{Source: "ashby", Count: len(ashbyOpps), Data: ashbyOpps})
	}

	greenhouseOpps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceGreenhouse, limit)
	if err != nil {
		results = append(results, ScrapeResult{Source: "greenhouse", Error: err.Error()})
	} else {
		results = append(results, ScrapeResult{Source: "greenhouse", Count: len(greenhouseOpps), Data: greenhouseOpps})
	}

	leverOpps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceLever, limit)
	if err != nil {
		results = append(results, ScrapeResult{Source: "lever", Error: err.Error()})
	} else {
		results = append(results, ScrapeResult{Source: "lever", Count: len(leverOpps), Data: leverOpps})
	}

	utils.SuccessResponse(w, http.StatusOK, results)
}
