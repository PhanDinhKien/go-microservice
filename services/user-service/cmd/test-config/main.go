package main

import (
	"fmt"

	"app-microservice/services/user-service/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	
	fmt.Println("=== Configuration Debug ===")
	fmt.Printf("DB Driver: %s\n", cfg.Database.Driver)
	fmt.Printf("DB Host: %s\n", cfg.Database.Host)
	fmt.Printf("DB Port: %s\n", cfg.Database.Port)
	fmt.Printf("DB User: %s\n", cfg.Database.User)
	fmt.Printf("DB Password: %s\n", cfg.Database.Password)
	fmt.Printf("DB Name: %s\n", cfg.Database.DBName)
	fmt.Printf("DB SSLMode: %s\n", cfg.Database.SSLMode)
	fmt.Printf("Max Open Connections: %d\n", cfg.Database.MaxOpenConnections)
	fmt.Printf("Max Idle Connections: %d\n", cfg.Database.MaxIdleConnections)
	fmt.Printf("Connection Max Lifetime: %s\n", cfg.Database.ConnectionMaxLifetime)
}
