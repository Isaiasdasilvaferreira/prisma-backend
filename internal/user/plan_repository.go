package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type PlanRepository interface {
	GetUserPlan(ctx context.Context, userID uuid.UUID) (*UserPlan, error)
	UpdateUserPlan(ctx context.Context, userID uuid.UUID, planType PlanType) error
	GetDailyUsage(ctx context.Context, userID uuid.UUID) (int, error)
	IncrementDailyUsage(ctx context.Context, userID uuid.UUID, count int) error
	ResetDailyUsage(ctx context.Context, userID uuid.UUID) error
}

type planRepository struct {
	supabase       *supabase.Client
	supabaseAdmin  *supabase.Client
}

func NewPlanRepository(supabase *supabase.Client, supabaseAdmin *supabase.Client) PlanRepository {
	return &planRepository{
		supabase:      supabase,
		supabaseAdmin: supabaseAdmin,
	}
}

func (r *planRepository) getClient() *supabase.Client {
	if r.supabaseAdmin != nil {
		return r.supabaseAdmin
	}
	return r.supabase
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

func (r *planRepository) UpdateUserPlan(ctx context.Context, userID uuid.UUID, planType PlanType) error {
	updates := map[string]interface{}{
		"plan_type": string(planType),
	}

	var result []map[string]interface{}
	err := r.supabase.DB.From("user_plans").
		Update(updates).
		Eq("user_id", userID.String()).
		Execute(&result)
	if err != nil {
		return fmt.Errorf("failed to update user plan: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("user plan not found for user_id: %s", userID)
	}

	return nil
}

func (r *planRepository) GetDailyUsage(ctx context.Context, userID uuid.UUID) (int, error) {
	today := time.Now().Format("2006-01-02")
	client := r.getClient()

	var result []struct {
		UsageCount int `json:"usage_count"`
	}
	err := client.DB.From("daily_usage").
		Select("usage_count").
		Eq("user_id", userID.String()).
		Eq("usage_date", today).
		Execute(&result)
	if err != nil {
		return 0, nil
	}

	if len(result) == 0 {
		return 0, nil
	}

	return result[0].UsageCount, nil
}

func (r *planRepository) IncrementDailyUsage(ctx context.Context, userID uuid.UUID, count int) error {
	today := time.Now().Format("2006-01-02")
	client := r.getClient()

	existing, err := r.GetDailyUsage(ctx, userID)
	if err != nil {
		return err
	}

	var result []map[string]interface{}

	if existing > 0 {
		newUsage := existing + 1
		updateData := map[string]interface{}{
			"usage_count": newUsage,
		}
		err = client.DB.From("daily_usage").
			Update(updateData).
			Eq("user_id", userID.String()).
			Eq("usage_date", today).
			Execute(&result)
		if err != nil {
			return fmt.Errorf("failed to update daily usage: %w", err)
		}
	} else {
		insertData := map[string]interface{}{
			"user_id":     userID.String(),
			"usage_date":  today,
			"usage_count": 1,
		}
		err = client.DB.From("daily_usage").
			Insert(insertData).
			Execute(&result)
		if err != nil {
			return fmt.Errorf("failed to insert daily usage: %w", err)
		}
	}

	return nil
}

func (r *planRepository) ResetDailyUsage(ctx context.Context, userID uuid.UUID) error {
	today := time.Now().Format("2006-01-02")

	updateData := map[string]interface{}{
		"usage_count": 0,
	}

	var result []map[string]interface{}
	err := r.supabase.DB.From("daily_usage").
		Update(updateData).
		Eq("user_id", userID.String()).
		Eq("usage_date", today).
		Execute(&result)
	if err != nil {
		return fmt.Errorf("failed to reset daily usage: %w", err)
	}

	return nil
}
