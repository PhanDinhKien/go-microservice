package repositories

import (
	"app-microservice/services/user-service/internal/domain/entities"
	"context"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entities.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int) (*entities.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)

	// GetAll retrieves all users with pagination
	GetAll(ctx context.Context, limit, offset int) ([]*entities.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entities.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id int) error

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)

	// GetByStatus retrieves users by status
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]*entities.User, error)

	// Search searches users by name or email
	Search(ctx context.Context, query string, limit, offset int) ([]*entities.User, error)

	// Exists checks if a user exists by ID
	Exists(ctx context.Context, id int) (bool, error)

	// EmailExists checks if an email already exists
	EmailExists(ctx context.Context, email string) (bool, error)
}
