package opportunity

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) GetUserOpportunities(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("🔴 GetUserOpportunities CHAMADO")

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.LogError("GetUserOpportunities - User not authenticated", nil)
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	utils.LogInfo(fmt.Sprintf("GetUserOpportunities - UserID do contexto: %s", userID))

	source := r.URL.Query().Get("source")
	limitStr := r.URL.Query().Get("limit")
	limit := 0
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	opps, err := c.service.GetUserOpportunities(r.Context(), userID, source, limit)
	if err != nil {
		utils.LogError("GetUserOpportunities - Erro ao buscar oportunidades", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.LogInfo(fmt.Sprintf("GetUserOpportunities - %d oportunidades encontradas", len(opps)))

	response := make([]OpportunityResponse, len(opps))
	for i, opp := range opps {
		response[i] = toResponse(opp)
	}

	utils.SuccessResponse(w, http.StatusOK, response)
}

func (c *Controller) GetUserOpportunityByExternalID(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("🔴 GetUserOpportunityByExternalID CHAMADO")

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.LogError("GetUserOpportunityByExternalID - User not authenticated", nil)
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	externalID := pathParts[len(pathParts)-1]

	utils.LogInfo(fmt.Sprintf("GetUserOpportunityByExternalID - UserID: %s, ExternalID: %s", userID, externalID))

	opp, err := c.service.GetUserOpportunityByExternalID(r.Context(), userID, externalID)
	if err != nil {
		utils.LogError("GetUserOpportunityByExternalID - Oportunidade não encontrada", err)
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, toResponse(*opp))
}

func (c *Controller) GetOpportunitiesBySource(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("🔴 GetOpportunitiesBySource CHAMADO")

	pathParts := strings.Split(r.URL.Path, "/")
	source := pathParts[len(pathParts)-1]

	limitStr := r.URL.Query().Get("limit")
	limit := 0
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	utils.LogInfo(fmt.Sprintf("GetOpportunitiesBySource - Source: %s, Limit: %d", source, limit))

	opps, err := c.service.GetOpportunitiesBySource(r.Context(), source, limit)
	if err != nil {
		utils.LogError("GetOpportunitiesBySource - Erro ao buscar oportunidades", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.LogInfo(fmt.Sprintf("GetOpportunitiesBySource - %d oportunidades encontradas", len(opps)))

	response := make([]OpportunityResponse, len(opps))
	for i, opp := range opps {
		response[i] = toResponse(opp)
	}

	utils.SuccessResponse(w, http.StatusOK, response)
}

func (c *Controller) GetOpportunitiesStats(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("🔴 GetOpportunitiesStats CHAMADO")

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.LogError("GetOpportunitiesStats - User not authenticated", nil)
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	utils.LogInfo(fmt.Sprintf("GetOpportunitiesStats - UserID: %s", userID))

	stats, err := c.service.GetOpportunitiesStats(r.Context(), userID)
	if err != nil {
		utils.LogError("GetOpportunitiesStats - Erro ao buscar estatísticas", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, stats)
}

func toResponse(opp Opportunity) OpportunityResponse {
	return OpportunityResponse{
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
