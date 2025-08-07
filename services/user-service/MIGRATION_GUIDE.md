# Database Migration Guide

## Overview
This guide explains how to use the database migration system for the User Service.

## Prerequisites
- PostgreSQL running (via Docker Compose or local installation)
- Go 1.23+ installed

## Migration Commands

### 1. Run Migrations
Creates or updates the database schema:
```bash
cd services/user-service
go run cmd/migrate/main.go migrate
# or
make -f Makefile.migrate migrate
```

### 2. Seed Data
Adds initial test data to the database:
```bash
go run cmd/migrate/main.go seed
# or  
make -f Makefile.migrate seed
```

### 3. Check Status
Shows current database status:
```bash
go run cmd/migrate/main.go status
# or
make -f Makefile.migrate status
```

## Configuration

### For Docker Environment
Use `.env` file with container hostname:
```bash
DB_HOST=postgres  # Container name
```

### For Local Development
Use `.env.migrate` file with localhost:
```bash
DB_HOST=localhost  # Local development
```

## Database Schema

The migration system uses **Ent ORM** for schema management:

- **Schema Definition**: `ent/schema/user.go`
- **Generated Code**: `ent/` directory
- **Migration Logic**: `internal/infrastructure/database/database.go`

### User Table Structure
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);
```

## Troubleshooting

### Connection Issues
1. Ensure PostgreSQL is running:
   ```bash
   docker-compose ps
   ```

2. Check configuration:
   ```bash
   go run cmd/test-config/main.go
   ```

### Common Errors

**Error**: `lookup postgres: no such host`
- **Solution**: Use `.env.migrate` for local development

**Error**: `no such table: users`  
- **Solution**: Run migration first: `make migrate`

**Error**: `sqlite: foreign_keys pragma is off`
- **Solution**: Check DB_DRIVER is set to `postgres`

## Development Workflow

1. **Start Services**:
   ```bash
   docker-compose up -d
   ```

2. **Run Migrations**:
   ```bash
   cd services/user-service
   make -f Makefile.migrate migrate
   ```

3. **Seed Data**:
   ```bash
   make -f Makefile.migrate seed
   ```

4. **Verify**:
   ```bash
   make -f Makefile.migrate status
   ```

## Schema Changes

When modifying the database schema:

1. **Update Ent Schema**:
   ```bash
   # Edit ent/schema/user.go
   ```

2. **Generate Code**:
   ```bash
   go generate ./ent
   ```

3. **Run Migration**:
   ```bash
   make -f Makefile.migrate migrate
   ```

## Files Structure

```
services/user-service/
├── cmd/
│   ├── migrate/main.go          # Migration command tool  
│   └── test-config/main.go      # Configuration debugger
├── .env                         # Docker environment
├── .env.migrate                 # Local development environment
├── Makefile.migrate            # Migration shortcuts
└── MIGRATION_GUIDE.md          # This file
```
