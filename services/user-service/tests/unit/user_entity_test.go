package entities

import (
	"app-microservice/services/user-service/internal/domain/entities"
	"testing"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    entities.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user",
			user: entities.User{
				Name:   "John Doe",
				Email:  "john@example.com",
				Phone:  "+1234567890",
				Status: string(entities.UserStatusActive),
			},
			wantErr: false,
		},
		{
			name: "empty name",
			user: entities.User{
				Name:  "",
				Email: "john@example.com",
			},
			wantErr: true,
			errMsg:  "user name is required",
		},
		{
			name: "short name",
			user: entities.User{
				Name:  "A",
				Email: "john@example.com",
			},
			wantErr: true,
			errMsg:  "user name must be at least 2 characters",
		},
		{
			name: "empty email",
			user: entities.User{
				Name:  "John Doe",
				Email: "",
			},
			wantErr: true,
			errMsg:  "user email is required",
		},
		{
			name: "invalid email",
			user: entities.User{
				Name:  "John Doe",
				Email: "invalid-email",
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "invalid phone",
			user: entities.User{
				Name:  "John Doe",
				Email: "john@example.com",
				Phone: "123",
			},
			wantErr: true,
			errMsg:  "invalid phone format",
		},
		{
			name: "invalid status",
			user: entities.User{
				Name:   "John Doe",
				Email:  "john@example.com",
				Status: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid user status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("User.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUser_IsActive(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{
			name:   "active user",
			status: string(entities.UserStatusActive),
			want:   true,
		},
		{
			name:   "inactive user",
			status: string(entities.UserStatusInactive),
			want:   false,
		},
		{
			name:   "suspended user",
			status: string(entities.UserStatusSuspended),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &entities.User{Status: tt.status}
			if got := user.IsActive(); got != tt.want {
				t.Errorf("User.IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_UpdateProfile(t *testing.T) {
	user := &entities.User{
		Name:  "John Doe",
		Email: "john@example.com",
		Phone: "+1234567890",
	}

	err := user.UpdateProfile("Jane Doe", "jane@example.com", "+9876543210")
	if err != nil {
		t.Errorf("User.UpdateProfile() error = %v", err)
	}

	if user.Name != "Jane Doe" {
		t.Errorf("User.UpdateProfile() name = %v, want %v", user.Name, "Jane Doe")
	}
	if user.Email != "jane@example.com" {
		t.Errorf("User.UpdateProfile() email = %v, want %v", user.Email, "jane@example.com")
	}
	if user.Phone != "+9876543210" {
		t.Errorf("User.UpdateProfile() phone = %v, want %v", user.Phone, "+9876543210")
	}
}
