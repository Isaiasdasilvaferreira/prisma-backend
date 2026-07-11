package plans

import (
	"net/http"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/user"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type PlanController struct {
	supabase   *supabase.Client
	planRepo   user.PlanRepository
	userSvc    user.Service
}

func NewPlanController(supabase *supabase.Client, planRepo user.PlanRepository, userSvc user.Service) *PlanController {
	return &PlanController{
		supabase: supabase,
		planRepo: planRepo,
		userSvc:  userSvc,
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

	canScrape, remaining, err := c.userSvc.CanScrapeOpportunities(r.Context(), userID)
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

	err = c.planRepo.UpdateUserPlan(r.Context(), userID, user.PlanProfessional)
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

	plan, err := c.userSvc.GetUserPlan(r.Context(), userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, plan)
}
