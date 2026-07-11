package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type Repository interface {
	GetUserPlan(ctx context.Context, userID uuid.UUID) (*UserPlan, error)
	UpdateUserPlan(ctx context.Context, userID uuid.UUID, planType PlanType) error
}

type repository struct {
	supabase *supabase.Client
}

func NewRepository(supabase *supabase.Client) Repository {
	return &repository{
		supabase: supabase,
	}
}

func (r *repository) GetUserPlan(ctx context.Context, userID uuid.UUID) (*UserPlan, error) {
	var result []UserPlan
	err := r.supabase.DB.From("user_plans").
		Select("*").
		Eq("user_id", userID.String()).
		Execute(&result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func (r *repository) UpdateUserPlan(ctx context.Context, userID uuid.UUID, planType PlanType) error {
	updates := map[string]interface{}{
		"plan_type": string(planType),
	}

	var result []map[string]interface{}
	err := r.supabase.DB.From("user_plans").
		Update(updates).
		Eq("user_id", userID.String()).
		Execute(&result)
	if err != nil {
		return err
	}

	return nil
}
