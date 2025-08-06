package main

// @title User Service API
// @version 1.0
// @description User Service API for microservice architecture
// @host localhost:8081
// @BasePath /

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"app-microservice/shared/config"
	"app-microservice/shared/models"

	_ "app-microservice/services/user-service/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var users = []models.User{
	{ID: 1, Name: "Nguyen Van A", Email: "a@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Name: "Tran Thi B", Email: "b@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

func main() {
	cfg := config.LoadConfig()

	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// Routes
	api := r.Group("/users")
	{
		api.GET("", getUsers)
		api.GET("/:id", getUserByID)
		api.POST("", createUser)
		api.PUT("/:id", updateUser)
		api.DELETE("/:id", deleteUser)
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "user-service",
			"time":    time.Now(),
		})
	})

	port := "8081"
	if cfg.Port != "" && cfg.Port != "8080" {
		port = cfg.Port
	}

	log.Printf("User Service đang chạy trên port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Get list of all users
// @Tags users
// @Produce json
// @Success 200 {object} models.Response
// @Router /users [get]
func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Lấy danh sách users thành công",
		Data:    users,
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /users/{id} [get]
func getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "ID không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, models.Response{
				Status:  "success",
				Message: "Lấy thông tin user thành công",
				Data:    user,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, models.ErrorResponse{
		Status:  "error",
		Message: "Không tìm thấy user",
		Error:   "User not found",
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the given data
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Router /users [post]
func createUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Dữ liệu không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	// Generate new ID
	newUser.ID = len(users) + 1
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	users = append(users, newUser)

	c.JSON(http.StatusCreated, models.Response{
		Status:  "success",
		Message: "Tạo user thành công",
		Data:    newUser,
	})
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update user data by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "Updated user data"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /users/{id} [put]
func updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "ID không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Dữ liệu không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	for i, user := range users {
		if user.ID == id {
			updatedUser.ID = id
			updatedUser.CreatedAt = user.CreatedAt
			updatedUser.UpdatedAt = time.Now()
			users[i] = updatedUser

			c.JSON(http.StatusOK, models.Response{
				Status:  "success",
				Message: "Cập nhật user thành công",
				Data:    updatedUser,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, models.ErrorResponse{
		Status:  "error",
		Message: "Không tìm thấy user",
		Error:   "User not found",
	})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /users/{id} [delete]
func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "ID không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, models.Response{
				Status:  "success",
				Message: "Xóa user thành công",
				Data:    nil,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, models.ErrorResponse{
		Status:  "error",
		Message: "Không tìm thấy user",
		Error:   "User not found",
	})
}
