package plans

import (
	"net/http"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/database"
	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/utils"
)

type PlanController struct {
	service *PlanService
}

func NewPlanController(db *database.Database) *PlanController {
	return &PlanController{
		service: NewPlanService(db),
	}
}

func (c *PlanController) CanScrape(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	canScrape, err := c.service.CanScrape(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Erro ao verificar plano")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"can_scrape": canScrape,
	})
}

func (c *PlanController) UpgradeToProfessional(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	if err := c.service.UpgradeToProfessional(r.Context(), userID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Erro ao fazer upgrade")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Plano atualizado para Profissional",
	})
}
