package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"app-microservice/services/user-service/ent"
	"app-microservice/services/user-service/ent/user"
	"app-microservice/services/user-service/internal/domain/entities"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3"    // SQLite driver
)

// MigrationRecord represents a migration record
type MigrationRecord struct {
	ID          int       `json:"id"`
	Version     string    `json:"version"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AppliedAt   time.Time `json:"applied_at"`
	Success     bool      `json:"success"`
}

// Config holds database configuration
type Config struct {
	Host                  string
	Port                  string
	User                  string
	Password              string
	DBName                string
	SSLMode               string
	Driver                string // postgres, sqlite
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime string
}

// DatabaseConnection holds both Ent client and raw DB connection
type DatabaseConnection struct {
	Client *ent.Client
	DB     *sql.DB
}

// NewConnection creates a new Ent client connection
func NewConnection(config *Config) (*ent.Client, error) {
	var drv *entsql.Driver

	switch config.Driver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

		db, err := sql.Open("pgx", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open postgres connection: %w", err)
		}

		// Configure connection pool
		db.SetMaxOpenConns(config.MaxOpenConnections)
		db.SetMaxIdleConns(config.MaxIdleConnections)

		// Parse connection max lifetime
		if duration, err := time.ParseDuration(config.ConnectionMaxLifetime); err == nil {
			db.SetConnMaxLifetime(duration)
		}

		// Test connection
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping postgres database: %w", err)
		}

		drv = entsql.OpenDB(dialect.Postgres, db)

	case "sqlite":
		db, err := sql.Open("sqlite3", config.DBName)
		if err != nil {
			return nil, fmt.Errorf("failed to open sqlite connection: %w", err)
		}

		// Test connection
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping sqlite database: %w", err)
		}

		drv = entsql.OpenDB(dialect.SQLite, db)

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}

	client := ent.NewClient(ent.Driver(drv))

	// Enable debug mode in development
	if config.Driver == "sqlite" {
		client = client.Debug()
	}

	return client, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(ctx context.Context, client *ent.Client) error {
	log.Println("Starting database migration...")
	
	// Log the schema that will be created
	log.Println("Creating/updating schema for entities: User")
	
	// Run the actual migration
	migrationStart := time.Now()
	err := client.Schema.Create(ctx)
	
	if err != nil {
		log.Printf("Failed to create database schema: %v", err)
		return fmt.Errorf("failed to create database schema: %w", err)
	}

	log.Printf("✅ Database migration completed successfully in %v", time.Since(migrationStart))
	log.Println("   - Table: users (with indexes: email, status, created_at)")
	return nil
}

// SeedData seeds initial data into the database
func SeedData(ctx context.Context, client *ent.Client) error {
	// Check if users already exist
	count, err := client.User.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	if count > 0 {
		log.Printf("Database already has %d users, skipping seed", count)
		return nil
	}

	log.Println("Seeding initial data...")

	// Seed initial users
	users := []struct {
		name   string
		email  string
		phone  string
		status user.Status
	}{
		{
			name:   "Nguyen Van A",
			email:  "a@example.com",
			phone:  "+84123456789",
			status: user.StatusActive,
		},
		{
			name:   "Tran Thi B",
			email:  "b@example.com",
			phone:  "+84987654321",
			status: user.StatusActive,
		},
		{
			name:   "Le Van C",
			email:  "c@example.com",
			phone:  "+84555666777",
			status: user.StatusInactive,
		},
	}

	for i, userData := range users {
		_, err := client.User.Create().
			SetName(userData.name).
			SetEmail(userData.email).
			SetPhone(userData.phone).
			SetStatus(userData.status).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to seed user data: %w", err)
		}
		log.Printf("   ✓ Created user %d: %s (%s)", i+1, userData.name, userData.email)
	}

	log.Printf("✅ Database seeded successfully with %d users", len(users))
	return nil
}

// EntUserToDomainUser converts Ent User to Domain User
func EntUserToDomainUser(entUser *ent.User) *entities.User {
	return &entities.User{
		ID:        entUser.ID,
		Name:      entUser.Name,
		Email:     entUser.Email,
		Phone:     entUser.Phone,
		Status:    string(entUser.Status),
		CreatedAt: entUser.CreatedAt,
		UpdatedAt: entUser.UpdatedAt,
	}
}

// DomainUserToEntCreate converts Domain User to Ent User Create builder
func DomainUserToEntCreate(client *ent.Client, domainUser *entities.User) *ent.UserCreate {
	create := client.User.Create().
		SetName(domainUser.Name).
		SetEmail(domainUser.Email).
		SetStatus(user.Status(domainUser.Status))

	if domainUser.Phone != "" {
		create = create.SetPhone(domainUser.Phone)
	}

	return create
}
