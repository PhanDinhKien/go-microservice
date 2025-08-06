#!/bin/bash

echo "🚀 Khởi động Microservices..."

# Function to run service
run_service() {
    local service_name=$1
    local service_path=$2
    
    echo "📦 Khởi động $service_name..."
    cd "$service_path"
    go run main.go &
    cd - > /dev/null
}

# Start services
run_service "User Service" "./services/user-service"
sleep 2
run_service "Product Service" "./services/product-service"
sleep 2
run_service "API Gateway" "./services/api-gateway"

echo "✅ Tất cả services đã được khởi động!"
echo ""
echo "📋 Thông tin services:"
echo "  🌐 API Gateway:    http://localhost:8080"
echo "  👤 User Service:   http://localhost:8081"
echo "  📦 Product Service: http://localhost:8082"
echo ""
echo "🔍 Health check:"
echo "  curl http://localhost:8080/services/health"
echo ""
echo "⚡ Để dừng tất cả services, nhấn Ctrl+C"

# Wait for all background processes
wait
