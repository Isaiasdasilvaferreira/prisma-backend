package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	GetUserPlan(ctx context.Context, userID uuid.UUID) (*UserPlan, error)
	CanScrapeOpportunities(ctx context.Context, userID uuid.UUID) (bool, int, error)
	IncrementUsedCount(ctx context.Context, userID uuid.UUID) error
	GetPlanLimit(ctx context.Context, userID uuid.UUID) (int, error)
	GetDailyLimit(planType PlanType) int
}

type service struct {
	planRepo PlanRepository
}

func NewService(planRepo PlanRepository) Service {
	return &service{
		planRepo: planRepo,
	}
}

func (s *service) GetUserPlan(ctx context.Context, userID uuid.UUID) (*UserPlan, error) {
	return s.planRepo.GetUserPlan(ctx, userID)
}

func (s *service) GetDailyLimit(planType PlanType) int {
	return planType.GetDailyLimit()
}

func (s *service) GetPlanLimit(ctx context.Context, userID uuid.UUID) (int, error) {
	plan, err := s.planRepo.GetUserPlan(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("error getting user plan: %w", err)
	}

	if plan == nil {
		return 10, nil
	}

	return s.GetDailyLimit(plan.PlanType), nil
}

func (s *service) CanScrapeOpportunities(ctx context.Context, userID uuid.UUID) (bool, int, error) {
	plan, err := s.planRepo.GetUserPlan(ctx, userID)
	if err != nil {
		return false, 0, fmt.Errorf("error getting user plan: %w", err)
	}

	if plan == nil {
		return true, 10, nil
	}

	dailyLimit := s.GetDailyLimit(plan.PlanType)
	usedToday, err := s.planRepo.GetDailyUsage(ctx, userID)
	if err != nil {
		return false, 0, fmt.Errorf("error getting daily usage: %w", err)
	}

	remaining := dailyLimit - usedToday
	canScrape := remaining > 0

	return canScrape, remaining, nil
}

func (s *service) IncrementUsedCount(ctx context.Context, userID uuid.UUID) error {
	return s.planRepo.IncrementDailyUsage(ctx, userID, 1)
}
