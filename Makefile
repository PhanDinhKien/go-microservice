.PHONY: help build run test clean docker-build docker-run docker-stop

# Variables
DOCKER_COMPOSE = docker-compose

help: ## Hiển thị trợ giúp
	@echo "📋 Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Cài đặt dependencies
	@echo "📦 Installing dependencies..."
	@go mod tidy
	@go mod download

build: ## Build tất cả services
	@echo "🔨 Building services..."
	@go build -o bin/user-service ./services/user-service
	@go build -o bin/product-service ./services/product-service
	@go build -o bin/api-gateway ./services/api-gateway
	@echo "✅ Build completed!"

run: ## Chạy tất cả services
	@echo "🚀 Starting services..."
	@./start-services.sh

test: ## Test API endpoints
	@echo "🧪 Testing API endpoints..."
	@./test-api.sh

clean: ## Dọn dẹp build files
	@echo "🧹 Cleaning up..."
	@rm -rf bin/
	@echo "✅ Clean completed!"

# Docker commands
docker-build: ## Build Docker images
	@echo "🐳 Building Docker images..."
	@$(DOCKER_COMPOSE) build

docker-run: ## Chạy services với Docker
	@echo "🐳 Starting services with Docker..."
	@$(DOCKER_COMPOSE) up -d
	@echo "✅ Services started!"
	@echo "📋 Check status: make docker-status"

docker-stop: ## Dừng Docker services
	@echo "🐳 Stopping Docker services..."
	@$(DOCKER_COMPOSE) down

docker-status: ## Kiểm tra trạng thái Docker services
	@echo "🐳 Docker services status:"
	@$(DOCKER_COMPOSE) ps

docker-logs: ## Xem logs của Docker services
	@echo "📄 Docker services logs:"
	@$(DOCKER_COMPOSE) logs -f

# Development commands
dev-user: ## Chạy User Service trong development mode
	@echo "👤 Starting User Service..."
	@cd services/user-service && go run main.go

dev-product: ## Chạy Product Service trong development mode
	@echo "📦 Starting Product Service..."
	@cd services/product-service && go run main.go

dev-gateway: ## Chạy API Gateway trong development mode
	@echo "🌐 Starting API Gateway..."
	@cd services/api-gateway && go run main.go
