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
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	PlanType  PlanType  `db:"plan_type" json:"plan_type"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
