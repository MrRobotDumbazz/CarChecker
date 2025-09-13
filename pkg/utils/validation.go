package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func ValidateUUID(uuidStr string) error {
	if uuidStr == "" {
		return fmt.Errorf("UUID cannot be empty")
	}

	_, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID format: %s", uuidStr)
	}

	return nil
}

func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

func ValidateStringLength(value, fieldName string, minLen, maxLen int) error {
	length := len(strings.TrimSpace(value))

	if minLen > 0 && length < minLen {
		return fmt.Errorf("%s must be at least %d characters long", fieldName, minLen)
	}

	if maxLen > 0 && length > maxLen {
		return fmt.Errorf("%s must be no more than %d characters long", fieldName, maxLen)
	}

	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func ValidateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL is required")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must include a scheme (http or https)")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must include a host")
	}

	return nil
}

func ValidateNumericRange(value int, fieldName string, min, max int) error {
	if value < min {
		return fmt.Errorf("%s must be at least %d", fieldName, min)
	}

	if value > max {
		return fmt.Errorf("%s must be no more than %d", fieldName, max)
	}

	return nil
}

func ValidateEnum(value, fieldName string, allowedValues []string) error {
	if value == "" {
		return fmt.Errorf("%s is required", fieldName)
	}

	for _, allowed := range allowedValues {
		if strings.EqualFold(value, allowed) {
			return nil
		}
	}

	return fmt.Errorf("%s must be one of: %s", fieldName, strings.Join(allowedValues, ", "))
}

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#x27;")
	input = strings.ReplaceAll(input, "&", "&amp;")
	return input
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, "; ")
}

func (ve *ValidationErrors) Add(field, message string) {
	*ve = append(*ve, ValidationError{Field: field, Message: message})
}

func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

func (ve ValidationErrors) ToMap() map[string]string {
	result := make(map[string]string)
	for _, err := range ve {
		result[err.Field] = err.Message
	}
	return result
}