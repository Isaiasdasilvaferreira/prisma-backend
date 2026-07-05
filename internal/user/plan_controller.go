package user

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
)

type PlanController struct {
	planRepo PlanRepository
}

func NewPlanController(planRepo PlanRepository) *PlanController {
	return &PlanController{
		planRepo: planRepo,
	}
}

func (c *PlanController) GetUserPlan(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	plan, err := c.planRepo.GetUserPlan(r.Context(), parsedUserID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to get user plan")
		return
	}

	if plan == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "User plan not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, plan)
}
