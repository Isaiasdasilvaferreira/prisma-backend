package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type PlanRepository interface {
	CreateUserPlan(ctx context.Context, userID uuid.UUID) (*UserPlan, error)
	GetUserPlan(ctx context.Context, userID uuid.UUID) (*UserPlan, error)
	UpdateUserPlan(ctx context.Context, userID uuid.UUID, planType PlanType) (*UserPlan, error)
	GetUserPlanByID(ctx context.Context, planID uuid.UUID) (*UserPlan, error)
}

type planRepository struct {
	supabase *supabase.Client
}

func NewPlanRepository(supabase *supabase.Client) PlanRepository {
	return &planRepository{
		supabase: supabase,
	}
}

func (r *planRepository) CreateUserPlan(ctx context.Context, userID uuid.UUID) (*UserPlan, error) {
	plan := &UserPlan{
		ID:        uuid.New(),
		UserID:    userID,
		PlanType:  PlanFree,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var result []UserPlan
	err := r.supabase.DB.From("user_plans").
		Insert(plan).
		Execute(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to create user plan: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no data returned from insert")
	}

	return &result[0], nil
}

func (r *planRepository) GetUserPlan(ctx context.Context, userID uuid.UUID) (*UserPlan, error) {
	var result []UserPlan
	err := r.supabase.DB.From("user_plans").
		Select("*").
		Eq("user_id", userID.String()).
		Execute(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to get user plan: %w", err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func (r *planRepository) UpdateUserPlan(ctx context.Context, userID uuid.UUID, planType PlanType) (*UserPlan, error) {
	updates := map[string]interface{}{
		"plan_type":  string(planType),
		"updated_at": time.Now(),
	}

	var result []UserPlan
	err := r.supabase.DB.From("user_plans").
		Update(updates).
		Eq("user_id", userID.String()).
		Execute(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to update user plan: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("user plan not found for user_id: %s", userID)
	}

	return &result[0], nil
}

func (r *planRepository) GetUserPlanByID(ctx context.Context, planID uuid.UUID) (*UserPlan, error) {
	var result []UserPlan
	err := r.supabase.DB.From("user_plans").
		Select("*").
		Eq("id", planID.String()).
		Execute(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to get user plan by id: %w", err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}
