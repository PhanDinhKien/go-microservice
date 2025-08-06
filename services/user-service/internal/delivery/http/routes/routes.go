package routes

import (
	"app-microservice/services/user-service/internal/delivery/http/handlers"
	"app-microservice/services/user-service/internal/delivery/http/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router holds the router configuration
type Router struct {
	userHandler *handlers.UserHandler
}

// NewRouter creates a new router instance
func NewRouter(userHandler *handlers.UserHandler) *Router {
	return &Router{
		userHandler: userHandler,
	}
}

// SetupRoutes configures all routes for the application
func (r *Router) SetupRoutes() *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.SecurityMiddleware())
	router.Use(gin.Recovery())
	router.Use(middleware.ValidationErrorMiddleware())

	// Health check endpoint
	router.GET("/health", r.healthCheck)

	// API version 1 routes
	v1 := router.Group("/api/v1")
	{
		r.setupUserRoutes(v1)
	}

	// User routes (for backward compatibility)
	userGroup := router.Group("/users")
	{
		r.setupUserRoutesCompat(userGroup)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Default route for undefined endpoints
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Endpoint not found",
			"error":   "The requested endpoint does not exist",
		})
	})

	return router
}

// setupUserRoutes sets up user-related routes for v1 API
func (r *Router) setupUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.POST("", r.userHandler.CreateUser)
		users.GET("", r.userHandler.GetUsers)
		users.GET("/:id", r.userHandler.GetUserByID)
		users.PUT("/:id", r.userHandler.UpdateUser)
		users.DELETE("/:id", r.userHandler.DeleteUser)
		users.PATCH("/:id/status", r.userHandler.UpdateUserStatus)
	}
}

// setupUserRoutesCompat sets up user routes for backward compatibility
func (r *Router) setupUserRoutesCompat(rg *gin.RouterGroup) {
	rg.POST("", r.userHandler.CreateUser)
	rg.GET("", r.userHandler.GetUsers)
	rg.GET("/:id", r.userHandler.GetUserByID)
	rg.PUT("/:id", r.userHandler.UpdateUser)
	rg.DELETE("/:id", r.userHandler.DeleteUser)
	rg.PATCH("/:id/status", r.userHandler.UpdateUserStatus)
}

// healthCheck provides a health check endpoint
func (r *Router) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "user-service",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC(),
		"uptime":    time.Since(startTime).String(),
	})
}

// startTime records when the service started
var startTime = time.Now()
