package plans

import (
	"context"
	"fmt"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/user"
	"github.com/nedpals/supabase-go"
)

type PlanService struct {
	supabase  *supabase.Client
	planRepo  user.PlanRepository
}

func NewPlanService(supabase *supabase.Client) *PlanService {
	return &PlanService{
		supabase: supabase,
		planRepo: user.NewPlanRepository(supabase),
	}
}

func (s *PlanService) CanScrape(ctx context.Context, userID string) (bool, error) {
	plan, err := s.planRepo.GetUserPlan(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user plan: %w", err)
	}

	if plan == nil {
		return false, nil
	}

	return plan.PlanType == user.PlanProfessional, nil
}

func (s *PlanService) UpgradeToProfessional(ctx context.Context, userID string) error {
	_, err := s.planRepo.UpdateUserPlan(ctx, userID, user.PlanProfessional)
	if err != nil {
		return fmt.Errorf("failed to upgrade plan: %w", err)
	}
	return nil
}
