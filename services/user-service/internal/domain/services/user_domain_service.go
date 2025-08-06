package services

import (
	"app-microservice/services/user-service/internal/domain/entities"
	"app-microservice/services/user-service/internal/domain/repositories"
	"context"
	"errors"
)

// UserDomainService provides domain business logic for users
type UserDomainService interface {
	// ValidateUserCreation validates business rules for user creation
	ValidateUserCreation(ctx context.Context, user *entities.User) error

	// ValidateUserUpdate validates business rules for user update
	ValidateUserUpdate(ctx context.Context, user *entities.User) error

	// CanDeleteUser checks if a user can be deleted
	CanDeleteUser(ctx context.Context, userID int) error
}

type userDomainService struct {
	userRepo repositories.UserRepository
}

// NewUserDomainService creates a new user domain service
func NewUserDomainService(userRepo repositories.UserRepository) UserDomainService {
	return &userDomainService{
		userRepo: userRepo,
	}
}

// ValidateUserCreation validates business rules for user creation
func (s *userDomainService) ValidateUserCreation(ctx context.Context, user *entities.User) error {
	// Validate entity rules
	if err := user.Validate(); err != nil {
		return err
	}

	// Check if email already exists
	exists, err := s.userRepo.EmailExists(ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	// Set default status if not provided
	if user.Status == "" {
		user.Status = string(entities.UserStatusActive)
	}

	return nil
}

// ValidateUserUpdate validates business rules for user update
func (s *userDomainService) ValidateUserUpdate(ctx context.Context, user *entities.User) error {
	// Validate entity rules
	if err := user.Validate(); err != nil {
		return err
	}

	// Check if user exists
	exists, err := s.userRepo.Exists(ctx, user.ID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user not found")
	}

	// Check if email is being changed and if it already exists for another user
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil && existingUser.ID != user.ID {
		return errors.New("email already exists for another user")
	}

	return nil
}

// CanDeleteUser checks if a user can be deleted
func (s *userDomainService) CanDeleteUser(ctx context.Context, userID int) error {
	// Check if user exists
	exists, err := s.userRepo.Exists(ctx, userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user not found")
	}

	// Add any business rules for deletion
	// For example: check if user has active orders, etc.

	return nil
}
