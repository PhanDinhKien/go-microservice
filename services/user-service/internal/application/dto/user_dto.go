package dto

import (
	"app-microservice/services/user-service/internal/domain/entities"
	"time"
)

// CreateUserRequest represents the request for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone,omitempty" binding:"omitempty"`
}

// UpdateUserRequest represents the request for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Email string `json:"email,omitempty" binding:"omitempty,email"`
	Phone string `json:"phone,omitempty" binding:"omitempty"`
}

// UserResponse represents the response format for user data
type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserListResponse represents the response format for user list
type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// SearchUsersRequest represents the request for searching users
type SearchUsersRequest struct {
	Query    string `form:"query" binding:"omitempty,min=1"`
	Status   string `form:"status" binding:"omitempty,oneof=active inactive suspended"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// UpdateUserStatusRequest represents the request for updating user status
type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive suspended"`
}

// Conversion methods

// ToEntity converts CreateUserRequest to User entity
func (req *CreateUserRequest) ToEntity() *entities.User {
	return &entities.User{
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		Status: string(entities.UserStatusActive),
	}
}

// ToUserResponse converts User entity to UserResponse
func ToUserResponse(user *entities.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserListResponse converts slice of User entities to UserListResponse
func ToUserListResponse(users []*entities.User, total int64, page, pageSize int) UserListResponse {
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = ToUserResponse(user)
	}

	totalPages := int(total / int64(pageSize))
	if total%int64(pageSize) > 0 {
		totalPages++
	}

	return UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// ApplyToEntity applies UpdateUserRequest changes to User entity
func (req *UpdateUserRequest) ApplyToEntity(user *entities.User) {
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	user.UpdatedAt = time.Now()
}
