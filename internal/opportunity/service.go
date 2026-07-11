package opportunity

import (
	"context"
	"fmt"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/user"
	"github.com/google/uuid"
)

type Service interface {
	GetUserOpportunities(ctx context.Context, userID string, source string, limit int) ([]Opportunity, error)
	GetUserOpportunityByID(ctx context.Context, userID string, oppID string) (*Opportunity, error)
	GetOpportunitiesBySource(ctx context.Context, source string, limit int) ([]Opportunity, error)
	GetOpportunitiesStats(ctx context.Context, userID string) (map[string]interface{}, error)
}

type service struct {
	repo    Repository
	userSvc user.Service
}

func NewService(repo Repository, userSvc user.Service) Service {
	return &service{
		repo:    repo,
		userSvc: userSvc,
	}
}

func (s *service) GetUserOpportunities(ctx context.Context, userID string, source string, limit int) ([]Opportunity, error) {
	return s.repo.GetByUserIDWithFilters(ctx, userID, source, limit)
}

func (s *service) GetUserOpportunityByID(ctx context.Context, userID string, oppID string) (*Opportunity, error) {
	if _, err := uuid.Parse(oppID); err != nil {
		return nil, fmt.Errorf("invalid opportunity ID")
	}

	opps, err := s.repo.GetByUserIDWithFilters(ctx, userID, "", 0)
	if err != nil {
		return nil, err
	}

	for _, opp := range opps {
		if opp.ID == oppID {
			return &opp, nil
		}
	}

	return nil, fmt.Errorf("opportunity not found")
}

func (s *service) GetOpportunitiesBySource(ctx context.Context, source string, limit int) ([]Opportunity, error) {
	return s.repo.GetBySource(ctx, source, limit)
}

func (s *service) GetOpportunitiesStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	opps, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	plan, err := s.userSvc.GetUserPlan(ctx, userUUID)
	if err != nil || plan == nil {
		plan = &user.UserPlan{
			PlanType: user.PlanFree,
		}
	}

	sourceStats := make(map[string]int)
	for _, opp := range opps {
		sourceStats[string(opp.Source)]++
	}

	contractStats := make(map[string]int)
	for _, opp := range opps {
		contractStats[string(opp.ContractType)]++
	}

	modalityStats := make(map[string]int)
	for _, opp := range opps {
		modalityStats[string(opp.Modality)]++
	}

	levelStats := make(map[string]int)
	for _, opp := range opps {
		levelStats[string(opp.Level)]++
	}

	stats := map[string]interface{}{
		"total":        len(opps),
		"plan_type":    plan.PlanType,
		"daily_limit":  plan.PlanType.GetDailyLimit(),
		"by_source":    sourceStats,
		"by_contract":  contractStats,
		"by_modality":  modalityStats,
		"by_level":     levelStats,
		"recent_count": len(opps),
	}

	return stats, nil
}
