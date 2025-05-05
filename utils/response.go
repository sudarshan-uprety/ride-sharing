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
func Success(message string, data interface{}, warning string, statusCode int) (Response, int) {
	return Response{
		Message: message,
		Success: true,
		Data:    data,
		Warning: warning,
	}, statusCode
}

// Error returns an error JSON response
func Error(message string, errors interface{}, warning string, statusCode int) (Response, int) {
	return Response{
		Message: message,
		Success: false,
		Errors:  errors,
	}, statusCode
}

// writeJSON handles encoding and header writing
func writeJSON(w http.ResponseWriter, statusCode int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
