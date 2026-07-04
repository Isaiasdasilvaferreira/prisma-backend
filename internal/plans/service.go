package plans

import (
	"context"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/database"
	"github.com/google/uuid"
)

type PlanService struct {
	db *database.Database
}

func NewPlanService(db *database.Database) *PlanService {
	return &PlanService{db: db}
}

func (s *PlanService) CanScrape(ctx context.Context, userID string) (bool, error) {
	var results []map[string]interface{}
	err := s.db.Supabase.DB.From("user_plans").
		Select("plan_type").
		Eq("user_id", userID).
		Execute(&results)

	if err != nil {
		return false, err
	}

	if len(results) == 0 {
		return false, s.createStarterPlan(ctx, userID)
	}

	planType := int(results[0]["plan_type"].(float64))
	return planType == 2, nil
}

func (s *PlanService) createStarterPlan(ctx context.Context, userID string) error {
	var result []map[string]interface{}
	err := s.db.Supabase.DB.From("user_plans").
		Insert(map[string]interface{}{
			"id":        uuid.New().String(),
			"user_id":   userID,
			"plan_type": 1,
		}).
		Execute(&result)
	return err
}

func (s *PlanService) UpgradeToProfessional(ctx context.Context, userID string) error {
	var result []map[string]interface{}
	
	var existing []map[string]interface{}
	err := s.db.Supabase.DB.From("user_plans").
		Select("id").
		Eq("user_id", userID).
		Execute(&existing)
	
	if err != nil {
		return err
	}
	
	if len(existing) == 0 {
		err = s.db.Supabase.DB.From("user_plans").
			Insert(map[string]interface{}{
				"id":        uuid.New().String(),
				"user_id":   userID,
				"plan_type": 2,
			}).
			Execute(&result)
		return err
	}
	
	err = s.db.Supabase.DB.From("user_plans").
		Update(map[string]interface{}{
			"plan_type": 2,
		}).
		Eq("user_id", userID).
		Execute(&result)
	
	return err
}
