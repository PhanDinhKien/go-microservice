package usecases

import (
	"app-microservice/services/user-service/internal/application/dto"
	"app-microservice/services/user-service/internal/domain/entities"
	"app-microservice/services/user-service/internal/domain/repositories"
	"app-microservice/services/user-service/internal/domain/services"
	"context"
)

// UserUseCase defines the interface for user use cases
type UserUseCase interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)

	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id int) (*dto.UserResponse, error)

	// GetUsers retrieves users with pagination and optional filters
	GetUsers(ctx context.Context, req *dto.SearchUsersRequest) (*dto.UserListResponse, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, id int, req *dto.UpdateUserRequest) (*dto.UserResponse, error)

	// DeleteUser deletes a user by ID
	DeleteUser(ctx context.Context, id int) error

	// UpdateUserStatus updates user status
	UpdateUserStatus(ctx context.Context, id int, req *dto.UpdateUserStatusRequest) (*dto.UserResponse, error)

	// SearchUsers searches users by query
	SearchUsers(ctx context.Context, req *dto.SearchUsersRequest) (*dto.UserListResponse, error)
}

type userUseCase struct {
	userRepo      repositories.UserRepository
	userDomainSvc services.UserDomainService
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(
	userRepo repositories.UserRepository,
	userDomainSvc services.UserDomainService,
) UserUseCase {
	return &userUseCase{
		userRepo:      userRepo,
		userDomainSvc: userDomainSvc,
	}
}

// CreateUser creates a new user
func (uc *userUseCase) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Convert DTO to entity
	user := req.ToEntity()

	// Validate business rules
	if err := uc.userDomainSvc.ValidateUserCreation(ctx, user); err != nil {
		return nil, err
	}

	// Create user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Convert to response DTO
	response := dto.ToUserResponse(user)
	return &response, nil
}

// GetUserByID retrieves a user by ID
func (uc *userUseCase) GetUserByID(ctx context.Context, id int) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

// GetUsers retrieves users with pagination and optional filters
func (uc *userUseCase) GetUsers(ctx context.Context, req *dto.SearchUsersRequest) (*dto.UserListResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize

	var users []*entities.User
	var err error

	// Filter by status if provided
	if req.Status != "" {
		users, err = uc.userRepo.GetByStatus(ctx, req.Status, req.PageSize, offset)
	} else {
		users, err = uc.userRepo.GetAll(ctx, req.PageSize, offset)
	}

	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	response := dto.ToUserListResponse(users, total, req.Page, req.PageSize)
	return &response, nil
}

// UpdateUser updates an existing user
func (uc *userUseCase) UpdateUser(ctx context.Context, id int, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Get existing user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	req.ApplyToEntity(user)

	// Validate business rules
	if err := uc.userDomainSvc.ValidateUserUpdate(ctx, user); err != nil {
		return nil, err
	}

	// Update user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

// DeleteUser deletes a user by ID
func (uc *userUseCase) DeleteUser(ctx context.Context, id int) error {
	// Validate deletion
	if err := uc.userDomainSvc.CanDeleteUser(ctx, id); err != nil {
		return err
	}

	// Delete user
	return uc.userRepo.Delete(ctx, id)
}

// UpdateUserStatus updates user status
func (uc *userUseCase) UpdateUserStatus(ctx context.Context, id int, req *dto.UpdateUserStatusRequest) (*dto.UserResponse, error) {
	// Get existing user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update status based on request
	switch req.Status {
	case string(entities.UserStatusActive):
		user.Activate()
	case string(entities.UserStatusInactive):
		user.Deactivate()
	case string(entities.UserStatusSuspended):
		user.Suspend()
	}

	// Update user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

// SearchUsers searches users by query
func (uc *userUseCase) SearchUsers(ctx context.Context, req *dto.SearchUsersRequest) (*dto.UserListResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize

	// Search users
	users, err := uc.userRepo.Search(ctx, req.Query, req.PageSize, offset)
	if err != nil {
		return nil, err
	}

	// Get total count for search
	total := int64(len(users)) // This is a simplified approach

	response := dto.ToUserListResponse(users, total, req.Page, req.PageSize)
	return &response, nil
}
