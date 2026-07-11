package plans

import (
	"net/http"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
	"github.com/google/uuid"
)

type PlanController struct {
	planService PlanService
}

func NewPlanController(planService PlanService) *PlanController {
	return &PlanController{
		planService: planService,
	}
}

func (c *PlanController) CanScrape(w http.ResponseWriter, r *http.Request) {
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

	canScrape, remaining, err := c.planService.CanScrape(r.Context(), userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"can_scrape": canScrape,
		"remaining":  remaining,
	})
}

func (c *PlanController) UpgradeToProfessional(w http.ResponseWriter, r *http.Request) {
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

	err = c.planService.UpgradeToProfessional(r.Context(), userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{
		"message": "Plan upgraded to professional",
	})
}

func (c *PlanController) GetUserPlan(w http.ResponseWriter, r *http.Request) {
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

	plan, err := c.planService.GetUserPlan(r.Context(), userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, plan)
}
