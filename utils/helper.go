package utils

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// HandleRequestErrors processes binding and validation errors
func HandleRequestErrors(c *gin.Context, err error) {
	// Case 1: JSON Unmarshal Type Errors
	if unmarshalTypeError, ok := err.(*json.UnmarshalTypeError); ok {
		fieldName := toSnakeCase(unmarshalTypeError.Field)
		ErrorResponse(c, NewErrorResponse(
			http.StatusBadRequest,
			"Invalid data type",
			fieldName+" must be a "+unmarshalTypeError.Type.String(),
			nil,
		))
		return
	}

	// Case 2: General JSON syntax errors
	if _, ok := err.(*json.SyntaxError); ok {
		ErrorResponse(c, NewErrorResponse(
			http.StatusBadRequest,
			"Invalid JSON format",
			"The request body contains invalid JSON",
			nil,
		))
		return
	}

	// Case 3: String pattern matching for other JSON errors
	if strings.Contains(err.Error(), "cannot unmarshal") {
		fieldName := extractFieldFromUnmarshalError(err.Error())
		ErrorResponse(c, NewErrorResponse(
			http.StatusBadRequest,
			"Invalid data type",
			fieldName+" has an invalid type",
			nil,
		))
		return
	}

	// Case 4: Standard validation errors
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errorMessages := FormatValidatorError(validationErrors)
		ErrorResponse(c, NewErrorResponse(
			http.StatusBadRequest,
			"Validation failed",
			"",
			errorMessages,
		))
		return
	}

	// Default case for other errors
	ErrorResponse(c, NewErrorResponse(
		http.StatusBadRequest,
		"Invalid request",
		err.Error(),
		nil,
	))
}

// FormatValidatorError converts validator errors to a user-friendly map
func FormatValidatorError(err validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)

	for _, e := range err {
		jsonName := toSnakeCase(e.Field())
		tag := e.Tag()
		param := e.Param()

		switch tag {
		case "required":
			errors[jsonName] = jsonName + " is required."
		case "min":
			errors[jsonName] = jsonName + " must be at least " + param + " characters."
		case "max":
			errors[jsonName] = jsonName + " must be less than " + param + " characters."
		case "email":
			errors[jsonName] = jsonName + " must be a valid email address."
		case "eqfield":
			targetField := toSnakeCase(param)
			errors[jsonName] = jsonName + " must match " + targetField + "."
		default:
			errors[jsonName] = jsonName + " is invalid (" + tag + ")."
		}
	}

	return errors
}

// Helper functions remain the same
func extractFieldFromUnmarshalError(errStr string) string {
	re := regexp.MustCompile(`field\s+\w+\.(\w+)`)
	matches := re.FindStringSubmatch(errStr)
	if len(matches) > 1 {
		return toSnakeCase(matches[1])
	}
	return "field"
}

func toSnakeCase(camel string) string {
	var result strings.Builder
	for i, char := range camel {
		if i > 0 && unicode.IsUpper(char) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(char))
	}
	return result.String()
}
