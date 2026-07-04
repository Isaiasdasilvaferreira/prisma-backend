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
type Level string

const (
	ContractCLT       ContractType = "CLT"
	ContractFreelancer ContractType = "Freelancer"
)

const (
	ModalityRemoto    Modality = "Remoto"
	ModalityPresencial Modality = "Presencial"
	ModalityHibrido   Modality = "Híbrido"
)

const (
	LevelJunior      Level = "Júnior"
	LevelPleno       Level = "Pleno"
	LevelSenior      Level = "Sênior"
	LevelEspecialista Level = "Especialista"
	LevelEstagiario  Level = "Estagiário / Trainee"
)

type Opportunity struct {
	ID             string       `json:"id"`
	ExternalID     string       `json:"external_id"`
	Source         Source       `json:"source"`
	Company        string       `json:"company"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	ContractType   ContractType `json:"contract_type"`
	Modality       Modality     `json:"modality"`
	Level          Level        `json:"level"`
	ServiceType    string       `json:"service_type"`
	Location       string       `json:"location"`
	SalaryRange    string       `json:"salary_range"`
	ApplicationURL string       `json:"application_url"`
	PostedAt       time.Time    `json:"posted_at"`
	IsActive       bool         `json:"is_active"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
