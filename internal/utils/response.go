package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func SuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	RespondWithJSON(w, status, Response{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	RespondWithJSON(w, status, Response{
		Success: false,
		Error:   message,
	})
}
