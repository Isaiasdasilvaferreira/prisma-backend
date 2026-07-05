package plans

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/user"
	"github.com/nedpals/supabase-go"
)

type PlanService struct {
	supabase *supabase.Client
	planRepo user.PlanRepository
}

func NewPlanService(supabase *supabase.Client) *PlanService {
	return &PlanService{
		supabase: supabase,
		planRepo: user.NewPlanRepository(supabase),
	}
}

func (s *PlanService) CanScrape(ctx context.Context, userID string) (bool, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user ID format: %w", err)
	}

	plan, err := s.planRepo.GetUserPlan(ctx, parsedUserID)
	if err != nil {
		return false, fmt.Errorf("failed to get user plan: %w", err)
	}

	if plan == nil {
		return false, nil
	}

	return plan.PlanType == user.PlanProfessional, nil
}

func (s *PlanService) UpgradeToProfessional(ctx context.Context, userID string) error {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	_, err = s.planRepo.UpdateUserPlan(ctx, parsedUserID, user.PlanProfessional)
	if err != nil {
		return fmt.Errorf("failed to upgrade plan: %w", err)
	}
	return nil
}
