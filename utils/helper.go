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

// HandleRequestErrors processes binding and validation errors and returns appropriate responses
func HandleRequestErrors(c *gin.Context, err error) bool {
	// If there's an error, process it and return true
	if err != nil {
		// Case 1: JSON Unmarshal Type Errors
		if unmarshalTypeError, ok := err.(*json.UnmarshalTypeError); ok {
			// Get the field name directly from the error
			fieldName := unmarshalTypeError.Field
			errorMap := map[string]string{
				fieldName: fieldName + " must be a " + unmarshalTypeError.Type.String() + ".",
			}
			Error(c.Writer, http.StatusBadRequest, "Invalid data type", errorMap)
			return true
		}

		// Case 2: General JSON syntax errors
		if _, ok := err.(*json.SyntaxError); ok {
			Error(c.Writer, http.StatusBadRequest, "Invalid JSON format", map[string]string{
				"error": "The request body contains invalid JSON.",
			})
			return true
		}

		// Case 3: String pattern matching for other JSON errors
		if strings.Contains(err.Error(), "cannot unmarshal") {
			// Extract field name from error message
			fieldName := extractFieldFromUnmarshalError(err.Error())
			errorMap := map[string]string{
				fieldName: fieldName + " has an invalid type.",
			}
			Error(c.Writer, http.StatusBadRequest, "Invalid data type", errorMap)
			return true
		}

		// Case 4: Standard validation errors
		errorMessages := FormatValidatorError(err)
		Error(c.Writer, http.StatusBadRequest, "Validation failed", errorMessages)
		return true
	}
	// No error
	return false
}

// Helper function to extract field name from unmarshal errors
func extractFieldFromUnmarshalError(errStr string) string {
	// Error format typically: "json: cannot unmarshal X into Go struct field StructName.field_name of type Y"
	re := regexp.MustCompile(`field\s+\w+\.(\w+)`) // This matches "field StructName.fieldName"
	matches := re.FindStringSubmatch(errStr)
	if len(matches) > 1 {
		// Convert to snake_case for JSON field name
		return toSnakeCase(matches[1])
	}
	return "error"
}

// FormatValidatorError converts validator errors to a user-friendly map
func FormatValidatorError(err error) map[string]string {
	errors := make(map[string]string)

	// Check if the error is from validator
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			// Get the JSON field name (snake_case)
			jsonName := toSnakeCase(e.Field())

			tag := e.Tag()     // e.g., "required", "min", "email"
			param := e.Param() // e.g., "8" for min=8

			// Use the JSON field name (e.g., "full_name" instead of "FullName")
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
			case "numeric":
				errors[jsonName] = jsonName + " must be a number."
			case "alphanum":
				errors[jsonName] = jsonName + " must contain only letters and numbers."
			default:
				errors[jsonName] = jsonName + " is invalid (" + tag + ")."
			}
		}
	} else {
		// Handle non-validator errors
		errors["error"] = err.Error()
	}

	return errors
}

// Helper function to convert CamelCase to snake_case
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
