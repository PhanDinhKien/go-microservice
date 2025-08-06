package fixtures

import (
	"app-microservice/services/user-service/internal/application/dto"
	"app-microservice/services/user-service/internal/domain/entities"
	"time"
)

// UserFixtures provides test data for users
type UserFixtures struct{}

// NewUserFixtures creates a new user fixtures instance
func NewUserFixtures() *UserFixtures {
	return &UserFixtures{}
}

// ValidUser returns a valid user entity
func (f *UserFixtures) ValidUser() *entities.User {
	return &entities.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Phone:     "+1234567890",
		Status:    string(entities.UserStatusActive),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// InvalidUser returns an invalid user entity
func (f *UserFixtures) InvalidUser() *entities.User {
	return &entities.User{
		ID:    1,
		Name:  "",              // Invalid: empty name
		Email: "invalid-email", // Invalid: bad email format
		Phone: "123",           // Invalid: bad phone format
	}
}

// MultipleUsers returns a slice of user entities
func (f *UserFixtures) MultipleUsers() []*entities.User {
	return []*entities.User{
		{
			ID:        1,
			Name:      "John Doe",
			Email:     "john@example.com",
			Phone:     "+1234567890",
			Status:    string(entities.UserStatusActive),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "Jane Smith",
			Email:     "jane@example.com",
			Phone:     "+9876543210",
			Status:    string(entities.UserStatusActive),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        3,
			Name:      "Bob Johnson",
			Email:     "bob@example.com",
			Phone:     "+5555666777",
			Status:    string(entities.UserStatusInactive),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// ValidCreateUserRequest returns a valid create user request
func (f *UserFixtures) ValidCreateUserRequest() *dto.CreateUserRequest {
	return &dto.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
		Phone: "+1234567890",
	}
}

// InvalidCreateUserRequest returns an invalid create user request
func (f *UserFixtures) InvalidCreateUserRequest() *dto.CreateUserRequest {
	return &dto.CreateUserRequest{
		Name:  "",              // Invalid: empty name
		Email: "invalid-email", // Invalid: bad email format
		Phone: "123",           // Invalid: bad phone format
	}
}

// ValidUpdateUserRequest returns a valid update user request
func (f *UserFixtures) ValidUpdateUserRequest() *dto.UpdateUserRequest {
	return &dto.UpdateUserRequest{
		Name:  "Jane Doe",
		Email: "jane@example.com",
		Phone: "+9876543210",
	}
}

// ValidSearchUsersRequest returns a valid search users request
func (f *UserFixtures) ValidSearchUsersRequest() *dto.SearchUsersRequest {
	return &dto.SearchUsersRequest{
		Query:    "John",
		Status:   string(entities.UserStatusActive),
		Page:     1,
		PageSize: 10,
	}
}

// ValidUpdateUserStatusRequest returns a valid update user status request
func (f *UserFixtures) ValidUpdateUserStatusRequest() *dto.UpdateUserStatusRequest {
	return &dto.UpdateUserStatusRequest{
		Status: string(entities.UserStatusSuspended),
	}
}

// UserResponse returns a user response
func (f *UserFixtures) UserResponse() *dto.UserResponse {
	user := f.ValidUser()
	response := dto.ToUserResponse(user)
	return &response
}

// UserListResponse returns a user list response
func (f *UserFixtures) UserListResponse() *dto.UserListResponse {
	users := f.MultipleUsers()
	response := dto.ToUserListResponse(users, int64(len(users)), 1, 10)
	return &response
}
