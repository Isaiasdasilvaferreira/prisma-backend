package user_opportunity

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type Repository struct {
	client *supabase.Client
	admin  *supabase.Client
}

func NewRepository(client *supabase.Client, admin *supabase.Client) *Repository {
	return &Repository{
		client: client,
		admin:  admin,
	}
}

func (r *Repository) Create(ctx context.Context, opp *UserOpportunity) error {
	if opp.ID == "" {
		opp.ID = uuid.New().String()
	}

	now := time.Now()
	opp.CreatedAt = now
	opp.UpdatedAt = now
	opp.IsActive = false
	opp.ApplicantIDs = []string{}

	if opp.AvailableRegistration != nil {
		opp.RemainingVacancies = opp.AvailableRegistration
	} else {
		remaining := 1
		opp.RemainingVacancies = &remaining
	}

	var result []UserOpportunity
	err := r.admin.DB.From("user_opportunities").
		Insert(opp).
		Execute(&result)

	if err != nil {
		return fmt.Errorf("failed to create user opportunity: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("failed to create user opportunity: no rows returned")
	}

	*opp = result[0]
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*UserOpportunity, error) {
	var result []UserOpportunity
	err := r.admin.DB.From("user_opportunities").
		Select("*").
		Eq("id", id).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to get user opportunity: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("user opportunity not found")
	}

	return &result[0], nil
}

func (r *Repository) GetAll(ctx context.Context, isActive *bool) ([]*UserOpportunity, error) {
	var result []UserOpportunity
	var err error

	if isActive != nil {
		err = r.admin.DB.From("user_opportunities").
			Select("*").
			Eq("is_active", fmt.Sprintf("%v", *isActive)).
			Execute(&result)
	} else {
		err = r.admin.DB.From("user_opportunities").
			Select("*").
			Execute(&result)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list user opportunities: %w", err)
	}

	opportunities := make([]*UserOpportunity, len(result))
	for i := range result {
		opportunities[i] = &result[i]
	}

	return opportunities, nil
}

func (r *Repository) Update(ctx context.Context, opp *UserOpportunity) error {
	opp.UpdatedAt = time.Now()

	var result []UserOpportunity
	err := r.admin.DB.From("user_opportunities").
		Update(opp).
		Eq("id", opp.ID).
		Execute(&result)

	if err != nil {
		return fmt.Errorf("failed to update user opportunity: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("user opportunity not found")
	}

	*opp = result[0]
	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	var result []UserOpportunity
	err := r.admin.DB.From("user_opportunities").
		Delete().
		Eq("id", id).
		Execute(&result)

	if err != nil {
		return fmt.Errorf("failed to delete user opportunity: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("user opportunity not found")
	}

	return nil
}

func (r *Repository) Approve(ctx context.Context, id string) error {
	updateData := map[string]interface{}{
		"is_active":  true,
		"updated_at": time.Now(),
	}

	var result []UserOpportunity
	err := r.admin.DB.From("user_opportunities").
		Update(updateData).
		Eq("id", id).
		Execute(&result)

	if err != nil {
		return fmt.Errorf("failed to approve user opportunity: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("user opportunity not found")
	}

	return nil
}

func (r *Repository) Reject(ctx context.Context, id string) error {
	updateData := map[string]interface{}{
		"is_active":  false,
		"updated_at": time.Now(),
	}

	var result []UserOpportunity
	err := r.admin.DB.From("user_opportunities").
		Update(updateData).
		Eq("id", id).
		Execute(&result)

	if err != nil {
		return fmt.Errorf("failed to reject user opportunity: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("user opportunity not found")
	}

	return nil
}

func (r *Repository) Apply(ctx context.Context, opportunityID string, userID string) (*UserOpportunity, error) {
	opp, err := r.GetByID(ctx, opportunityID)
	if err != nil {
		return nil, err
	}

	if opp.RemainingVacancies == nil || *opp.RemainingVacancies <= 0 {
		return nil, fmt.Errorf("no vacancies available")
	}

	for _, id := range opp.ApplicantIDs {
		if id == userID {
			return nil, fmt.Errorf("user has already applied for this opportunity")
		}
	}

	newApplicantIDs := append(opp.ApplicantIDs, userID)
	newRemaining := *opp.RemainingVacancies - 1

	updateData := map[string]interface{}{
		"applicant_ids":       newApplicantIDs,
		"remaining_vacancies": newRemaining,
		"updated_at":          time.Now(),
	}

	if newRemaining <= 0 {
		updateData["is_active"] = false
	}

	var result []UserOpportunity
	err = r.admin.DB.From("user_opportunities").
		Update(updateData).
		Eq("id", opportunityID).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to update opportunity: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("opportunity not found")
	}

	if newRemaining <= 0 {
		return nil, fmt.Errorf("opportunity fully booked and deactivated")
	}

	return &result[0], nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID string) ([]*UserOpportunity, error) {
	var result []UserOpportunity
	err := r.admin.DB.From("user_opportunities").
		Select("*").
		Filter("applicant_ids", "cs", fmt.Sprintf("{%s}", userID)).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to get user applications: %w", err)
	}

	opportunities := make([]*UserOpportunity, len(result))
	for i := range result {
		opportunities[i] = &result[i]
	}

	return opportunities, nil
}
