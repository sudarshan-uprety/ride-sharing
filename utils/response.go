package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Warning string      `json:"warning,omitempty"`
}

// Success returns a successful JSON response
func Success(w http.ResponseWriter, statusCode int, message string, data interface{}, warning string) {
	resp := Response{
		Message: message,
		Success: true,
		Data:    data,
		Warning: warning,
	}
	writeJSON(w, statusCode, resp)
}

// Error returns an error JSON response
func Error(w http.ResponseWriter, statusCode int, message string, errors interface{}) {
	resp := Response{
		Message: message,
		Success: false,
		Errors:  errors,
	}
	writeJSON(w, statusCode, resp)
}

// writeJSON handles encoding and header writing
func writeJSON(w http.ResponseWriter, statusCode int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
