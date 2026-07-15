package opportunity

import (
	"context"
	"fmt"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/user"
	"github.com/google/uuid"
)

type Service interface {
	GetUserOpportunities(ctx context.Context, userID string, source string, limit int) ([]Opportunity, error)
	GetUserOpportunityByExternalID(ctx context.Context, userID string, externalID string) (*Opportunity, error)
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
	opps, err := s.repo.GetByUserIDWithFilters(ctx, userID, source, 0)
	if err != nil {
		return nil, err
	}

	if limit > 0 && len(opps) > limit {
		opps = opps[:limit]
	}

	return opps, nil
}

func (s *service) GetUserOpportunityByExternalID(ctx context.Context, userID string, externalID string) (*Opportunity, error) {
	opp, err := s.repo.GetByExternalID(ctx, externalID)
	if err != nil {
		return nil, err
	}
	if opp == nil {
		return nil, fmt.Errorf("opportunity not found")
	}
	return opp, nil
}

func (s *service) GetOpportunitiesBySource(ctx context.Context, source string, limit int) ([]Opportunity, error) {
	opps, err := s.repo.GetBySource(ctx, source, 0)
	if err != nil {
		return nil, err
	}

	if limit > 0 && len(opps) > limit {
		opps = opps[:limit]
	}

	return opps, nil
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

	stats := map[string]interface{}{
		"total":        len(opps),
		"plan_type":    plan.PlanType,
		"daily_limit":  plan.PlanType.GetDailyLimit(),
		"by_source":    sourceStats,
		"by_contract":  contractStats,
		"by_modality":  modalityStats,
		"recent_count": len(opps),
	}

	return stats, nil
}
