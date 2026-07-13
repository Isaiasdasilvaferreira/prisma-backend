package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/opportunity"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/user"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
	"github.com/rs/zerolog/log"
)

type Scraper interface {
	Scrape(ctx context.Context) ([]opportunity.Opportunity, error)
	GetSource() opportunity.Source
}

type ScraperService struct {
	supabase   *supabase.Client
	userSvc    user.Service
	oppRepo    opportunity.Repository
	scrapers   map[opportunity.Source]Scraper
}

func NewScraperService(supabase *supabase.Client, userSvc user.Service, oppRepo opportunity.Repository) *ScraperService {
	s := &ScraperService{
		supabase: supabase,
		userSvc:  userSvc,
		oppRepo:  oppRepo,
		scrapers: make(map[opportunity.Source]Scraper),
	}

	s.scrapers[opportunity.SourceGreenhouse] = NewGreenhouseScraper(supabase)
	s.scrapers[opportunity.SourceAshby] = NewAshbyScraper(supabase)

	return s
}

type BaseScraper struct {
	client   *http.Client
	baseURL  string
	supabase *supabase.Client
}

func NewBaseScraper(supabase *supabase.Client, baseURL string) *BaseScraper {
	return &BaseScraper{
		client:   &http.Client{Timeout: 30 * time.Second},
		baseURL:  baseURL,
		supabase: supabase,
	}
}

func (b *BaseScraper) DetermineContractType(title string) opportunity.ContractType {
	titleLower := strings.ToLower(title)
	if strings.Contains(titleLower, "freelance") || strings.Contains(titleLower, "freelancer") {
		return opportunity.ContractFreelancer
	}
	return opportunity.ContractCLT
}

func (b *BaseScraper) DetermineModality(title, location string) opportunity.Modality {
	titleLower := strings.ToLower(title)
	locationLower := strings.ToLower(location)

	if strings.Contains(titleLower, "remote") || strings.Contains(titleLower, "remoto") {
		return opportunity.ModalityRemoto
	}
	if strings.Contains(titleLower, "hybrid") || strings.Contains(titleLower, "híbrido") {
		return opportunity.ModalityHibrido
	}
	if strings.Contains(locationLower, "remote") || strings.Contains(locationLower, "remoto") {
		return opportunity.ModalityRemoto
	}
	return opportunity.ModalityPresencial
}

func (b *BaseScraper) DetermineServiceType(title string) string {
	titleLower := strings.ToLower(title)

	services := map[string]string{
		"ui":             "UI Design",
		"ux":             "UX Design",
		"product design": "Product Design",
		"branding":       "Branding / Identidade Visual",
		"motion":         "Motion Design",
		"graphic":        "Design Gráfico (generalista)",
		"editorial":      "Design Editorial",
		"packaging":      "Packaging",
		"social media":   "Social Media Design",
		"identity":       "Criação de Identidade Visual / Branding",
		"landing page":   "Landing Page Design",
		"presentation":   "Apresentações",
	}

	for key, value := range services {
		if strings.Contains(titleLower, key) {
			return value
		}
	}
	return "Design"
}

func (b *BaseScraper) IsDesignRelated(title string) bool {
	titleLower := strings.ToLower(title)

	designKeywords := []string{
		"design", "designer",
		"ui", "ux",
		"product design", "product designer",
		"graphic", "graphic designer",
		"visual", "visual designer",
		"branding", "brand designer",
		"motion", "motion designer",
		"editorial", "editorial designer",
		"packaging", "packaging designer",
		"social media", "social media designer",
		"identidade visual",
		"landing page", "landing page designer",
		"web design", "web designer",
		"interaction", "interaction designer",
		"user interface", "user interface designer",
		"user experience", "user experience designer",
		"creative", "creative designer",
		"art director",
		"design gráfico",
		"ilustração", "illustration",
		"ux design", "ui design",
	}

	for _, keyword := range designKeywords {
		if strings.Contains(titleLower, keyword) {
			return true
		}
	}
	return false
}

func (b *BaseScraper) IsExcludedRole(title string) bool {
	titleLower := strings.ToLower(title)

	excludedKeywords := []string{
		"account", "account manager", "account executive",
		"analista de contas", "analista de reclame aqui",
		"sales", "vendas", "comercial",
		"farmer", "farming",
		"reclame aqui", "escalados",
		"customer success", "sucesso do cliente",
		"marketing", "mercado",
		"finance", "finanças",
		"hr", "rh", "recursos humanos",
		"admin", "administrativo",
		"legal", "jurídico",
		"operation", "operações",
		"support", "suporte",
		"it", "ti", "tecnologia",
	}

	for _, keyword := range excludedKeywords {
		if strings.Contains(titleLower, keyword) {
			return true
		}
	}
	return false
}

func IsLocationInBrazil(location string) bool {
	locationLower := strings.ToLower(location)
	return strings.Contains(locationLower, "brasil") ||
		strings.Contains(locationLower, "brazil") ||
		strings.Contains(locationLower, "são paulo") ||
		strings.Contains(locationLower, "são paulo") ||
		strings.Contains(locationLower, "rio de janeiro") ||
		strings.Contains(locationLower, "rio") ||
		strings.Contains(locationLower, "sp") ||
		strings.Contains(locationLower, "rj")
}

func IsRemote(location string) bool {
	locationLower := strings.ToLower(location)
	return strings.Contains(locationLower, "remote") ||
		strings.Contains(locationLower, "remoto")
}

type GreenhouseScraper struct {
	*BaseScraper
}

func NewGreenhouseScraper(supabase *supabase.Client) *GreenhouseScraper {
	return &GreenhouseScraper{
		BaseScraper: NewBaseScraper(supabase, "https://boards-api.greenhouse.io/v1/boards"),
	}
}

func (g *GreenhouseScraper) GetSource() opportunity.Source {
	return opportunity.SourceGreenhouse
}

func (g *GreenhouseScraper) Scrape(ctx context.Context) ([]opportunity.Opportunity, error) {
	companies := []string{
		"nubank",
		"figma",
		"notion",
		"uber",
		"airbnb",
		"google",
		"apple",
		"microsoft",
		"amazon",
		"adobe",
		"canva",
		"spotify",
		"netflix",
		"shopify",
		"stripe",
		"dropbox",
		"pinterest",
		"invision",
		"webflow",
		"sketch",
		"marvelapp",
		"protopie",
		"squarespace",
		"wix",
		"godaddy",
	}
	var allOpps []opportunity.Opportunity

	for _, company := range companies {
		url := fmt.Sprintf("%s/%s/jobs", g.baseURL, company)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			log.Error().Err(err).Str("company", company).Msg("Failed to create request")
			continue
		}

		resp, err := g.client.Do(req)
		if err != nil {
			log.Error().Err(err).Str("company", company).Msg("Failed to fetch jobs")
			continue
		}
		defer resp.Body.Close()

		var result struct {
			Jobs []struct {
				ID          int                   `json:"id"`
				Title       string                `json:"title"`
				Location    struct{ Name string } `json:"location"`
				AbsoluteURL string                `json:"absolute_url"`
				CreatedAt   string                `json:"created_at"`
			} `json:"jobs"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Error().Err(err).Str("company", company).Msg("Failed to decode")
			continue
		}

		for _, job := range result.Jobs {
			if !g.IsDesignRelated(job.Title) {
				continue
			}

			if g.IsExcludedRole(job.Title) {
				continue
			}

			opp := opportunity.Opportunity{
				ExternalID:     fmt.Sprintf("greenhouse-%d", job.ID),
				Source:         opportunity.SourceGreenhouse,
				Company:        company,
				Title:          job.Title,
				ContractType:   g.DetermineContractType(job.Title),
				Modality:       g.DetermineModality(job.Title, job.Location.Name),
				ServiceType:    g.DetermineServiceType(job.Title),
				Location:       job.Location.Name,
				ApplicationURL: job.AbsoluteURL,
				IsActive:       true,
			}
			allOpps = append(allOpps, opp)
		}
	}

	return allOpps, nil
}

type AshbyScraper struct {
	*BaseScraper
}

func NewAshbyScraper(supabase *supabase.Client) *AshbyScraper {
	return &AshbyScraper{
		BaseScraper: NewBaseScraper(supabase, "https://api.ashbyhq.com/posting-api/job-board"),
	}
}

func (a *AshbyScraper) GetSource() opportunity.Source {
	return opportunity.SourceAshby
}

func (a *AshbyScraper) Scrape(ctx context.Context) ([]opportunity.Opportunity, error) {
	companies := []struct {
		Name string
		Slug string
	}{
		{"cursor", "cursor"},
		{"notion", "notion"},
	}

	var allOpps []opportunity.Opportunity

	for _, company := range companies {
		url := fmt.Sprintf("%s/%s?includeCompensation=true", a.baseURL, company.Slug)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			log.Error().Err(err).Str("company", company.Name).Msg("Failed to create request")
			continue
		}

		req.Header.Set("Accept", "application/json; version=1")
		req.Header.Set("Content-Type", "application/json")

		resp, err := a.client.Do(req)
		if err != nil {
			log.Error().Err(err).Str("company", company.Name).Msg("Failed to fetch jobs")
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 401 || resp.StatusCode == 403 {
			log.Warn().Str("company", company.Name).Int("status", resp.StatusCode).Msg("Unauthorized or forbidden")
			continue
		}

		if resp.StatusCode == 404 {
			log.Warn().Str("company", company.Name).Msg("Job board not found")
			continue
		}

		var result struct {
			Jobs []struct {
				ID          string `json:"id"`
				Title       string `json:"title"`
				Description string `json:"description"`
				Location    struct {
					City        string `json:"city"`
					Country     string `json:"country"`
					Remote      bool   `json:"remote"`
					Hybrid      bool   `json:"hybrid"`
					DisplayName string `json:"displayName"`
				} `json:"location"`
				URL struct {
					Application string `json:"application"`
				} `json:"url"`
				ListedAt     int64  `json:"listedAt"`
				IsListed     bool   `json:"isListed"`
				Department   string `json:"department"`
				Employment   string `json:"employment"`
				Compensation struct {
					Currency string `json:"currency"`
					Min      int    `json:"min"`
					Max      int    `json:"max"`
				} `json:"compensation"`
			} `json:"jobs"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Error().Err(err).Str("company", company.Name).Msg("Failed to decode")
			continue
		}

		for _, job := range result.Jobs {
			if !job.IsListed {
				continue
			}

			if !a.IsDesignRelated(job.Title) {
				continue
			}

			if a.IsExcludedRole(job.Title) {
				continue
			}

			locationDisplay := job.Location.DisplayName
			if !IsLocationInBrazil(locationDisplay) && !job.Location.Remote && !job.Location.Hybrid {
				continue
			}

			var modality opportunity.Modality
			if job.Location.Remote {
				modality = opportunity.ModalityRemoto
			} else if job.Location.Hybrid {
				modality = opportunity.ModalityHibrido
			} else {
				modality = opportunity.ModalityPresencial
			}

			if locationDisplay == "" {
				locationDisplay = "Remote"
			}

			opp := opportunity.Opportunity{
				ExternalID:     fmt.Sprintf("ashby-%s", job.ID),
				Source:         opportunity.SourceAshby,
				Company:        company.Name,
				Title:          job.Title,
				ContractType:   a.DetermineContractType(job.Title),
				Modality:       modality,
				ServiceType:    a.DetermineServiceType(job.Title),
				Location:       locationDisplay,
				ApplicationURL: job.URL.Application,
				IsActive:       true,
			}
			allOpps = append(allOpps, opp)
		}
	}

	return allOpps, nil
}

func (s *ScraperService) RunScraping(ctx context.Context) error {
	utils.LogInfo("RunScraping iniciado")
	log.Info().Msg("Starting scraping...")

	for source, scraper := range s.scrapers {
		utils.LogInfo(fmt.Sprintf("Raspando fonte: %s", source))
		log.Info().Str("source", string(source)).Msg("Scraping source")
		opps, err := scraper.Scrape(ctx)

		if err != nil {
			utils.LogError(fmt.Sprintf("Scraping falhou para %s", source), err)
			log.Error().Err(err).Str("source", string(source)).Msg("Scraping failed")
			continue
		}

		if err := s.saveOpportunities(ctx, opps); err != nil {
			utils.LogError(fmt.Sprintf("Failed to save %s", source), err)
			log.Error().Err(err).Str("source", string(source)).Msg("Failed to save")
			continue
		}

		utils.LogInfo(fmt.Sprintf("Fonte %s retornou %d oportunidades", source, len(opps)))
		log.Info().Str("source", string(source)).Int("count", len(opps)).Msg("Scraping completed")
	}

	utils.LogInfo("RunScraping finalizado")
	return nil
}

func (s *ScraperService) saveOpportunities(ctx context.Context, opps []opportunity.Opportunity) error {
	utils.LogInfo(fmt.Sprintf("saveOpportunities chamado com %d oportunidades", len(opps)))

	var brazilOpps []opportunity.Opportunity
	var internationalOpps []opportunity.Opportunity

	for _, opp := range opps {
		existing, err := s.oppRepo.GetByExternalID(ctx, opp.ExternalID)
		if err != nil {
			utils.LogError(fmt.Sprintf("Erro ao verificar existência de %s", opp.ExternalID), err)
			return err
		}

		if existing != nil {
			continue
		}

		if IsLocationInBrazil(opp.Location) {
			brazilOpps = append(brazilOpps, opp)
		} else {
			internationalOpps = append(internationalOpps, opp)
		}
	}

	var finalOpps []opportunity.Opportunity

	if len(brazilOpps) > 0 {
		finalOpps = brazilOpps
		utils.LogInfo(fmt.Sprintf("Adicionando %d vagas do Brasil primeiro", len(brazilOpps)))
	}

	if len(finalOpps) < 10 && len(internationalOpps) > 0 {
		needed := 10 - len(finalOpps)
		if len(internationalOpps) > needed {
			finalOpps = append(finalOpps, internationalOpps[:needed]...)
			utils.LogInfo(fmt.Sprintf("Adicionando %d vagas internacionais para completar 10", needed))
		} else {
			finalOpps = append(finalOpps, internationalOpps...)
			utils.LogInfo(fmt.Sprintf("Adicionando %d vagas internacionais (total disponível)", len(internationalOpps)))
		}
	}

	companyCount := make(map[string]int)
	var filteredOpps []opportunity.Opportunity

	for _, opp := range finalOpps {
		if companyCount[opp.Company] >= 2 {
			continue
		}
		companyCount[opp.Company]++
		filteredOpps = append(filteredOpps, opp)
	}

	if len(filteredOpps) > 10 {
		filteredOpps = filteredOpps[:10]
	}

	utils.LogInfo(fmt.Sprintf("Após filtro: %d oportunidades (máx 2 por empresa, máximo 10)", len(filteredOpps)))

	for _, opp := range filteredOpps {
		utils.LogInfo(fmt.Sprintf("Criando nova oportunidade: %s", opp.ExternalID))
		if err := s.oppRepo.Create(ctx, &opp); err != nil {
			utils.LogError(fmt.Sprintf("Erro ao criar %s", opp.ExternalID), err)
			return err
		}
		utils.LogInfo(fmt.Sprintf("Oportunidade criada com sucesso: %s", opp.ExternalID))
	}

	utils.LogInfo("saveOpportunities finalizado")
	return nil
}

func (s *ScraperService) ScrapeForUser(ctx context.Context, userID uuid.UUID, source opportunity.Source, limit int) ([]opportunity.Opportunity, error) {
	utils.LogInfo(fmt.Sprintf("ScrapeForUser - UserID: %s, Source: %s, Limit: %d", userID.String(), source, limit))

	canScrape, remaining, err := s.userSvc.CanScrapeOpportunities(ctx, userID)
	if err != nil {
		utils.LogError("Erro ao verificar permissões de scraping", err)
		return nil, fmt.Errorf("error checking scrape permissions: %w", err)
	}

	if !canScrape {
		utils.LogInfo(fmt.Sprintf("Limite diário atingido. Restam: %d", remaining))
		return nil, fmt.Errorf("daily limit reached. Remaining: %d", remaining)
	}

	scraper, exists := s.scrapers[source]
	if !exists {
		utils.LogError(fmt.Sprintf("Scraper não encontrado: %s", source), nil)
		return nil, fmt.Errorf("scraper not found for source: %s", source)
	}

	allOpps, err := scraper.Scrape(ctx)
	if err != nil {
		utils.LogError(fmt.Sprintf("Erro ao raspar %s", source), err)
		return nil, fmt.Errorf("error scraping %s: %w", source, err)
	}

	if limit > 0 && len(allOpps) > limit {
		allOpps = allOpps[:limit]
	}

	if len(allOpps) > 0 {
		if err := s.userSvc.IncrementUsedCount(ctx, userID); err != nil {
			utils.LogError("Erro ao incrementar contagem de uso", err)
			log.Error().Err(err).Msg("Failed to increment usage count")
		}
	}

	utils.LogInfo(fmt.Sprintf("ScrapeForUser retornando %d oportunidades", len(allOpps)))
	return allOpps, nil
}
