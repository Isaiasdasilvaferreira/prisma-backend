package opportunity

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type Repository interface {
	Create(ctx context.Context, opp *Opportunity) error
	CreateMany(ctx context.Context, opps []*Opportunity) error
	GetByExternalID(ctx context.Context, externalID string) (*Opportunity, error)
}

type repository struct {
	supabase *supabase.Client
}

func NewRepository(supabase *supabase.Client) Repository {
	return &repository{
		supabase: supabase,
	}
}

func (r *repository) Create(ctx context.Context, opp *Opportunity) error {
	var result []Opportunity
	err := r.supabase.DB.From("opportunities").
		Insert(map[string]interface{}{
			"id":              uuid.New().String(),
			"external_id":     opp.ExternalID,
			"source":          opp.Source,
			"company":         opp.Company,
			"title":           opp.Title,
			"description":     opp.Description,
			"contract_type":   opp.ContractType,
			"modality":        opp.Modality,
			"level":           opp.Level,
			"service_type":    opp.ServiceType,
			"location":        opp.Location,
			"salary_range":    opp.SalaryRange,
			"application_url": opp.ApplicationURL,
			"posted_at":       opp.PostedAt,
			"is_active":       opp.IsActive,
		}).
		Execute(&result)
	if err != nil {
		return fmt.Errorf("error creating opportunity: %w", err)
	}

	return nil
}

func (r *repository) CreateMany(ctx context.Context, opps []*Opportunity) error {
	if len(opps) == 0 {
		return nil
	}

	for _, opp := range opps {
		if err := r.Create(ctx, opp); err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) GetByExternalID(ctx context.Context, externalID string) (*Opportunity, error) {
	var result []Opportunity
	err := r.supabase.DB.From("opportunities").
		Select("*").
		Eq("external_id", externalID).
		Execute(&result)
	if err != nil {
		return nil, fmt.Errorf("error getting opportunity: %w", err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}
