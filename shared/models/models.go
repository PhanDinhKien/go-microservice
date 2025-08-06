package models

import "time"

// User represents a user in the system
// @Description User information
type User struct {
	ID        int       `json:"id" example:"1"`
	Name      string    `json:"name" example:"Nguyen Van A"`
	Email     string    `json:"email" example:"a@example.com"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Product represents a product in the system
// @Description Product information
type Product struct {
	ID          int       `json:"id" example:"1"`
	Name        string    `json:"name" example:"Laptop Dell"`
	Description string    `json:"description" example:"Laptop Dell Inspiron 15"`
	Price       float64   `json:"price" example:"15000000"`
	Stock       int       `json:"stock" example:"10"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Response represents a successful API response
// @Description Successful response
type Response struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error API response
// @Description Error response
type ErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Something went wrong"`
	Error   string `json:"error" example:"Detailed error message"`
}
