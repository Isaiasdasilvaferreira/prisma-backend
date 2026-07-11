package user

import (
	"time"

	"github.com/google/uuid"
)

type PlanType string

const (
	PlanFree        PlanType = "free"
	PlanProfessional PlanType = "professional"
)

type UserPlan struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	PlanType  PlanType  `json:"plan_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DailyUsage struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	UsageDate  string    `json:"usage_date"`
	UsageCount int       `json:"usage_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (p PlanType) GetDailyLimit() int {
	switch p {
	case PlanProfessional:
		return 30
	case PlanFree:
		return 10
	default:
		return 10
	}
}
