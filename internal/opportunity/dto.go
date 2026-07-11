package opportunity

import "time"

type OpportunityResponse struct {
	ID             string    `json:"id"`
	ExternalID     string    `json:"external_id"`
	Source         string    `json:"source"`
	Company        string    `json:"company"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	ContractType   string    `json:"contract_type"`
	Modality       string    `json:"modality"`
	Level          string    `json:"level"`
	ServiceType    string    `json:"service_type"`
	Location       string    `json:"location"`
	SalaryRange    string    `json:"salary_range"`
	ApplicationURL string    `json:"application_url"`
	PostedAt       time.Time `json:"posted_at"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

type ScrapeRequest struct {
	Sources []string `json:"sources"`
	Limit   int      `json:"limit"`
}

type ScrapeResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message,omitempty"`
	Data    []OpportunityResponse    `json:"data,omitempty"`
	Errors  map[string]string        `json:"errors,omitempty"`
}
