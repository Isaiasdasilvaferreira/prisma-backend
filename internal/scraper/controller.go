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

	limit := 10
	opps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceAshby, limit)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]opportunity.OpportunityResponse, len(opps))
	for i, opp := range opps {
		response[i] = opportunity.OpportunityResponse{
			ExternalID:     opp.ExternalID,
			Source:         string(opp.Source),
			Company:        opp.Company,
			Title:          opp.Title,
			ContractType:   string(opp.ContractType),
			Modality:       string(opp.Modality),
			ServiceType:    opp.ServiceType,
			Location:       opp.Location,
			ApplicationURL: opp.ApplicationURL,
			IsActive:       opp.IsActive,
		}
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"source":        "ashby",
		"count":         len(response),
		"opportunities": response,
	})
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

	limit := 10
	opps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceGreenhouse, limit)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]opportunity.OpportunityResponse, len(opps))
	for i, opp := range opps {
		response[i] = opportunity.OpportunityResponse{
			ExternalID:     opp.ExternalID,
			Source:         string(opp.Source),
			Company:        opp.Company,
			Title:          opp.Title,
			ContractType:   string(opp.ContractType),
			Modality:       string(opp.Modality),
			ServiceType:    opp.ServiceType,
			Location:       opp.Location,
			ApplicationURL: opp.ApplicationURL,
			IsActive:       opp.IsActive,
		}
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"source":        "greenhouse",
		"count":         len(response),
		"opportunities": response,
	})
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

	limit := 10
	opps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceLever, limit)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]opportunity.OpportunityResponse, len(opps))
	for i, opp := range opps {
		response[i] = opportunity.OpportunityResponse{
			ExternalID:     opp.ExternalID,
			Source:         string(opp.Source),
			Company:        opp.Company,
			Title:          opp.Title,
			ContractType:   string(opp.ContractType),
			Modality:       string(opp.Modality),
			ServiceType:    opp.ServiceType,
			Location:       opp.Location,
			ApplicationURL: opp.ApplicationURL,
			IsActive:       opp.IsActive,
		}
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"source":        "lever",
		"count":         len(response),
		"opportunities": response,
	})
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

	limit := 10
	allOpps := make([]opportunity.Opportunity, 0)

	ashbyOpps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceAshby, limit)
	if err == nil {
		allOpps = append(allOpps, ashbyOpps...)
	}

	greenhouseOpps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceGreenhouse, limit)
	if err == nil {
		allOpps = append(allOpps, greenhouseOpps...)
	}

	leverOpps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceLever, limit)
	if err == nil {
		allOpps = append(allOpps, leverOpps...)
	}

	response := make([]opportunity.OpportunityResponse, len(allOpps))
	for i, opp := range allOpps {
		response[i] = opportunity.OpportunityResponse{
			ExternalID:     opp.ExternalID,
			Source:         string(opp.Source),
			Company:        opp.Company,
			Title:          opp.Title,
			ContractType:   string(opp.ContractType),
			Modality:       string(opp.Modality),
			ServiceType:    opp.ServiceType,
			Location:       opp.Location,
			ApplicationURL: opp.ApplicationURL,
			IsActive:       opp.IsActive,
		}
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"count":         len(response),
		"opportunities": response,
	})
}
