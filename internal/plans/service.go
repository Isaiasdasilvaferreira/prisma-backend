package plans

import (
	"context"
	"fmt"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/user"
	"github.com/google/uuid"
)

type PlanService interface {
	GetUserPlan(ctx context.Context, userID uuid.UUID) (*user.UserPlan, error)
	UpgradeToProfessional(ctx context.Context, userID uuid.UUID) error
	CanScrape(ctx context.Context, userID uuid.UUID) (bool, int, error)
	GetDailyLimit(ctx context.Context, userID uuid.UUID) (int, error)
}

type planService struct {
	planRepo user.PlanRepository
}

func NewPlanService(planRepo user.PlanRepository) PlanService {
	return &planService{
		planRepo: planRepo,
	}
}

func (s *planService) GetUserPlan(ctx context.Context, userID uuid.UUID) (*user.UserPlan, error) {
	return s.planRepo.GetUserPlan(ctx, userID)
}

func (s *planService) UpgradeToProfessional(ctx context.Context, userID uuid.UUID) error {
	return s.planRepo.UpdateUserPlan(ctx, userID, user.PlanProfessional)
}

func (s *planService) CanScrape(ctx context.Context, userID uuid.UUID) (bool, int, error) {
	plan, err := s.planRepo.GetUserPlan(ctx, userID)
	if err != nil {
		return false, 0, fmt.Errorf("error getting user plan: %w", err)
	}

	if plan == nil {
		return true, 10, nil
	}

	dailyLimit := plan.PlanType.GetDailyLimit()
	usedToday, err := s.planRepo.GetDailyUsage(ctx, userID)
	if err != nil {
		return false, 0, fmt.Errorf("error getting daily usage: %w", err)
	}

	remaining := dailyLimit - usedToday
	canScrape := remaining > 0

	return canScrape, remaining, nil
}

func (s *planService) GetDailyLimit(ctx context.Context, userID uuid.UUID) (int, error) {
	plan, err := s.planRepo.GetUserPlan(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("error getting user plan: %w", err)
	}

	if plan == nil {
		return 10, nil
	}

	return plan.PlanType.GetDailyLimit(), nil
}
