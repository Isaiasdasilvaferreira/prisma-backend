package user

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
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
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	plan, err := c.planRepo.GetUserPlan(r.Context(), userID)
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
