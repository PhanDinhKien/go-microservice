# Microservice Architecture with Go Gin

Dự án này bao gồm 2 microservice sử dụng Go Gin framework với tích hợp Swagger để testing và documentation API.

## Cấu trúc dự án

```
app-microservice/
├── services/
│   ├── user-service/          # Service quản lý người dùng
│   │   ├── main.go           # Main application file
│   │   ├── docs/             # Swagger documentation
│   │   │   └── docs.go       # Generated docs
│   │   ├── .env              # Environment variables
│   │   └── Dockerfile        # Docker configuration
│   ├── product-service/       # Service quản lý sản phẩm
│   │   ├── main.go           # Main application file
│   │   ├── docs/             # Swagger documentation
│   │   │   └── docs.go       # Generated docs
│   │   ├── .env              # Environment variables
│   │   └── Dockerfile        # Docker configuration
│   └── api-gateway/          # API Gateway
│       ├── main.go           # Main application file
│       ├── .env              # Environment variables
│       └── Dockerfile        # Docker configuration
├── shared/                   # Thư viện dùng chung
│   └── models/
│       └── models.go         # Common data models
├── docker-compose.yml        # Docker compose config
├── go.mod                    # Go module file
├── go.sum                    # Go dependencies
└── README.md
```

## Services

### 1. User Service (Port 8081)
- **Chức năng**: Quản lý thông tin người dùng (CRUD operations)
- **Port**: 8081
- **Endpoints**: `/users/*`
- **Swagger UI**: `http://localhost:8081/swagger/index.html`
- **Health Check**: `http://localhost:8081/health`

### 2. Product Service (Port 8082)
- **Chức năng**: Quản lý thông tin sản phẩm (CRUD operations)
- **Port**: 8082
- **Endpoints**: `/products/*`
- **Swagger UI**: `http://localhost:8082/swagger/index.html`
- **Health Check**: `http://localhost:8082/health`

### 3. API Gateway (Port 8080)
- **Chức năng**: Route requests đến các service tương ứng
- **Port**: 8080
- **Features**: Load balancing, CORS, Health monitoring
- **Health Check**: `http://localhost:8080/health`

## Cách chạy dự án

### Yêu cầu hệ thống
- Go 1.22+
- Docker & Docker Compose (optional)

### 1. Cài đặt dependencies
```bash
# Từ thư mục root của dự án
go mod tidy
```

### 2. Chạy tất cả services bằng Docker
```bash
# Build và chạy tất cả services
docker-compose up -d

# Xem logs
docker-compose logs -f

# Dừng services
docker-compose down
```

### 3. Chạy từng service riêng lẻ

#### User Service
```bash
# Terminal 1: Chạy User Service
cd services/user-service
go run main.go

# Service sẽ chạy trên: http://localhost:8081
# Swagger UI: http://localhost:8081/swagger/index.html
```

#### Product Service
```bash
# Terminal 2: Chạy Product Service  
cd services/product-service
go run main.go

# Service sẽ chạy trên: http://localhost:8082
# Swagger UI: http://localhost:8082/swagger/index.html
```

#### API Gateway
```bash
# Terminal 3: Chạy API Gateway
cd services/api-gateway
go run main.go

# Gateway sẽ chạy trên: http://localhost:8080
```

### 4. Chạy từ root directory
```bash
# Từ thư mục root, có thể chạy:
go run services/user-service/main.go &
go run services/product-service/main.go &
go run services/api-gateway/main.go &
```

## Truy cập Swagger UI

### User Service Swagger
- **URL**: `http://localhost:8081/swagger/index.html`
- **Mô tả**: API documentation cho User Service
- **Features**: Test CRUD operations cho Users

### Product Service Swagger
- **URL**: `http://localhost:8082/swagger/index.html`
- **Mô tả**: API documentation cho Product Service
- **Features**: Test CRUD operations cho Products

### Sử dụng Swagger UI
1. Truy cập URL Swagger tương ứng
2. Explore các API endpoints
3. Click "Try it out" để test API
4. Nhập parameters và click "Execute"
5. Xem response data và status code

## API Endpoints

### User Service (Port 8081)
- `GET /users` - Lấy danh sách users
- `GET /users/:id` - Lấy thông tin user theo ID
- `POST /users` - Tạo user mới
- `PUT /users/:id` - Cập nhật user
- `DELETE /users/:id` - Xóa user
- `GET /health` - Health check
- `GET /swagger/*` - Swagger UI

### Product Service (Port 8082)
- `GET /products` - Lấy danh sách products
- `GET /products/:id` - Lấy thông tin product theo ID
- `POST /products` - Tạo product mới
- `PUT /products/:id` - Cập nhật product
- `DELETE /products/:id` - Xóa product
- `GET /health` - Health check
- `GET /swagger/*` - Swagger UI

### API Gateway (Port 8080)
- `GET /user/*` - Proxy đến User Service
- `GET /product/*` - Proxy đến Product Service
- `GET /health` - Health check tổng thể

## Testing APIs

### 1. Qua Swagger UI (Recommended)
- User Service: `http://localhost:8081/swagger/index.html`
- Product Service: `http://localhost:8082/swagger/index.html`

### 2. Qua cURL
```bash
# Test User Service
curl -X GET http://localhost:8081/users
curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com"}'

# Test Product Service
curl -X GET http://localhost:8082/products
curl -X POST http://localhost:8082/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Product","price":100,"description":"Test Description"}'

# Test qua API Gateway
curl -X GET http://localhost:8080/user/users
curl -X GET http://localhost:8080/product/products
```

### 3. Qua Postman
Import các endpoint từ Swagger JSON:
- User Service: `http://localhost:8081/swagger/doc.json`
- Product Service: `http://localhost:8082/swagger/doc.json`

## Troubleshooting

### Port đã được sử dụng
```bash
# Kiểm tra process đang dùng port
lsof -i :8081
lsof -i :8082
lsof -i :8080

# Kill process nếu cần
kill -9 <PID>
```

### Swagger không hiển thị
1. Kiểm tra service đã chạy chưa
2. Truy cập đúng URL Swagger
3. Kiểm tra docs/docs.go file có tồn tại không

### Dependencies issues
```bash
# Clean và reinstall
go clean -modcache
go mod tidy
go mod download
```
