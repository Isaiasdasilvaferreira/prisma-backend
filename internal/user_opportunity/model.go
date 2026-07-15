package user_opportunity

import (
	"time"
)

type UserOpportunity struct {
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

type CreateUserOpportunityRequest struct {
	Title                 string  `json:"title"`
	Company               string  `json:"company"`
	ContractType          string  `json:"contract_type"`
	Modality              string  `json:"modality"`
	Location              *string `json:"location"`
	Salary                *string `json:"salary"`
	AvailableRegistration *int    `json:"available_registration"`
	WhatsApp              *string `json:"whatsapp"`
	Email                 string  `json:"email"`
	Description           string  `json:"description"`
	Responsibilities      *string `json:"responsibilities"`
	Requirements          *string `json:"requirements"`
}

type UpdateUserOpportunityRequest struct {
	Title                 *string `json:"title"`
	Company               *string `json:"company"`
	ContractType          *string `json:"contract_type"`
	Modality              *string `json:"modality"`
	Location              *string `json:"location"`
	Salary                *string `json:"salary"`
	AvailableRegistration *int    `json:"available_registration"`
	WhatsApp              *string `json:"whatsapp"`
	Email                 *string `json:"email"`
	Description           *string `json:"description"`
	Responsibilities      *string `json:"responsibilities"`
	Requirements          *string `json:"requirements"`
	IsActive              *bool   `json:"is_active"`
}
