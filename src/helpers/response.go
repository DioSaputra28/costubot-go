package helpers

import (
	"contact-management/src/apps"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func SuccessResponse(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := Response{
		Status:  "success",
		Message: message,
		Data:    data,
		Error:   nil,
	}
	json.NewEncoder(w).Encode(response)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string, err any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := Response{
		Status:  "error",
		Message: message,
		Data:    nil,
		Error:   err,
	}

	apps.LoggingApp().WithFields(logrus.Fields{
		"status" : statusCode,
		"message" : message,
	}).Info("Success response")

	json.NewEncoder(w).Encode(response)
}

func UnauthorizedResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	response := Response{
		Status:  "error",
		Message: message,
		Data:    nil,
		Error:   "Unauthorized access",
	}
	json.NewEncoder(w).Encode(response)
}

func ForbiddenResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	response := Response{
		Status:  "error",
		Message: message,
		Data:    nil,
		Error:   "Access forbidden",
	}
	json.NewEncoder(w).Encode(response)
}

func NotFoundResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	response := Response{
		Status:  "error",
		Message: message,
		Data:    nil,
		Error:   "Resource not found",
	}
	json.NewEncoder(w).Encode(response)
}

func BadRequestResponse(w http.ResponseWriter, message string, err any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	response := Response{
		Status:  "error",
		Message: message,
		Data:    nil,
		Error:   err,
	}
	json.NewEncoder(w).Encode(response)
}

func ConflictResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	response := Response{
		Status:  "error",
		Message: message,
		Data:    nil,
		Error:   "Conflict occurred",
	}
	json.NewEncoder(w).Encode(response)
}

func InternalServerErrorResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	response := Response{
		Status:  "error",
		Message: message,
		Data:    nil,
		Error:   "Internal server error",
	}
	json.NewEncoder(w).Encode(response)
}

func CreatedResponse(w http.ResponseWriter, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := Response{
		Status:  "success",
		Message: message,
		Data:    data,
		Error:   nil,
	}
	json.NewEncoder(w).Encode(response)
}

func NoContentResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func ValidationErrorResponse(w http.ResponseWriter, message string, validationErrors any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	response := Response{
		Status:  "error",
		Message: message,
		Data:    nil,
		Error:   validationErrors,
	}
	json.NewEncoder(w).Encode(response)
}

func TooManyRequestsResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	response := Response{
		Status:  "error",
		Message: message,
		Data:    nil,
		Error:   "Rate limit exceeded",
	}
	json.NewEncoder(w).Encode(response)
}