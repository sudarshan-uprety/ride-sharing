package utils

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// FormatValidatorError formats validator.ValidationErrors into a user-friendly map
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
		// Handle non-validator errors (e.g., JSON parsing errors)
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
