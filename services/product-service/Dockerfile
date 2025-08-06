# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy shared modules
COPY shared/ ./shared/

# Copy service source
COPY services/product-service/ ./services/product-service/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o product-service ./services/product-service

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/product-service .

# Expose port
EXPOSE 8082

# Run the binary
CMD ["./product-service"]
