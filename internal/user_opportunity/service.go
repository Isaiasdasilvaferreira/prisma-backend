package user_opportunity

import (
	"context"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUserOpportunity(ctx context.Context, req *CreateUserOpportunityRequest) (*UserOpportunity, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.Company == "" {
		return nil, fmt.Errorf("company is required")
	}
	if req.ContractType == "" {
		return nil, fmt.Errorf("contract type is required")
	}
	if req.Modality == "" {
		return nil, fmt.Errorf("modality is required")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if req.Description == "" {
		return nil, fmt.Errorf("description is required")
	}

	opp := &UserOpportunity{
		Title:                 req.Title,
		Company:               req.Company,
		ContractType:          req.ContractType,
		Modality:              req.Modality,
		Location:              req.Location,
		Salary:                req.Salary,
		AvailableRegistration: req.AvailableRegistration,
		WhatsApp:              req.WhatsApp,
		Email:                 req.Email,
		Description:           req.Description,
		Responsibilities:      req.Responsibilities,
		Requirements:          req.Requirements,
	}

	err := s.repo.Create(ctx, opp)
	if err != nil {
		return nil, err
	}

	return opp, nil
}

func (s *Service) GetUserOpportunity(ctx context.Context, id string) (*UserOpportunity, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetAllUserOpportunities(ctx context.Context, isActive *bool) ([]*UserOpportunity, error) {
	return s.repo.GetAll(ctx, isActive)
}

func (s *Service) UpdateUserOpportunity(ctx context.Context, id string, req *UpdateUserOpportunityRequest) (*UserOpportunity, error) {
	opp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		opp.Title = *req.Title
	}
	if req.Company != nil {
		opp.Company = *req.Company
	}
	if req.ContractType != nil {
		opp.ContractType = *req.ContractType
	}
	if req.Modality != nil {
		opp.Modality = *req.Modality
	}
	if req.Location != nil {
		opp.Location = req.Location
	}
	if req.Salary != nil {
		opp.Salary = req.Salary
	}
	if req.AvailableRegistration != nil {
		opp.AvailableRegistration = req.AvailableRegistration
	}
	if req.WhatsApp != nil {
		opp.WhatsApp = req.WhatsApp
	}
	if req.Email != nil {
		opp.Email = *req.Email
	}
	if req.Description != nil {
		opp.Description = *req.Description
	}
	if req.Responsibilities != nil {
		opp.Responsibilities = req.Responsibilities
	}
	if req.Requirements != nil {
		opp.Requirements = req.Requirements
	}
	if req.IsActive != nil {
		opp.IsActive = *req.IsActive
	}

	err = s.repo.Update(ctx, opp)
	if err != nil {
		return nil, err
	}

	return opp, nil
}

func (s *Service) DeleteUserOpportunity(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) ApproveUserOpportunity(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Approve(ctx, id)
}

func (s *Service) RejectUserOpportunity(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Reject(ctx, id)
}

func (s *Service) ApplyToOpportunity(ctx context.Context, opportunityID string, userID string) (*UserOpportunity, error) {
	opp, err := s.repo.GetByID(ctx, opportunityID)
	if err != nil {
		return nil, err
	}

	if !opp.IsActive {
		return nil, fmt.Errorf("opportunity is not active")
	}

	return s.repo.Apply(ctx, opportunityID, userID)
}

func (s *Service) GetUserApplications(ctx context.Context, userID string) ([]*UserOpportunity, error) {
	return s.repo.GetByUserID(ctx, userID)
}
