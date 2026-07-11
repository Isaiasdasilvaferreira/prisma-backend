package opportunity

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type Repository interface {
	Create(ctx context.Context, opp *Opportunity) error
	CreateMany(ctx context.Context, opps []*Opportunity) error
	GetByExternalID(ctx context.Context, externalID string) (*Opportunity, error)
	GetByUserID(ctx context.Context, userID string) ([]Opportunity, error)
	GetByUserIDWithFilters(ctx context.Context, userID string, source string, limit int) ([]Opportunity, error)
	GetAllActive(ctx context.Context, limit int) ([]Opportunity, error)
	GetBySource(ctx context.Context, source string, limit int) ([]Opportunity, error)
	CountByUser(ctx context.Context, userID string) (int, error)
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

func (r *repository) GetByUserID(ctx context.Context, userID string) ([]Opportunity, error) {
	var result []Opportunity
	err := r.supabase.DB.From("opportunities").
		Select("*").
		Eq("user_id", userID).
		Order("created_at", &supabase.OrderOpts{Descending: true}).
		Execute(&result)
	if err != nil {
		return nil, fmt.Errorf("error getting opportunities by user: %w", err)
	}

	return result, nil
}

func (r *repository) GetByUserIDWithFilters(ctx context.Context, userID string, source string, limit int) ([]Opportunity, error) {
	query := r.supabase.DB.From("opportunities").
		Select("*").
		Eq("user_id", userID).
		Eq("is_active", true)

	if source != "" {
		query = query.Eq("source", source)
	}

	if limit > 0 {
		query = query.Limit(uint64(limit))
	}

	query = query.Order("created_at", &supabase.OrderOpts{Descending: true})

	var result []Opportunity
	err := query.Execute(&result)
	if err != nil {
		return nil, fmt.Errorf("error getting opportunities with filters: %w", err)
	}

	return result, nil
}

func (r *repository) GetAllActive(ctx context.Context, limit int) ([]Opportunity, error) {
	query := r.supabase.DB.From("opportunities").
		Select("*").
		Eq("is_active", true).
		Order("created_at", &supabase.OrderOpts{Descending: true})

	if limit > 0 {
		query = query.Limit(uint64(limit))
	}

	var result []Opportunity
	err := query.Execute(&result)
	if err != nil {
		return nil, fmt.Errorf("error getting active opportunities: %w", err)
	}

	return result, nil
}

func (r *repository) GetBySource(ctx context.Context, source string, limit int) ([]Opportunity, error) {
	query := r.supabase.DB.From("opportunities").
		Select("*").
		Eq("source", source).
		Eq("is_active", true)

	if limit > 0 {
		query = query.Limit(uint64(limit))
	}

	var result []Opportunity
	err := query.Execute(&result)
	if err != nil {
		return nil, fmt.Errorf("error getting opportunities by source: %w", err)
	}

	return result, nil
}

func (r *repository) CountByUser(ctx context.Context, userID string) (int, error) {
	var result []struct {
		Count int `json:"count"`
	}
	err := r.supabase.DB.From("opportunities").
		Select("count", "exact").
		Eq("user_id", userID).
		Execute(&result)
	if err != nil {
		return 0, fmt.Errorf("error counting opportunities: %w", err)
	}

	if len(result) == 0 {
		return 0, nil
	}

	return result[0].Count, nil
}
