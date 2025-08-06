# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy shared modules
COPY shared/ ./shared/

# Copy service source
COPY services/api-gateway/ ./services/api-gateway/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./services/api-gateway

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/api-gateway .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./api-gateway"]
