package entities

import (
	"errors"
	"regexp"
	"time"
)

// User represents the user domain entity
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" validate:"required,min=2,max=100"`
	Email     string    `json:"email" validate:"required,email"`
	Phone     string    `json:"phone,omitempty" validate:"omitempty,phone"`
	Status    string    `json:"status" validate:"omitempty,oneof=active inactive suspended"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserStatus represents possible user statuses
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// Business rules and validation methods

// Validate performs domain-level validation
func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("user name is required")
	}

	if len(u.Name) < 2 {
		return errors.New("user name must be at least 2 characters")
	}

	if len(u.Name) > 100 {
		return errors.New("user name must not exceed 100 characters")
	}

	if u.Email == "" {
		return errors.New("user email is required")
	}

	if !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}

	if u.Phone != "" && !isValidPhone(u.Phone) {
		return errors.New("invalid phone format")
	}

	if u.Status != "" && !isValidStatus(string(u.Status)) {
		return errors.New("invalid user status")
	}

	return nil
}

// IsActive checks if user is active
func (u *User) IsActive() bool {
	return u.Status == string(UserStatusActive)
}

// Activate sets user status to active
func (u *User) Activate() {
	u.Status = string(UserStatusActive)
	u.UpdatedAt = time.Now()
}

// Deactivate sets user status to inactive
func (u *User) Deactivate() {
	u.Status = string(UserStatusInactive)
	u.UpdatedAt = time.Now()
}

// Suspend sets user status to suspended
func (u *User) Suspend() {
	u.Status = string(UserStatusSuspended)
	u.UpdatedAt = time.Now()
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(name, email, phone string) error {
	if name != "" {
		u.Name = name
	}
	if email != "" {
		u.Email = email
	}
	if phone != "" {
		u.Phone = phone
	}

	u.UpdatedAt = time.Now()

	return u.Validate()
}

// Helper functions for validation
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^[\+]?[\d\s\-\(\)]{10,15}$`)
	return phoneRegex.MatchString(phone)
}

func isValidStatus(status string) bool {
	validStatuses := []string{
		string(UserStatusActive),
		string(UserStatusInactive),
		string(UserStatusSuspended),
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}
