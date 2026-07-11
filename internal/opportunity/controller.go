package opportunity

import (
	"net/http"
	"strconv"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
	"github.com/gorilla/mux"
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
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

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
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]OpportunityResponse, len(opps))
	for i, opp := range opps {
		response[i] = toResponse(opp)
	}

	utils.SuccessResponse(w, http.StatusOK, response)
}

func (c *Controller) GetUserOpportunityByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	oppID := vars["id"]

	opp, err := c.service.GetUserOpportunityByID(r.Context(), userID, oppID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, toResponse(*opp))
}

func (c *Controller) GetOpportunitiesBySource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	source := vars["source"]

	limitStr := r.URL.Query().Get("limit")
	limit := 0
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	opps, err := c.service.GetOpportunitiesBySource(r.Context(), source, limit)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]OpportunityResponse, len(opps))
	for i, opp := range opps {
		response[i] = toResponse(opp)
	}

	utils.SuccessResponse(w, http.StatusOK, response)
}

func (c *Controller) GetOpportunitiesStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	stats, err := c.service.GetOpportunitiesStats(r.Context(), userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, stats)
}

func toResponse(opp Opportunity) OpportunityResponse {
	return OpportunityResponse{
		ID:             opp.ID,
		ExternalID:     opp.ExternalID,
		Source:         string(opp.Source),
		Company:        opp.Company,
		Title:          opp.Title,
		Description:    opp.Description,
		ContractType:   string(opp.ContractType),
		Modality:       string(opp.Modality),
		Level:          string(opp.Level),
		ServiceType:    opp.ServiceType,
		Location:       opp.Location,
		SalaryRange:    opp.SalaryRange,
		ApplicationURL: opp.ApplicationURL,
		PostedAt:       opp.PostedAt,
		IsActive:       opp.IsActive,
		CreatedAt:      opp.CreatedAt,
	}
}
