package scraper

import (
	"net/http"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
	"github.com/nedpals/supabase-go"
)

type ScraperController struct {
	service *ScraperService
}

func NewScraperController(supabase *supabase.Client) *ScraperController {
	return &ScraperController{
		service: NewScraperService(supabase),
	}
}

func (c *ScraperController) TriggerScraping(w http.ResponseWriter, r *http.Request) {
	userRole := r.Context().Value("user_role")
	if userRole != "admin" {
		utils.RespondWithError(w, http.StatusForbidden, "Apenas administradores")
		return
	}

	go func() {
		ctx := r.Context()
		c.service.RunScraping(ctx)
	}()

	utils.RespondWithJSON(w, http.StatusAccepted, map[string]string{
		"message": "Scraping iniciado",
	})
}
