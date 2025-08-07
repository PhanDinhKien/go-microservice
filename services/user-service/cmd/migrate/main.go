package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"app-microservice/services/user-service/internal/config"
	"app-microservice/services/user-service/internal/infrastructure/database"
)

func main() {
	// Load .env.migrate for local development migrations
	if err := godotenv.Load(".env.migrate"); err != nil {
		// Fall back to .env if .env.migrate doesn't exist
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Could not load any .env file: %v", err)
		}
	}
	// Load configuration
	cfg := config.LoadConfig()

	// Create database connection
	client, err := database.NewConnection(&database.Config{
		Driver:                cfg.Database.Driver,
		Host:                  cfg.Database.Host,
		Port:                  cfg.Database.Port,
		User:                  cfg.Database.User,
		Password:              cfg.Database.Password,
		DBName:                cfg.Database.DBName,
		SSLMode:               cfg.Database.SSLMode,
		MaxOpenConnections:    cfg.Database.MaxOpenConnections,
		MaxIdleConnections:    cfg.Database.MaxIdleConnections,
		ConnectionMaxLifetime: cfg.Database.ConnectionMaxLifetime,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Get command argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/migrate/main.go [migrate|seed|status]")
	}

	command := os.Args[1]

	switch command {
	case "migrate":
		if err := database.AutoMigrate(ctx, client); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("✅ Migration completed successfully")

	case "seed":
		if err := database.SeedData(ctx, client); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
		log.Println("✅ Seeding completed successfully")

	case "status":
		// Show current database status
		count, err := client.User.Query().Count(ctx)
		if err != nil {
			log.Fatalf("Failed to count users: %v", err)
		}
		
		log.Println("📊 Database Status:")
		log.Printf("   Users: %d records", count)
		log.Println("   Migration: Users table exists with Ent schema")

	default:
		log.Fatalf("Unknown command: %s. Use: migrate, seed, or status", command)
	}
}
