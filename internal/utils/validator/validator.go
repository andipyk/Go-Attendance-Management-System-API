package validator

import (
	"golang-tes/internal/domain"
	"net/mail"
	"strings"
	"unicode"
)

// ValidateEmail checks if the email format is valid
func ValidateEmail(email string) error {
	if len(email) > domain.MaxEmailLength {
		return domain.ErrInvalidEmail
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return domain.ErrInvalidEmail
	}
	return nil
}

// ValidatePassword checks if the password meets security requirements
func ValidatePassword(password string) error {
	if len(password) < domain.MinPasswordLength {
		return domain.ErrInvalidPassword
	}
	if len(password) > domain.MaxPasswordLength {
		return domain.ErrInvalidPassword
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return domain.ErrInvalidPassword
	}

	return nil
}

// ValidateName checks if the name is valid
func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > domain.MaxNameLength {
		return domain.ErrInvalidInput
	}
	return nil
}

// ValidateAttendanceStatus checks if the attendance status is valid
func ValidateAttendanceStatus(status string) error {
	if !domain.ValidAttendanceStatuses[status] {
		return domain.ErrInvalidAttendanceStatus
	}
	return nil
}

// ValidateUserRole checks if the user role is valid
func ValidateUserRole(role string) error {
	if !domain.ValidUserRoles[role] {
		return domain.ErrInvalidInput
	}
	return nil
}
