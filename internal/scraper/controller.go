package scraper

import (
	"net/http"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/database"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
)

type ScraperController struct {
	service *ScraperService
}

func NewScraperController(db *database.Database) *ScraperController {
	return &ScraperController{
		service: NewScraperService(db),
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
