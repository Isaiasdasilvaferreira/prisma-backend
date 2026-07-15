package user_opportunity

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Isaiasdasilvaferreira/prisma-backend/internal/auth"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) CreateUserOpportunity(w http.ResponseWriter, r *http.Request) {
	var req CreateUserOpportunityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	opp, err := c.service.CreateUserOpportunity(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(opp.ToResponse())
}

func (c *Controller) GetUserOpportunity(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/user-opportunities/")
	parts := strings.Split(path, "/")
	id := parts[0]

	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	opp, err := c.service.GetUserOpportunity(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(opp.ToResponse())
}

func (c *Controller) GetAllUserOpportunities(w http.ResponseWriter, r *http.Request) {
	var isActive *bool
	if isActiveParam := r.URL.Query().Get("is_active"); isActiveParam != "" {
		val, err := strconv.ParseBool(isActiveParam)
		if err == nil {
			isActive = &val
		}
	}

	opportunities, err := c.service.GetAllUserOpportunities(r.Context(), isActive)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(opportunities) == 0 {
		json.NewEncoder(w).Encode([]UserOpportunityResponse{})
		return
	}
	json.NewEncoder(w).Encode(ToResponseList(opportunities))
}

func (c *Controller) UpdateUserOpportunity(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/user-opportunities/")
	parts := strings.Split(path, "/")
	id := parts[0]

	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req UpdateUserOpportunityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	opp, err := c.service.UpdateUserOpportunity(r.Context(), id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(opp.ToResponse())
}

func (c *Controller) DeleteUserOpportunity(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/user-opportunities/")
	parts := strings.Split(path, "/")
	id := parts[0]

	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err := c.service.DeleteUserOpportunity(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) ApproveUserOpportunity(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/user-opportunities/")
	parts := strings.Split(path, "/")
	id := parts[0]

	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err := c.service.ApproveUserOpportunity(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	opp, err := c.service.GetUserOpportunity(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(opp.ToResponse())
}

func (c *Controller) RejectUserOpportunity(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/user-opportunities/")
	parts := strings.Split(path, "/")
	id := parts[0]

	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err := c.service.RejectUserOpportunity(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	opp, err := c.service.GetUserOpportunity(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(opp.ToResponse())
}

func (c *Controller) ApplyToOpportunity(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not authenticated"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/user-opportunities/")
	parts := strings.Split(path, "/")
	id := parts[0]

	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	opp, err := c.service.ApplyToOpportunity(r.Context(), id, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if opp != nil {
		json.NewEncoder(w).Encode(opp.ToResponse())
	} else {
		json.NewEncoder(w).Encode(map[string]string{"message": "Application successful"})
	}
}

func (c *Controller) GetUserApplications(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	opps, err := c.service.GetUserApplications(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToResponseList(opps))
}
