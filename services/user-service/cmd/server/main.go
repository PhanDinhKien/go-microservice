package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app-microservice/services/user-service/internal/application/usecases"
	"app-microservice/services/user-service/internal/config"
	"app-microservice/services/user-service/internal/delivery/http/handlers"
	"app-microservice/services/user-service/internal/delivery/http/routes"
	"app-microservice/services/user-service/internal/domain/services"
	"app-microservice/services/user-service/internal/infrastructure/database"
	"app-microservice/services/user-service/internal/infrastructure/repositories"
	"app-microservice/services/user-service/pkg/logger"

	_ "app-microservice/services/user-service/docs"
)

// @title User Service API
// @version 1.0
// @description User Service API for microservice architecture with Clean Architecture pattern
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /api/v1
// @schemes http https

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	appLogger := logger.New(logger.Config{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
	})
	logger.SetGlobalLogger(appLogger)

	logger.Info("Starting User Service...")
	logger.Infof("Configuration loaded: Server Port=%s, DB Driver=%s", cfg.Server.Port, cfg.Database.Driver)

	// Initialize database
	dbConfig := &database.Config{
		Host:                  cfg.Database.Host,
		Port:                  cfg.Database.Port,
		User:                  cfg.Database.User,
		Password:              cfg.Database.Password,
		DBName:                cfg.Database.DBName,
		SSLMode:               cfg.Database.SSLMode,
		Driver:                cfg.Database.Driver,
		MaxOpenConnections:    cfg.Database.MaxOpenConnections,
		MaxIdleConnections:    cfg.Database.MaxIdleConnections,
		ConnectionMaxLifetime: cfg.Database.ConnectionMaxLifetime,
	}

	entClient, err := database.NewConnection(dbConfig)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := entClient.Close(); err != nil {
			logger.Errorf("Failed to close database connection: %v", err)
		} else {
			logger.Info("Database connection closed")
		}
	}()
	logger.Info("Database connection established")

	// Run migrations
	ctx := context.Background()
	if err := database.AutoMigrate(ctx, entClient); err != nil {
		logger.Fatalf("Failed to run database migrations: %v", err)
	}
	logger.Info("Database migrations completed")

	// Seed initial data
	if err := database.SeedData(ctx, entClient); err != nil {
		logger.Warnf("Failed to seed database: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(entClient)
	logger.Info("Repositories initialized")

	// Initialize domain services
	userDomainService := services.NewUserDomainService(userRepo)
	logger.Info("Domain services initialized")

	// Initialize use cases
	userUseCase := usecases.NewUserUseCase(userRepo, userDomainService)
	logger.Info("Use cases initialized")

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUseCase)
	logger.Info("Handlers initialized")

	// Initialize router
	router := routes.NewRouter(userHandler)
	ginEngine := router.SetupRoutes()
	logger.Info("Routes configured")

	// Setup HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      ginEngine,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("User Service starting on port %s", cfg.Server.Port)
		logger.Infof("Swagger UI available at: http://localhost:%s/swagger/index.html", cfg.Server.Port)
		logger.Infof("Health check available at: http://localhost:%s/health", cfg.Server.Port)
		logger.Infof("API endpoints available at: http://localhost:%s/api/v1/users", cfg.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down User Service...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("User Service stopped")
}
