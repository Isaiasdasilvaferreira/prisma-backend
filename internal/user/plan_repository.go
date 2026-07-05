package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type PlanRepository interface {
	CreateUserPlan(ctx context.Context, userID uuid.UUID) (*UserPlan, error)
	GetUserPlan(ctx context.Context, userID uuid.UUID) (*UserPlan, error)
	UpdateUserPlan(ctx context.Context, userID uuid.UUID, planType PlanType) (*UserPlan, error)
	GetUserPlanByID(ctx context.Context, planID uuid.UUID) (*UserPlan, error)
}

type planRepository struct {
	db *sql.DB
}

func NewPlanRepository(db *sql.DB) PlanRepository {
	return &planRepository{
		db: db,
	}
}

func (r *planRepository) CreateUserPlan(ctx context.Context, userID uuid.UUID) (*UserPlan, error) {
	query := `
		INSERT INTO user_plans (user_id, plan_type, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, user_id, plan_type, created_at, updated_at
	`

	var plan UserPlan
	err := r.db.QueryRowContext(ctx, query, userID, PlanFree).Scan(
		&plan.ID,
		&plan.UserID,
		&plan.PlanType,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user plan: %w", err)
	}

	return &plan, nil
}

func (r *planRepository) GetUserPlan(ctx context.Context, userID uuid.UUID) (*UserPlan, error) {
	query := `
		SELECT id, user_id, plan_type, created_at, updated_at
		FROM user_plans
		WHERE user_id = $1
	`

	var plan UserPlan
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&plan.ID,
		&plan.UserID,
		&plan.PlanType,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user plan: %w", err)
	}

	return &plan, nil
}

func (r *planRepository) UpdateUserPlan(ctx context.Context, userID uuid.UUID, planType PlanType) (*UserPlan, error) {
	query := `
		UPDATE user_plans
		SET plan_type = $1, updated_at = NOW()
		WHERE user_id = $2
		RETURNING id, user_id, plan_type, created_at, updated_at
	`

	var plan UserPlan
	err := r.db.QueryRowContext(ctx, query, planType, userID).Scan(
		&plan.ID,
		&plan.UserID,
		&plan.PlanType,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user plan not found for user_id: %s", userID)
		}
		return nil, fmt.Errorf("failed to update user plan: %w", err)
	}

	return &plan, nil
}

func (r *planRepository) GetUserPlanByID(ctx context.Context, planID uuid.UUID) (*UserPlan, error) {
	query := `
		SELECT id, user_id, plan_type, created_at, updated_at
		FROM user_plans
		WHERE id = $1
	`

	var plan UserPlan
	err := r.db.QueryRowContext(ctx, query, planID).Scan(
		&plan.ID,
		&plan.UserID,
		&plan.PlanType,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user plan by id: %w", err)
	}

	return &plan, nil
}
