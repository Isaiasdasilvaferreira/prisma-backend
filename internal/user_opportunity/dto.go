package user_opportunity

import "time"

type UserOpportunityResponse struct {
	ID                     string     `json:"id"`
	Title                  string     `json:"title"`
	Company                string     `json:"company"`
	ContractType           string     `json:"contract_type"`
	Modality               string     `json:"modality"`
	Location               *string    `json:"location"`
	Salary                 *string    `json:"salary"`
	AvailableRegistration  *int       `json:"available_registration"`
	WhatsApp               *string    `json:"whatsapp"`
	Email                  string     `json:"email"`
	Description            string     `json:"description"`
	Responsibilities       *string    `json:"responsibilities"`
	Requirements           *string    `json:"requirements"`
	IsActive               bool       `json:"is_active"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

func (u *UserOpportunity) ToResponse() UserOpportunityResponse {
	return UserOpportunityResponse{
		ID:                    u.ID,
		Title:                 u.Title,
		Company:               u.Company,
		ContractType:          u.ContractType,
		Modality:              u.Modality,
		Location:              u.Location,
		Salary:                u.Salary,
		AvailableRegistration: u.AvailableRegistration,
		WhatsApp:              u.WhatsApp,
		Email:                 u.Email,
		Description:           u.Description,
		Responsibilities:      u.Responsibilities,
		Requirements:          u.Requirements,
		IsActive:              u.IsActive,
		CreatedAt:             u.CreatedAt,
		UpdatedAt:             u.UpdatedAt,
	}
}

func ToResponseList(opportunities []*UserOpportunity) []UserOpportunityResponse {
	responses := make([]UserOpportunityResponse, len(opportunities))
	for i, opp := range opportunities {
		responses[i] = opp.ToResponse()
	}
	return responses
}
