package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"app-microservice/services/user-service/ent"
	"app-microservice/services/user-service/ent/user"
	"app-microservice/services/user-service/internal/domain/entities"
	"app-microservice/services/user-service/internal/domain/repositories"
	"app-microservice/services/user-service/internal/infrastructure/database"
)

// userRepository implements the UserRepository interface using Ent
type userRepository struct {
	client *ent.Client
}

// NewUserRepository creates a new user repository
func NewUserRepository(client *ent.Client) repositories.UserRepository {
	return &userRepository{
		client: client,
	}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, domainUser *entities.User) error {
	entUser, err := database.DomainUserToEntCreate(r.client, domainUser).Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return errors.New("email already exists")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Update domain user with generated ID and timestamps
	domainUser.ID = entUser.ID
	domainUser.CreatedAt = entUser.CreatedAt
	domainUser.UpdatedAt = entUser.UpdatedAt

	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id int) (*entities.User, error) {
	entUser, err := r.client.User.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return database.EntUserToDomainUser(entUser), nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	entUser, err := r.client.User.Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return database.EntUserToDomainUser(entUser), nil
}

// GetAll retrieves all users with pagination
func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	entUsers, err := r.client.User.Query().
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(user.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	domainUsers := make([]*entities.User, len(entUsers))
	for i, entUser := range entUsers {
		domainUsers[i] = database.EntUserToDomainUser(entUser)
	}

	return domainUsers, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, domainUser *entities.User) error {
	updateBuilder := r.client.User.UpdateOneID(domainUser.ID).
		SetName(domainUser.Name).
		SetEmail(domainUser.Email).
		SetStatus(user.Status(domainUser.Status))

	if domainUser.Phone != "" {
		updateBuilder = updateBuilder.SetPhone(domainUser.Phone)
	}

	entUser, err := updateBuilder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("user not found")
		}
		if ent.IsConstraintError(err) {
			return errors.New("email already exists")
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Update domain user with new timestamps
	domainUser.UpdatedAt = entUser.UpdatedAt

	return nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id int) error {
	err := r.client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// Count returns the total number of users
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.client.User.Query().Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return int64(count), nil
}

// GetByStatus retrieves users by status
func (r *userRepository) GetByStatus(ctx context.Context, status string, limit, offset int) ([]*entities.User, error) {
	entUsers, err := r.client.User.Query().
		Where(user.StatusEQ(user.Status(status))).
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(user.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by status: %w", err)
	}

	domainUsers := make([]*entities.User, len(entUsers))
	for i, entUser := range entUsers {
		domainUsers[i] = database.EntUserToDomainUser(entUser)
	}

	return domainUsers, nil
}

// Search searches users by name or email
func (r *userRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entities.User, error) {
	entUsers, err := r.client.User.Query().
		Where(
			user.Or(
				user.NameContains(strings.ToLower(query)),
				user.EmailContains(strings.ToLower(query)),
				user.NameContainsFold(query),
				user.EmailContainsFold(query),
			),
		).
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(user.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	domainUsers := make([]*entities.User, len(entUsers))
	for i, entUser := range entUsers {
		domainUsers[i] = database.EntUserToDomainUser(entUser)
	}

	return domainUsers, nil
}

// Exists checks if a user exists by ID
func (r *userRepository) Exists(ctx context.Context, id int) (bool, error) {
	exists, err := r.client.User.Query().
		Where(user.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return exists, nil
}

// EmailExists checks if an email already exists
func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := r.client.User.Query().
		Where(user.EmailEQ(email)).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check if email exists: %w", err)
	}
	return exists, nil
}
