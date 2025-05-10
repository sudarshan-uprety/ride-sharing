package validation

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// ProcessValidationError creates a structured error response from validator errors
func ProcessValidationError(err error) map[string]string {
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
				errors[jsonName] = "Required field"
			case "min":
				errors[jsonName] = "Must be at least " + param + " characters"
			case "max":
				errors[jsonName] = "Must be less than " + param + " characters"
			case "email":
				errors[jsonName] = "Must be a valid email address"
			case "eqfield":
				targetField := toSnakeCase(param)
				errors[jsonName] = "Must match " + targetField
			case "numeric":
				errors[jsonName] = "Must be a number"
			case "alphanum":
				errors[jsonName] = "Must contain only letters and numbers"
			case "e164":
				errors[jsonName] = "Must be a valid phone number in E.164 format"
			case "strongpassword":
				errors[jsonName] = GetPasswordRules()
			case "otpvalidation":
				errors[jsonName] = GetOTPRules()
			default:
				errors[jsonName] = "Invalid value (" + tag + ")"
			}
		}
	} else if unmarshalTypeError, ok := err.(*json.UnmarshalTypeError); ok {
		// JSON Unmarshal Type Errors
		fieldName := toSnakeCase(unmarshalTypeError.Field)
		errors[fieldName] = "Must be a " + unmarshalTypeError.Type.String()
	} else if syntaxError, ok := err.(*json.SyntaxError); ok {
		// JSON Syntax errors
		errors["_error"] = "Invalid JSON format at position " + string(rune(syntaxError.Offset))
	} else if strings.Contains(err.Error(), "cannot unmarshal") {
		// Other JSON unmarshal errors
		fieldName := extractFieldFromUnmarshalError(err.Error())
		errors[fieldName] = "Has an invalid type"
	} else {
		// Handle any other errors
		errors["_error"] = err.Error()
	}

	return errors
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

// RegisterCustomValidators adds all custom validators to the provided validator instance
func RegisterCustomValidators(v *validator.Validate) error {
	// Register strongPassword validator with proper error handling
	if err := v.RegisterValidation("strongpassword", validateStrongPassword); err != nil {
		return err
	}
	if err := v.RegisterValidation("otpvalidation", validateOTP); err != nil {
		return err
	}
	return nil
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Optimized checks using regex for better performance
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	return len(password) >= 8 && hasUpper && hasLower && hasDigit && hasSpecial
}

func validateOTP(fl validator.FieldLevel) bool {
	otp := fl.Field().String()

	// Regex: ^\d{6}$ means exactly 6 digits
	match := regexp.MustCompile(`^\d{6}$`).MatchString(otp)

	return match
}

// GetPasswordRules returns a description of password requirements for API docs/errors
func GetPasswordRules() string {
	return "Password must contain at least 8 characters including: 1 uppercase, 1 lowercase, 1 digit, and 1 special character"
}

func GetOTPRules() string {
	return "OTP must br 6 digit"
}
