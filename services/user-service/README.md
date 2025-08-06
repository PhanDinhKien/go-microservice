# User Service

User Service là một microservice được xây dựng theo kiến trúc Clean Architecture, cung cấp các API để quản lý người dùng trong hệ thống microservice.

## Kiến trúc

Dự án được tổ chức theo **Clean Architecture** với **Ent ORM**, chia thành các layer rõ ràng và tuân theo nguyên tắc Dependency Inversion:

```
user-service/
├── .env                        # File biến môi trường runtime (KHÔNG commit vào git)
├── .env.example               # Template mẫu cho file biến môi trường
├── Dockerfile                 # Cấu hình container Docker cho deployment
├── Makefile                   # Tập lệnh tự động hóa build và development
├── README.md                  # Tài liệu chính của dự án
├── bin/                       # 📁 Thư mục chứa file binary đã compile
│   └── user-service          # File thực thi chính sau khi build
├── cmd/                       # 🎯 Điểm khởi đầu ứng dụng (Entry Points)
│   └── server/               # Server chính của ứng dụng
│       └── main.go          # ⭐ Khởi tạo server, dependency injection, wiring
├── ent/                      # 🤖 Code được generate bởi Ent ORM (KHÔNG sửa tay)
│   ├── schema/              # Định nghĩa schema cơ sở dữ liệu
│   │   └── user.go         # Schema entity User với fields, edges, indexes
│   ├── client.go           # Client kết nối cơ sở dữ liệu
│   ├── user.go             # Model User được generate
│   ├── user_create.go      # Các thao tác tạo User
│   ├── user_query.go       # Các thao tác truy vấn User
│   ├── user_update.go      # Các thao tác cập nhật User
│   └── ...                 # Các file khác được generate bởi Ent
├── internal/                 # 🔒 Code riêng tư của ứng dụng (không thể import từ ngoài)
│   ├── config/              # 🔧 Lớp Cấu hình (Configuration Layer)
│   │   └── config.go       # Đọc biến môi trường và cấu hình ứng dụng
│   ├── domain/              # 🏛️ Lớp Miền (Domain Layer - Logic nghiệp vụ cốt lõi)
│   │   ├── entities/        # Các thực thể nghiệp vụ cốt lõi
│   │   │   └── user.go     # Entity User với quy tắc nghiệp vụ và validation
│   │   ├── repositories/    # Giao diện Repository (contracts/ports)
│   │   │   └── user_repository.go  # Interface định nghĩa Repository cho User
│   │   └── services/        # Dịch vụ miền (Domain Services)
│   │       └── user_domain_service.go  # Logic nghiệp vụ phức tạp của User
│   ├── application/         # 📋 Lớp Ứng dụng (Application Layer - Use Cases)
│   │   ├── dto/            # Đối tượng truyền dữ liệu (Data Transfer Objects)
│   │   │   └── user_dto.go # Cấu trúc Request/Response cho API
│   │   ├── services/       # Dịch vụ ứng dụng (orchestration workflow)
│   │   └── usecases/       # Use Cases (quy trình nghiệp vụ)
│   │       └── user_usecase.go  # Triển khai các use case của User
│   ├── infrastructure/     # 🔌 Lớp Hạ tầng (Infrastructure Layer - Dependencies ngoài)
│   │   ├── database/       # Hạ tầng cơ sở dữ liệu
│   │   │   └── database.go # Cấu hình Ent client, migration, converter
│   │   ├── http/           # Hạ tầng HTTP (hiện tại không dùng)
│   │   └── repositories/   # Triển khai Repository (adapters)
│   │       └── user_repository.go # Repository User với Ent implementation
│   └── delivery/           # 🚀 Lớp Giao tiếp (Delivery Layer - Interface Adapters)
│       └── http/           # Cơ chế giao tiếp HTTP
│           ├── handlers/   # Xử lý HTTP request
│           │   └── user_handler.go  # Handler API cho User
│           ├── middleware/ # Middleware HTTP
│           │   └── middleware.go    # CORS, logging, security headers
│           └── routes/     # Định nghĩa route
│               └── routes.go        # Cấu hình API routes và middleware
├── pkg/                     # 📦 Gói công khai (có thể tái sử dụng)
│   ├── logger/             # Tiện ích logging
│   │   └── logger.go       # Triển khai structured logging với các level
│   └── utils/              # Tiện ích chung
│       └── utils.go        # Hàm helper, validator, converter
├── docs/                    # 📚 Tài liệu API (auto-generated)
│   ├── docs.go             # Cấu hình Swagger và metadata
│   ├── swagger.json        # Đặc tả OpenAPI định dạng JSON
│   └── swagger.yaml        # Đặc tả OpenAPI định dạng YAML
├── migrations/              # 🗄️ Migration cơ sở dữ liệu
│   └── 001_create_users_table.sql  # Script migration SQL
└── tests/                   # 🧪 Bộ test (Test Suite)
    ├── fixtures/           # Dữ liệu test và mock
    │   └── user_fixtures.go # Generator dữ liệu test cho User
    ├── integration/        # Test tích hợp (test toàn bộ workflow)
    └── unit/               # Unit test (test từng component riêng lẻ)
        └── user_entity_test.go  # Test logic nghiệp vụ của User entity
```

### 📋 Giải thích vai trò từng Layer (Clean Architecture)

#### 🏛️ **Lớp Miền - Domain Layer** (`internal/domain/`)
- **Entities (Thực thể)**: Các đối tượng nghiệp vụ cốt lõi chứa quy tắc kinh doanh và validation logic
- **Repositories (Kho dữ liệu)**: Giao diện (interface) định nghĩa cách truy cập dữ liệu (không có implementation cụ thể)
- **Services (Dịch vụ)**: Các dịch vụ miền chứa logic nghiệp vụ phức tạp không thuộc về một entity cụ thể

#### 📋 **Lớp Ứng dụng - Application Layer** (`internal/application/`)
- **DTOs (Data Transfer Objects)**: Cấu trúc dữ liệu để truyền thông tin giữa các layer
- **Use Cases (Trường hợp sử dụng)**: Điều phối các workflow nghiệp vụ, kết hợp các domain object
- **Services (Dịch vụ)**: Logic cụ thể của ứng dụng (không phải business logic)

#### 🔌 **Lớp Hạ tầng - Infrastructure Layer** (`internal/infrastructure/`)
- **Database (Cơ sở dữ liệu)**: Thiết lập Ent client, migration, quản lý schema
- **Repositories (Kho dữ liệu)**: Triển khai cụ thể các repository interface từ domain layer
- **External Services (Dịch vụ bên ngoài)**: Tích hợp với bên thứ ba (email, cache, v.v.)

#### 🚀 **Lớp Giao tiếp - Delivery Layer** (`internal/delivery/`)
- **HTTP Handlers (Xử lý HTTP)**: Chuyển đổi HTTP request thành các lời gọi use case
- **Middleware (Phần mềm trung gian)**: Các mối quan tâm chung (auth, logging, CORS)
- **Routes (Định tuyến)**: Định nghĩa API endpoint và thiết lập middleware

#### 📦 **Lớp Gói công khai - Package Layer** (`pkg/`)
- **Logger (Ghi log)**: Structured logging với các level khác nhau
- **Utils (Tiện ích)**: Các utility chung, validator, helper function

### 🔄 Data Flow
```
HTTP Request → Middleware → Handler → Use Case → Domain Service → Repository → Database
                   ↓           ↓         ↓           ↓            ↓
HTTP Response ← Handler ← DTO ← Use Case ← Domain ← Repository ← Ent Client
```

### 🛠️ **Technology Stack**
- **Framework**: Gin (HTTP router)
- **ORM**: Ent (Type-safe ORM with code generation)
- **Database**: SQLite (default) / PostgreSQL
- **Documentation**: Swagger/OpenAPI
- **Logging**: Custom structured logger
- **Testing**: Go testing + fixtures
- **Build**: Makefile + Go modules

## 📁 Detailed Directory Structure

### 🚀 **Root Level Files**
```
├── .env                    # Runtime environment variables (không commit vào git)
├── .env.example           # Template cho environment variables
├── Dockerfile             # Container definition cho deployment
├── Makefile              # Build automation và development tasks
└── README.md             # Documentation chính của project
```

### 📦 **Build & Binary**
```
├── bin/                   # Output directory cho compiled binaries
│   └── user-service      # Executable file sau khi build
```

### 🎯 **Application Entry Points**
```
├── cmd/                   # Application entry points (main packages)
│   └── server/           # Main server application
│       └── main.go      # ⭐ Dependency injection, server startup logic
```
- **Ý nghĩa**: Điểm khởi đầu của application, setup dependencies, run server

### 🗄️ **Ent ORM (Generated Code)**
```
├── ent/                   # Ent ORM generated code (KHÔNG edit manually)
│   ├── schema/           # Database schema definitions
│   │   └── user.go      # User entity schema với fields, edges, indexes
│   ├── client.go        # Database client với connection methods
│   ├── user.go          # Generated User entity model
│   ├── user_create.go   # Generated User creation operations
│   ├── user_query.go    # Generated User query operations
│   ├── user_update.go   # Generated User update operations
│   └── ...              # Các files khác được generate bởi Ent
```
- **Ý nghĩa**: Type-safe ORM code, auto-generated từ schema definitions

### 🏗️ **Private Application Code (Clean Architecture Layers)**
```
├── internal/             # Private application code (không thể import từ bên ngoài)
│   ├── config/          # 🔧 Configuration Layer
│   │   └── config.go   # Environment variables loading & app configuration
│   ├── domain/          # 🏛️ Domain Layer (Core Business Logic)
│   │   ├── entities/    # Core business entities
│   │   │   └── user.go # User domain entity với business rules & validation
│   │   ├── repositories/ # Repository interfaces (contracts/ports)
│   │   │   └── user_repository.go # User repository interface definition
│   │   └── services/    # Domain services (complex business logic)
│   │       └── user_domain_service.go # User business rules & validation
│   ├── application/     # 📋 Application Layer (Use Cases & Orchestration)
│   │   ├── dto/        # Data Transfer Objects
│   │   │   └── user_dto.go # Request/Response structures cho API
│   │   ├── services/   # Application services (workflow orchestration)
│   │   └── usecases/   # Use cases (business workflows implementation)
│   │       └── user_usecase.go # User use cases coordination
│   ├── infrastructure/ # 🔌 Infrastructure Layer (External Dependencies)
│   │   ├── database/   # Database infrastructure
│   │   │   └── database.go # Ent client setup, migrations, converters
│   │   ├── http/       # HTTP infrastructure (currently unused)
│   │   └── repositories/ # Repository implementations (adapters)
│   │       └── user_repository.go # User repository với Ent implementation
│   └── delivery/       # 🚀 Delivery Layer (Interface Adapters)
│       └── http/       # HTTP delivery mechanism
│           ├── handlers/ # HTTP request handlers
│           │   └── user_handler.go # User API handlers & HTTP logic
│           ├── middleware/ # HTTP middleware
│           │   └── middleware.go # CORS, logging, security headers
│           └── routes/ # Route definitions
│               └── routes.go # API route setup & middleware chaining
```

#### 🏛️ **Domain Layer Explained**
- **Entities**: Core business objects, chứa business rules và validation logic
- **Repositories**: Interface definitions, không có implementation cụ thể
- **Services**: Complex business logic không thuộc về 1 entity cụ thể

#### 📋 **Application Layer Explained**
- **DTOs**: Structures cho data transfer between layers
- **Use Cases**: Coordinate domain objects để thực hiện business workflows
- **Services**: Application-specific logic, orchestration

#### 🔌 **Infrastructure Layer Explained**
- **Database**: Concrete implementations cho data persistence
- **Repositories**: Implement domain repository interfaces
- **External Services**: Third-party integrations

#### 🚀 **Delivery Layer Explained**
- **Handlers**: Convert HTTP requests thành use case calls
- **Middleware**: Cross-cutting concerns như auth, logging
- **Routes**: API endpoint definitions

### 📦 **Reusable Packages**
```
├── pkg/                  # Public packages (có thể reuse ở projects khác)
│   ├── logger/          # Logging utilities
│   │   └── logger.go   # Structured logging implementation với levels
│   └── utils/          # Common utilities
│       └── utils.go    # Helper functions, validators, converters
```
- **Ý nghĩa**: Reusable components, có thể share across services

### 📚 **Documentation**
```
├── docs/                # API documentation (auto-generated)
│   ├── docs.go         # Swagger configuration & metadata
│   ├── swagger.json    # OpenAPI specification (JSON format)
│   └── swagger.yaml    # OpenAPI specification (YAML format)
```
- **Ý nghĩa**: Auto-generated API documentation từ code annotations

### 🗄️ **Database Migrations**
```
├── migrations/          # Database migration files
│   └── 001_create_users_table.sql # SQL migration scripts
```
- **Ý nghĩa**: Database schema version control, migration scripts

### 🧪 **Testing Suite**
```
└── tests/               # Test files organization
    ├── fixtures/        # Test data & mocks
    │   └── user_fixtures.go # User test data generators
    ├── integration/     # Integration tests (test full workflows)
    └── unit/           # Unit tests (test individual components)
        └── user_entity_test.go # User entity business logic tests
```
- **Ý nghĩa**: Comprehensive testing strategy với test data management

## Tính năng

- ✅ CRUD operations cho User
- ✅ Tìm kiếm và lọc người dùng
- ✅ Quản lý trạng thái người dùng (active/inactive/suspended)
- ✅ Validation dữ liệu đầu vào
- ✅ Swagger documentation
- ✅ Structured logging
- ✅ Health check endpoint
- ✅ CORS middleware
- ✅ Security headers
- ✅ Graceful shutdown
- ✅ Database migrations
- ✅ Unit tests

## API Endpoints

### Base URL: `/api/v1`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users` | Lấy danh sách users với pagination |
| GET | `/users/{id}` | Lấy thông tin user theo ID |
| POST | `/users` | Tạo user mới |
| PUT | `/users/{id}` | Cập nhật thông tin user |
| DELETE | `/users/{id}` | Xóa user |
| PATCH | `/users/{id}/status` | Cập nhật trạng thái user |

### Backward Compatibility

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users` | Compatible với API cũ |
| GET | `/users/{id}` | Compatible với API cũ |
| POST | `/users` | Compatible với API cũ |
| PUT | `/users/{id}` | Compatible với API cũ |
| DELETE | `/users/{id}` | Compatible với API cũ |

### Other Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/swagger/*` | Swagger documentation |

## Cấu hình

### Environment Variables

Copy file `.env.example` thành `.env` và cấu hình:

```bash
cp .env.example .env
```

Các biến môi trường quan trọng:

- `PORT`: Port của service (default: 8081)
- `DB_DRIVER`: Database driver (sqlite hoặc postgres)
- `DB_NAME`: Tên database
- `LOG_LEVEL`: Mức độ log (debug, info, warn, error, fatal)

### Database

Service hỗ trợ 2 loại database:

#### SQLite (Default)
```env
DB_DRIVER=sqlite
DB_NAME=user_service.db
```

#### PostgreSQL
```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=user_service
DB_SSLMODE=disable
```

## Cách chạy

### 1. Từ source code

```bash
# Cài đặt dependencies
go mod tidy

# Chạy từ cmd/server
cd cmd/server
go run main.go

# Hoặc chạy từ root directory
go run cmd/server/main.go
```

### 2. Build binary

```bash
# Build binary
go build -o bin/user-service cmd/server/main.go

# Chạy binary
./bin/user-service
```

### 3. Docker

```bash
# Build Docker image
docker build -t user-service .

# Chạy container
docker run -p 8081:8081 user-service
```

## Testing

### Unit Tests

```bash
# Chạy tất cả unit tests
go test ./tests/unit/...

# Chạy tests với coverage
go test -cover ./tests/unit/...

# Chạy tests với verbose output
go test -v ./tests/unit/...
```

### Integration Tests

```bash
# Chạy integration tests
go test ./tests/integration/...
```

### Test Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## API Documentation

Sau khi service đang chạy:

- **Swagger UI**: http://localhost:8081/swagger/index.html
- **OpenAPI JSON**: http://localhost:8081/swagger/doc.json

## Examples

### Create User

```bash
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890"
  }'
```

### Get Users

```bash
# Get all users
curl http://localhost:8081/api/v1/users

# Get users with pagination
curl "http://localhost:8081/api/v1/users?page=1&page_size=5"

# Search users
curl "http://localhost:8081/api/v1/users?query=John"

# Filter by status
curl "http://localhost:8081/api/v1/users?status=active"
```

### Update User

```bash
curl -X PUT http://localhost:8081/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com"
  }'
```

### Update User Status

```bash
curl -X PATCH http://localhost:8081/api/v1/users/1/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "suspended"
  }'
```

## Monitoring

### Health Check

```bash
curl http://localhost:8081/health
```

Response:
```json
{
  "status": "ok",
  "service": "user-service",
  "version": "1.0.0",
  "timestamp": "2025-01-01T12:00:00Z",
  "uptime": "5m30s"
}
```

### Logs

Service sử dụng structured logging với các mức độ:

- `DEBUG`: Debug information
- `INFO`: General information
- `WARN`: Warning messages
- `ERROR`: Error messages
- `FATAL`: Fatal errors (service will exit)

## Development

### Adding New Features

1. **Domain Layer**: Thêm entities, repository interfaces, domain services
2. **Application Layer**: Thêm DTOs, use cases
3. **Infrastructure Layer**: Implement repositories, database operations
4. **Delivery Layer**: Thêm HTTP handlers, routes
5. **Tests**: Viết unit tests và integration tests

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run linter (if golangci-lint is installed)
golangci-lint run
```

## Troubleshooting

### Common Issues

1. **Port already in use**: Thay đổi PORT trong .env file
2. **Database connection failed**: Kiểm tra cấu hình database
3. **Permission denied**: Đảm bảo user có quyền ghi file database (SQLite)

### Debug Mode

Set log level thành debug:

```env
LOG_LEVEL=debug
```

### Database Issues

```bash
# Check SQLite database
sqlite3 user_service.db ".schema"

# Check PostgreSQL connection
psql -h localhost -U postgres -d user_service -c "\dt"
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes following Clean Architecture principles
4. Add tests for new features
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License.
