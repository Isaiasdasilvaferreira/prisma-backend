package user_opportunity

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, opp *UserOpportunity) error {
	if opp.ID == "" {
		opp.ID = uuid.New().String()
	}

	query := `
		INSERT INTO user_opportunities (
			id, title, company, contract_type, modality, 
			location, salary, available_registration, whatsapp, email, 
			description, responsibilities, requirements, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9, $10, 
			$11, $12, $13, $14, $15, $16
		)
	`

	now := time.Now()
	opp.CreatedAt = now
	opp.UpdatedAt = now
	opp.IsActive = false

	_, err := r.db.ExecContext(
		ctx,
		query,
		opp.ID,
		opp.Title,
		opp.Company,
		opp.ContractType,
		opp.Modality,
		opp.Location,
		opp.Salary,
		opp.AvailableRegistration,
		opp.WhatsApp,
		opp.Email,
		opp.Description,
		opp.Responsibilities,
		opp.Requirements,
		opp.IsActive,
		opp.CreatedAt,
		opp.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user opportunity: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*UserOpportunity, error) {
	query := `
		SELECT id, title, company, contract_type, modality, 
			location, salary, available_registration, whatsapp, email, 
			description, responsibilities, requirements, is_active, created_at, updated_at
		FROM user_opportunities
		WHERE id = $1
	`

	var opp UserOpportunity
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&opp.ID,
		&opp.Title,
		&opp.Company,
		&opp.ContractType,
		&opp.Modality,
		&opp.Location,
		&opp.Salary,
		&opp.AvailableRegistration,
		&opp.WhatsApp,
		&opp.Email,
		&opp.Description,
		&opp.Responsibilities,
		&opp.Requirements,
		&opp.IsActive,
		&opp.CreatedAt,
		&opp.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user opportunity not found")
		}
		return nil, fmt.Errorf("failed to get user opportunity: %w", err)
	}

	return &opp, nil
}

func (r *Repository) GetAll(ctx context.Context, isActive *bool) ([]*UserOpportunity, error) {
	query := `
		SELECT id, title, company, contract_type, modality, 
			location, salary, available_registration, whatsapp, email, 
			description, responsibilities, requirements, is_active, created_at, updated_at
		FROM user_opportunities
	`

	args := []interface{}{}
	whereClause := ""

	if isActive != nil {
		whereClause = " WHERE is_active = $1"
		args = append(args, *isActive)
	}

	query += whereClause + " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list user opportunities: %w", err)
	}
	defer rows.Close()

	var opportunities []*UserOpportunity
	for rows.Next() {
		var opp UserOpportunity
		err := rows.Scan(
			&opp.ID,
			&opp.Title,
			&opp.Company,
			&opp.ContractType,
			&opp.Modality,
			&opp.Location,
			&opp.Salary,
			&opp.AvailableRegistration,
			&opp.WhatsApp,
			&opp.Email,
			&opp.Description,
			&opp.Responsibilities,
			&opp.Requirements,
			&opp.IsActive,
			&opp.CreatedAt,
			&opp.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user opportunity: %w", err)
		}
		opportunities = append(opportunities, &opp)
	}

	return opportunities, nil
}

func (r *Repository) Update(ctx context.Context, opp *UserOpportunity) error {
	query := `
		UPDATE user_opportunities
		SET title = $1, company = $2, contract_type = $3, modality = $4,
			location = $5, salary = $6, available_registration = $7, 
			whatsapp = $8, email = $9, description = $10, 
			responsibilities = $11, requirements = $12, is_active = $13, updated_at = $14
		WHERE id = $15
	`

	opp.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(
		ctx,
		query,
		opp.Title,
		opp.Company,
		opp.ContractType,
		opp.Modality,
		opp.Location,
		opp.Salary,
		opp.AvailableRegistration,
		opp.WhatsApp,
		opp.Email,
		opp.Description,
		opp.Responsibilities,
		opp.Requirements,
		opp.IsActive,
		opp.UpdatedAt,
		opp.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user opportunity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user opportunity not found")
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM user_opportunities WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user opportunity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user opportunity not found")
	}

	return nil
}

func (r *Repository) Approve(ctx context.Context, id string) error {
	query := `UPDATE user_opportunities SET is_active = true, updated_at = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to approve user opportunity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user opportunity not found")
	}

	return nil
}

func (r *Repository) Reject(ctx context.Context, id string) error {
	query := `UPDATE user_opportunities SET is_active = false, updated_at = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to reject user opportunity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user opportunity not found")
	}

	return nil
}
