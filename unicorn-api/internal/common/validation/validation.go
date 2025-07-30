package validation

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode"

	"unicorn-api/internal/common/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Validator provides validation utilities
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateUUID validates if a string is a valid UUID
func (v *Validator) ValidateUUID(id string) error {
	if id == "" {
		return errors.ErrBadRequest.WithDetails("ID is required")
	}

	if _, err := uuid.Parse(id); err != nil {
		return errors.ErrBadRequest.WithDetails("Invalid UUID format")
	}

	return nil
}

// ValidateEmail validates email format
func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return errors.ErrBadRequest.WithDetails("Email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.ErrBadRequest.WithDetails("Invalid email format")
	}

	return nil
}

// ValidatePassword validates password strength
func (v *Validator) ValidatePassword(password string) error {
	if password == "" {
		return errors.ErrBadRequest.WithDetails("Password is required")
	}

	if len(password) < 8 {
		return errors.ErrBadRequest.WithDetails("Password must be at least 8 characters long")
	}

	var (
		hasUpper  bool
		hasLower  bool
		hasNumber bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return errors.ErrBadRequest.WithDetails("Password must contain uppercase, lowercase, and numeric characters")
	}

	return nil
}

// ValidateName validates name format
func (v *Validator) ValidateName(name string) error {
	if name == "" {
		return errors.ErrBadRequest.WithDetails("Name is required")
	}

	if len(strings.TrimSpace(name)) < 2 {
		return errors.ErrBadRequest.WithDetails("Name must be at least 2 characters long")
	}

	if len(name) > 100 {
		return errors.ErrBadRequest.WithDetails("Name must be less than 100 characters")
	}

	return nil
}

// ValidateSecretName validates secret name format
func (v *Validator) ValidateSecretName(name string) error {
	if name == "" {
		return errors.ErrBadRequest.WithDetails("Secret name is required")
	}

	if len(name) > 50 {
		return errors.ErrBadRequest.WithDetails("Secret name must be less than 50 characters")
	}

	// Only allow alphanumeric characters, hyphens, and underscores
	validNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validNameRegex.MatchString(name) {
		return errors.ErrBadRequest.WithDetails("Secret name can only contain alphanumeric characters, hyphens, and underscores")
	}

	return nil
}

// ValidateSecretValue validates secret value
func (v *Validator) ValidateSecretValue(value string) error {
	if value == "" {
		return errors.ErrBadRequest.WithDetails("Secret value is required")
	}

	if len(value) > 10000 {
		return errors.ErrBadRequest.WithDetails("Secret value must be less than 10KB")
	}

	return nil
}

// ValidateDockerImage validates Docker image name
func (v *Validator) ValidateDockerImage(image string) error {
	if image == "" {
		return errors.ErrBadRequest.WithDetails("Docker image is required")
	}

	// Basic Docker image validation
	imageRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*/[a-zA-Z0-9][a-zA-Z0-9._-]*:[a-zA-Z0-9._-]+$|^[a-zA-Z0-9][a-zA-Z0-9._-]*:[a-zA-Z0-9._-]+$`)
	if !imageRegex.MatchString(image) {
		return errors.ErrBadRequest.WithDetails("Invalid Docker image format")
	}

	return nil
}

// ValidatePort validates port number
func (v *Validator) ValidatePort(port string) error {
	if port == "" {
		return errors.ErrBadRequest.WithDetails("Port is required")
	}

	portRegex := regexp.MustCompile(`^\d+$`)
	if !portRegex.MatchString(port) {
		return errors.ErrBadRequest.WithDetails("Port must be a number")
	}

	return nil
}

// ValidateJSON validates if a string is valid JSON
func (v *Validator) ValidateJSON(jsonStr string) error {
	if jsonStr == "" {
		return nil // Empty JSON is valid
	}

	var js interface{}
	if err := json.Unmarshal([]byte(jsonStr), &js); err != nil {
		return errors.ErrBadRequest.WithDetails("Invalid JSON format")
	}

	return nil
}

// BindAndValidate binds JSON request and validates it
func (v *Validator) BindAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.ErrBadRequest.WithDetails(err.Error())
	}

	return nil
}
