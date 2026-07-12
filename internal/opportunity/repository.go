package opportunity

import (
	"context"
	"fmt"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
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
	supabase       *supabase.Client
	supabaseAdmin  *supabase.Client
}

func NewRepository(supabase *supabase.Client, supabaseAdmin *supabase.Client) Repository {
	return &repository{
		supabase:      supabase,
		supabaseAdmin: supabaseAdmin,
	}
}

func (r *repository) getClient() *supabase.Client {
	if r.supabaseAdmin != nil {
		return r.supabaseAdmin
	}
	return r.supabase
}

func (r *repository) Create(ctx context.Context, opp *Opportunity) error {
	utils.LogInfo(fmt.Sprintf("[Create] Iniciando - ExternalID: %s", opp.ExternalID))
	utils.LogData(opp)

	var result []Opportunity
	err := r.getClient().DB.From("opportunities").
		Insert(map[string]interface{}{
			"external_id":     opp.ExternalID,
			"source":          opp.Source,
			"company":         opp.Company,
			"title":           opp.Title,
			"contract_type":   opp.ContractType,
			"modality":        opp.Modality,
			"service_type":    opp.ServiceType,
			"location":        opp.Location,
			"application_url": opp.ApplicationURL,
			"is_active":       opp.IsActive,
		}).
		Execute(&result)
	if err != nil {
		utils.LogError(fmt.Sprintf("[Create] Erro ao inserir - ExternalID: %s", opp.ExternalID), err)
		return fmt.Errorf("error creating opportunity: %w", err)
	}

	utils.LogInfo(fmt.Sprintf("[Create] Sucesso - ExternalID: %s", opp.ExternalID))
	return nil
}

func (r *repository) CreateMany(ctx context.Context, opps []*Opportunity) error {
	if len(opps) == 0 {
		return nil
	}

	var inserts []map[string]interface{}
	for _, opp := range opps {
		inserts = append(inserts, map[string]interface{}{
			"external_id":     opp.ExternalID,
			"source":          opp.Source,
			"company":         opp.Company,
			"title":           opp.Title,
			"contract_type":   opp.ContractType,
			"modality":        opp.Modality,
			"service_type":    opp.ServiceType,
			"location":        opp.Location,
			"application_url": opp.ApplicationURL,
			"is_active":       opp.IsActive,
		})
	}

	var result []Opportunity
	err := r.getClient().DB.From("opportunities").
		Insert(inserts).
		Execute(&result)
	if err != nil {
		utils.LogError("[CreateMany] Erro ao inserir em lote", err)
		return fmt.Errorf("error creating many opportunities: %w", err)
	}

	utils.LogInfo(fmt.Sprintf("[CreateMany] %d oportunidades inseridas com sucesso", len(result)))
	return nil
}

func (r *repository) GetByExternalID(ctx context.Context, externalID string) (*Opportunity, error) {
	utils.LogInfo(fmt.Sprintf("GetByExternalID - Buscando: %s", externalID))

	var result []Opportunity
	err := r.supabaseAdmin.DB.From("opportunities").
		Select("*").
		Filter("external_id", "eq", externalID).
		Execute(&result)
	if err != nil {
		utils.LogError(fmt.Sprintf("GetByExternalID - Erro ao buscar %s", externalID), err)
		return nil, fmt.Errorf("error getting opportunity: %w", err)
	}

	if len(result) == 0 {
		utils.LogInfo(fmt.Sprintf("GetByExternalID - Nenhum resultado para %s", externalID))
		return nil, nil
	}

	utils.LogInfo(fmt.Sprintf("GetByExternalID - Encontrado: %s", externalID))
	return &result[0], nil
}

func (r *repository) GetByUserID(ctx context.Context, userID string) ([]Opportunity, error) {
	utils.LogInfo("GetByUserID - Buscando TODAS as oportunidades (global)")

	var result []Opportunity
	err := r.supabaseAdmin.DB.From("opportunities").
		Select("*").
		Execute(&result)
	if err != nil {
		utils.LogError("GetByUserID - Erro", err)
		return nil, fmt.Errorf("error getting opportunities: %w", err)
	}

	utils.LogInfo(fmt.Sprintf("GetByUserID - %d oportunidades encontradas", len(result)))
	return result, nil
}

func (r *repository) GetByUserIDWithFilters(ctx context.Context, userID string, source string, limit int) ([]Opportunity, error) {
	utils.LogInfo(fmt.Sprintf("GetByUserIDWithFilters - source: %s, limit: %d (global)", source, limit))

	var result []Opportunity
	query := r.supabaseAdmin.DB.From("opportunities").Select("*")

	if source != "" {
		result, err := query.Filter("source", "eq", source).Execute(&result)
		if err != nil {
			utils.LogError("GetByUserIDWithFilters - Erro", err)
			return nil, fmt.Errorf("error getting opportunities with filters: %w", err)
		}
	} else {
		err := query.Execute(&result)
		if err != nil {
			utils.LogError("GetByUserIDWithFilters - Erro", err)
			return nil, fmt.Errorf("error getting opportunities: %w", err)
		}
	}

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	utils.LogInfo(fmt.Sprintf("GetByUserIDWithFilters - %d oportunidades encontradas", len(result)))
	return result, nil
}

func (r *repository) GetAllActive(ctx context.Context, limit int) ([]Opportunity, error) {
	utils.LogInfo(fmt.Sprintf("GetAllActive - limit: %d", limit))

	var result []Opportunity
	err := r.supabaseAdmin.DB.From("opportunities").
		Select("*").
		Filter("is_active", "eq", "true").
		Execute(&result)
	if err != nil {
		utils.LogError("GetAllActive - Erro", err)
		return nil, fmt.Errorf("error getting active opportunities: %w", err)
	}

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	utils.LogInfo(fmt.Sprintf("GetAllActive - %d oportunidades encontradas", len(result)))
	return result, nil
}

func (r *repository) GetBySource(ctx context.Context, source string, limit int) ([]Opportunity, error) {
	utils.LogInfo(fmt.Sprintf("GetBySource - source: %s, limit: %d", source, limit))

	var result []Opportunity
	err := r.supabaseAdmin.DB.From("opportunities").
		Select("*").
		Filter("source", "eq", source).
		Execute(&result)
	if err != nil {
		utils.LogError("GetBySource - Erro", err)
		return nil, fmt.Errorf("error getting opportunities by source: %w", err)
	}

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	utils.LogInfo(fmt.Sprintf("GetBySource - %d oportunidades encontradas", len(result)))
	return result, nil
}

func (r *repository) CountByUser(ctx context.Context, userID string) (int, error) {
	utils.LogInfo("CountByUser - Contando TODAS as oportunidades (global)")

	var result []Opportunity
	err := r.supabaseAdmin.DB.From("opportunities").
		Select("external_id").
		Execute(&result)
	if err != nil {
		utils.LogError("CountByUser - Erro", err)
		return 0, fmt.Errorf("error counting opportunities: %w", err)
	}

	utils.LogInfo(fmt.Sprintf("CountByUser - %d oportunidades encontradas", len(result)))
	return len(result), nil
}
