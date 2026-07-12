package opportunity

import "time"

type OpportunityResponse struct {
	ID             string    `json:"id"`
	ExternalID     string    `json:"external_id"`
	Source         string    `json:"source"`
	Company        string    `json:"company"`
	Title          string    `json:"title"`
	ContractType   string    `json:"contract_type"`
	Modality       string    `json:"modality"`
	ServiceType    string    `json:"service_type"`
	Location       string    `json:"location"`
	ApplicationURL string    `json:"application_url"`
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
