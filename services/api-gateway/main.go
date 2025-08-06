package main

// @title API Gateway
// @version 1.0
// @description API Gateway for microservice architecture
// @host localhost:8080
// @BasePath /

import (
	"io"
	"log"
	"net/http"
	"time"

	"app-microservice/shared/config"

	"github.com/gin-gonic/gin"
)

const (
	USER_SERVICE_URL    = "http://localhost:8081"
	PRODUCT_SERVICE_URL = "http://localhost:8082"
)

func main() {
	cfg := config.LoadConfig()

	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// API Documentation Hub
	r.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
<!DOCTYPE html>
<html>
<head>
    <title>Microservice API Documentation</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .service { border: 1px solid #ddd; margin: 20px 0; padding: 20px; border-radius: 5px; }
        .service h2 { color: #333; }
        .service a { color: #007bff; text-decoration: none; font-weight: bold; }
        .service a:hover { text-decoration: underline; }
        .status { padding: 5px 10px; border-radius: 3px; font-size: 12px; }
        .up { background: #d4edda; color: #155724; }
        .down { background: #f8d7da; color: #721c24; }
    </style>
</head>
<body>
    <h1>🚀 Microservice API Documentation</h1>
    <p>Welcome to the API documentation hub for our microservice architecture.</p>
    
    <div class="service">
        <h2>🌐 API Gateway</h2>
        <p><strong>Port:</strong> 8080</p>
        <p><strong>Description:</strong> Main entry point that routes requests to microservices</p>
        <p><strong>Health Check:</strong> <a href="/health" target="_blank">Gateway Health</a></p>
        <p><strong>Services Status:</strong> <a href="/services/health" target="_blank">All Services Health</a></p>
    </div>

    <div class="service">
        <h2>👤 User Service</h2>
        <p><strong>Port:</strong> 8081</p>
        <p><strong>Description:</strong> Manages user data and operations</p>
        <p><strong>Swagger UI:</strong> <a href="http://localhost:8081/swagger/index.html" target="_blank">User Service API Docs</a></p>
        <p><strong>Direct Access:</strong> <a href="/users" target="_blank">GET /users</a></p>
    </div>

    <div class="service">
        <h2>📦 Product Service</h2>
        <p><strong>Port:</strong> 8082</p>
        <p><strong>Description:</strong> Manages product catalog and operations</p>
        <p><strong>Swagger UI:</strong> <a href="http://localhost:8082/swagger/index.html" target="_blank">Product Service API Docs</a></p>
        <p><strong>Direct Access:</strong> <a href="/products" target="_blank">GET /products</a></p>
    </div>

    <div class="service">
        <h2>🔗 Quick Links</h2>
        <ul>
            <li><a href="/users" target="_blank">All Users</a></li>
            <li><a href="/products" target="_blank">All Products</a></li>
            <li><a href="/services/health" target="_blank">System Health</a></li>
        </ul>
    </div>

    <div class="service">
        <h2>📋 API Examples</h2>
        <h3>Users API:</h3>
        <ul>
            <li>GET /users - Get all users</li>
            <li>GET /users/{id} - Get user by ID</li>
            <li>POST /users - Create new user</li>
            <li>PUT /users/{id} - Update user</li>
            <li>DELETE /users/{id} - Delete user</li>
        </ul>
        
        <h3>Products API:</h3>
        <ul>
            <li>GET /products - Get all products</li>
            <li>GET /products/{id} - Get product by ID</li>
            <li>POST /products - Create new product</li>
            <li>PUT /products/{id} - Update product</li>
            <li>DELETE /products/{id} - Delete product</li>
        </ul>
    </div>
</body>
</html>`)
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "api-gateway",
			"time":    time.Now(),
		})
	})

	// Routes to User Service
	r.Any("/users/*path", proxyToService(USER_SERVICE_URL))
	r.Any("/users", proxyToService(USER_SERVICE_URL))

	// Routes to Product Service
	r.Any("/products/*path", proxyToService(PRODUCT_SERVICE_URL))
	r.Any("/products", proxyToService(PRODUCT_SERVICE_URL))

	// Service health checks
	r.GET("/services/health", func(c *gin.Context) {
		userHealth := checkServiceHealth(USER_SERVICE_URL + "/health")
		productHealth := checkServiceHealth(PRODUCT_SERVICE_URL + "/health")

		c.JSON(http.StatusOK, gin.H{
			"gateway": "ok",
			"services": gin.H{
				"user-service":    userHealth,
				"product-service": productHealth,
			},
			"time": time.Now(),
		})
	})

	port := "8080"
	if cfg.Port != "" {
		port = cfg.Port
	}

	log.Printf("API Gateway đang chạy trên port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func proxyToService(serviceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Construct target URL
		targetURL := serviceURL + c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			targetURL += "?" + c.Request.URL.RawQuery
		}

		// Create new request
		req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create request",
				"details": err.Error(),
			})
			return
		}

		// Copy headers
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Make request to service
		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Service unavailable",
				"details": err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// Copy response body
		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}
}

func checkServiceHealth(healthURL string) string {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(healthURL)
	if err != nil {
		return "down"
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "up"
	}
	return "down"
}
