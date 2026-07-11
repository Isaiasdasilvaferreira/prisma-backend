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
	s.scrapers[opportunity.SourceLever] = NewLeverScraper(supabase)
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

func (b *BaseScraper) DetermineLevel(title string) opportunity.Level {
	titleLower := strings.ToLower(title)

	if strings.Contains(titleLower, "senior") || strings.Contains(titleLower, "sênior") {
		return opportunity.LevelSenior
	}
	if strings.Contains(titleLower, "pleno") || strings.Contains(titleLower, "mid") {
		return opportunity.LevelPleno
	}
	if strings.Contains(titleLower, "junior") || strings.Contains(titleLower, "júnior") {
		return opportunity.LevelJunior
	}
	if strings.Contains(titleLower, "especialista") || strings.Contains(titleLower, "specialist") {
		return opportunity.LevelEspecialista
	}
	if strings.Contains(titleLower, "estagiario") || strings.Contains(titleLower, "trainee") {
		return opportunity.LevelEstagiario
	}
	return ""
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
	companies := []string{"company1", "company2"}
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
			postedAt, _ := time.Parse(time.RFC3339, job.CreatedAt)

			opp := opportunity.Opportunity{
				ExternalID:     fmt.Sprintf("greenhouse-%d", job.ID),
				Source:         opportunity.SourceGreenhouse,
				Company:        company,
				Title:          job.Title,
				ContractType:   g.DetermineContractType(job.Title),
				Modality:       g.DetermineModality(job.Title, job.Location.Name),
				Level:          g.DetermineLevel(job.Title),
				ServiceType:    g.DetermineServiceType(job.Title),
				Location:       job.Location.Name,
				ApplicationURL: job.AbsoluteURL,
				PostedAt:       postedAt,
				IsActive:       true,
			}
			allOpps = append(allOpps, opp)
		}
	}

	return allOpps, nil
}

type LeverScraper struct {
	*BaseScraper
}

func NewLeverScraper(supabase *supabase.Client) *LeverScraper {
	return &LeverScraper{
		BaseScraper: NewBaseScraper(supabase, "https://api.lever.co/v0"),
	}
}

func (l *LeverScraper) GetSource() opportunity.Source {
	return opportunity.SourceLever
}

func (l *LeverScraper) Scrape(ctx context.Context) ([]opportunity.Opportunity, error) {
	companies := []string{"company1", "company2"}
	var allOpps []opportunity.Opportunity

	for _, company := range companies {
		url := fmt.Sprintf("%s/postings/%s", l.baseURL, company)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			log.Error().Err(err).Str("company", company).Msg("Failed to create request")
			continue
		}

		resp, err := l.client.Do(req)
		if err != nil {
			log.Error().Err(err).Str("company", company).Msg("Failed to fetch jobs")
			continue
		}
		defer resp.Body.Close()

		var result struct {
			Data []struct {
				ID         string `json:"id"`
				Text       string `json:"text"`
				CreatedAt  int64  `json:"createdAt"`
				Categories struct {
					Location   string `json:"location"`
					Commitment string `json:"commitment"`
				} `json:"categories"`
				URL string `json:"url"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Error().Err(err).Str("company", company).Msg("Failed to decode")
			continue
		}

		for _, posting := range result.Data {
			opp := opportunity.Opportunity{
				ExternalID:     fmt.Sprintf("lever-%s", posting.ID),
				Source:         opportunity.SourceLever,
				Company:        company,
				Title:          posting.Text,
				ContractType:   l.DetermineContractType(posting.Text),
				Modality:       l.DetermineModality(posting.Text, posting.Categories.Location),
				Level:          l.DetermineLevel(posting.Text),
				ServiceType:    l.DetermineServiceType(posting.Text),
				Location:       posting.Categories.Location,
				ApplicationURL: posting.URL,
				PostedAt:       time.Unix(posting.CreatedAt/1000, 0),
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
		BaseScraper: NewBaseScraper(supabase, "https://api.ashbyhq.com/posting-api"),
	}
}

func (a *AshbyScraper) GetSource() opportunity.Source {
	return opportunity.SourceAshby
}

func (a *AshbyScraper) Scrape(ctx context.Context) ([]opportunity.Opportunity, error) {
	return []opportunity.Opportunity{}, nil
}

func (s *ScraperService) RunScraping(ctx context.Context) error {
	log.Info().Msg("Starting scraping...")

	for source, scraper := range s.scrapers {
		log.Info().Str("source", string(source)).Msg("Scraping source")
		opps, err := scraper.Scrape(ctx)

		if err != nil {
			log.Error().Err(err).Str("source", string(source)).Msg("Scraping failed")
			continue
		}

		if err := s.saveOpportunities(ctx, opps); err != nil {
			log.Error().Err(err).Str("source", string(source)).Msg("Failed to save")
			continue
		}

		log.Info().Str("source", string(source)).Int("count", len(opps)).Msg("Scraping completed")
	}

	return nil
}

func (s *ScraperService) saveOpportunities(ctx context.Context, opps []opportunity.Opportunity) error {
	for _, opp := range opps {
		existing, err := s.oppRepo.GetByExternalID(ctx, opp.ExternalID)
		if err != nil {
			return err
		}

		if existing == nil {
			if err := s.oppRepo.Create(ctx, &opp); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ScraperService) ScrapeForUser(ctx context.Context, userID uuid.UUID, source opportunity.Source, limit int) ([]opportunity.Opportunity, error) {
	canScrape, remaining, err := s.userSvc.CanScrapeOpportunities(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error checking scrape permissions: %w", err)
	}

	if !canScrape {
		return nil, fmt.Errorf("daily limit reached. Remaining: %d", remaining)
	}

	scraper, exists := s.scrapers[source]
	if !exists {
		return nil, fmt.Errorf("scraper not found for source: %s", source)
	}

	allOpps, err := scraper.Scrape(ctx)
	if err != nil {
		return nil, fmt.Errorf("error scraping %s: %w", source, err)
	}

	if limit > 0 && len(allOpps) > limit {
		allOpps = allOpps[:limit]
	}

	for range allOpps {
		if err := s.userSvc.IncrementUsedCount(ctx, userID); err != nil {
			log.Error().Err(err).Msg("Failed to increment usage count")
		}
	}

	return allOpps, nil
}
