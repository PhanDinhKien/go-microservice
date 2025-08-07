# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy shared modules
COPY shared/ ./shared/

# Copy service source
COPY services/user-service/ ./services/user-service/

# Build the application
WORKDIR /app/services/user-service
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/user-service/user-service .

# Expose port
EXPOSE 8081

# Run the binary
CMD ["./user-service"]
