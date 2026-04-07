package utils

import (
	"teras-vps/backend/models"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// RegisterInput represents user registration input
type RegisterInput struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Phone    string `json:"phone" validate:"omitempty,max=20"`
}

// LoginInput represents user login input
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateProfileInput represents profile update input
type UpdateProfileInput struct {
	Phone string `json:"phone" validate:"omitempty,max=20"`
}

// ChangePasswordInput represents password change input
type ChangePasswordInput struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}

// Validate validates input struct
func Validate(s interface{}) error {
	return validate.Struct(s)
}

// ValidateUsername checks if username is valid
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	// Check if username is alphanumeric
	for _, c := range username {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}

// ValidateEmail checks if email is valid
func ValidateEmail(email string) bool {
	err := validate.Var(email, "required,email")
	return err == nil
}

// ValidatePassword checks if password meets requirements
func ValidatePassword(password string) bool {
	if len(password) < 8 || len(password) > 100 {
		return false
	}
	return true
}

// ValidateRole checks if role is valid
func ValidateRole(role string) bool {
	return role == string(models.RoleCustomer) || role == string(models.RoleAdmin)
}
