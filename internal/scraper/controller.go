package scraper

import (
	"fmt"
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
	utils.LogInfo("🔴 ScrapeAshby INICIADO")

	userIDStr, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.LogError("ScrapeAshby - User not authenticated", nil)
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.LogError("ScrapeAshby - Invalid user ID", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	utils.LogInfo(fmt.Sprintf("ScrapeAshby - UserID: %s", userID.String()))

	limit := 10
	opps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceAshby, limit)
	if err != nil {
		utils.LogError("ScrapeAshby - Erro no ScrapeForUser", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.LogInfo(fmt.Sprintf("ScrapeAshby - %d oportunidades encontradas", len(opps)))

	if err := c.service.saveOpportunities(r.Context(), opps); err != nil {
		utils.LogError("ScrapeAshby - Erro ao salvar oportunidades", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Erro ao salvar oportunidades")
		return
	}

	utils.LogInfo("ScrapeAshby - Oportunidades salvas com sucesso")

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

	utils.LogInfo("ScrapeAshby - FINALIZADO")
	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"source":        "ashby",
		"count":         len(response),
		"opportunities": response,
	})
}

func (c *ScraperController) ScrapeGreenhouse(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("🔴 ScrapeGreenhouse INICIADO")

	userIDStr, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.LogError("ScrapeGreenhouse - User not authenticated", nil)
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.LogError("ScrapeGreenhouse - Invalid user ID", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	utils.LogInfo(fmt.Sprintf("ScrapeGreenhouse - UserID: %s", userID.String()))

	limit := 10
	opps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceGreenhouse, limit)
	if err != nil {
		utils.LogError("ScrapeGreenhouse - Erro no ScrapeForUser", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.LogInfo(fmt.Sprintf("ScrapeGreenhouse - %d oportunidades encontradas", len(opps)))

	if err := c.service.saveOpportunities(r.Context(), opps); err != nil {
		utils.LogError("ScrapeGreenhouse - Erro ao salvar oportunidades", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Erro ao salvar oportunidades")
		return
	}

	utils.LogInfo("ScrapeGreenhouse - Oportunidades salvas com sucesso")

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

	utils.LogInfo("ScrapeGreenhouse - FINALIZADO")
	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"source":        "greenhouse",
		"count":         len(response),
		"opportunities": response,
	})
}

func (c *ScraperController) ScrapeLever(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("🔴 ScrapeLever INICIADO")

	userIDStr, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.LogError("ScrapeLever - User not authenticated", nil)
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.LogError("ScrapeLever - Invalid user ID", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	utils.LogInfo(fmt.Sprintf("ScrapeLever - UserID: %s", userID.String()))

	limit := 10
	opps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceLever, limit)
	if err != nil {
		utils.LogError("ScrapeLever - Erro no ScrapeForUser", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.LogInfo(fmt.Sprintf("ScrapeLever - %d oportunidades encontradas", len(opps)))

	if err := c.service.saveOpportunities(r.Context(), opps); err != nil {
		utils.LogError("ScrapeLever - Erro ao salvar oportunidades", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Erro ao salvar oportunidades")
		return
	}

	utils.LogInfo("ScrapeLever - Oportunidades salvas com sucesso")

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

	utils.LogInfo("ScrapeLever - FINALIZADO")
	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"source":        "lever",
		"count":         len(response),
		"opportunities": response,
	})
}

func (c *ScraperController) ScrapeAll(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("🔴 ScrapeAll INICIADO")

	userIDStr, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.LogError("ScrapeAll - User not authenticated", nil)
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.LogError("ScrapeAll - Invalid user ID", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	utils.LogInfo(fmt.Sprintf("ScrapeAll - UserID: %s", userID.String()))

	limit := 10
	allOpps := make([]opportunity.Opportunity, 0)

	utils.LogInfo("ScrapeAll - Raspando Ashby...")
	ashbyOpps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceAshby, limit)
	if err == nil {
		allOpps = append(allOpps, ashbyOpps...)
		utils.LogInfo(fmt.Sprintf("ScrapeAll - Ashby retornou %d oportunidades", len(ashbyOpps)))
	} else {
		utils.LogError("ScrapeAll - Erro no Ashby", err)
	}

	utils.LogInfo("ScrapeAll - Raspando Greenhouse...")
	greenhouseOpps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceGreenhouse, limit)
	if err == nil {
		allOpps = append(allOpps, greenhouseOpps...)
		utils.LogInfo(fmt.Sprintf("ScrapeAll - Greenhouse retornou %d oportunidades", len(greenhouseOpps)))
	} else {
		utils.LogError("ScrapeAll - Erro no Greenhouse", err)
	}

	utils.LogInfo("ScrapeAll - Raspando Lever...")
	leverOpps, err := c.service.ScrapeForUser(r.Context(), userID, opportunity.SourceLever, limit)
	if err == nil {
		allOpps = append(allOpps, leverOpps...)
		utils.LogInfo(fmt.Sprintf("ScrapeAll - Lever retornou %d oportunidades", len(leverOpps)))
	} else {
		utils.LogError("ScrapeAll - Erro no Lever", err)
	}

	utils.LogInfo(fmt.Sprintf("ScrapeAll - Total de oportunidades: %d", len(allOpps)))

	if err := c.service.saveOpportunities(r.Context(), allOpps); err != nil {
		utils.LogError("ScrapeAll - Erro ao salvar oportunidades", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Erro ao salvar oportunidades")
		return
	}

	utils.LogInfo("ScrapeAll - Oportunidades salvas com sucesso")

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

	utils.LogInfo("ScrapeAll - FINALIZADO com sucesso")
	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"count":         len(response),
		"opportunities": response,
	})
}
