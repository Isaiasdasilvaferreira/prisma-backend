package opportunity

import "time"

type Source string

const (
	SourceGreenhouse Source = "greenhouse"
	SourceLever      Source = "lever"
	SourceAshby      Source = "ashby"
)

type ContractType string
type Modality string

const (
	ContractCLT       ContractType = "CLT"
	ContractFreelancer ContractType = "Freelancer"
)

const (
	ModalityRemoto    Modality = "Remoto"
	ModalityPresencial Modality = "Presencial"
	ModalityHibrido   Modality = "Híbrido"
)

type Opportunity struct {
	ExternalID     string       `json:"external_id"`
	Source         Source       `json:"source"`
	Company        string       `json:"company"`
	Title          string       `json:"title"`
	ContractType   ContractType `json:"contract_type"`
	Modality       Modality     `json:"modality"`
	ServiceType    string       `json:"service_type"`
	Location       string       `json:"location"`
	ApplicationURL string       `json:"application_url"`
	IsActive       bool         `json:"is_active"`
	UserID         string       `json:"user_id"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
