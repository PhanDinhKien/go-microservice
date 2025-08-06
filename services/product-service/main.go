package main

// @title Product Service API
// @version 1.0
// @description Product Service API for microservice architecture
// @host localhost:8082
// @BasePath /

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"app-microservice/shared/config"
	"app-microservice/shared/models"

	_ "app-microservice/services/product-service/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var products = []models.Product{
	{ID: 1, Name: "Laptop Dell", Description: "Laptop Dell Inspiron 15", Price: 15000000, Stock: 10, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Name: "iPhone 14", Description: "Apple iPhone 14 Pro", Price: 25000000, Stock: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

func main() {
	cfg := config.LoadConfig()

	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// Routes
	api := r.Group("/products")
	{
		api.GET("", getProducts)
		api.GET("/:id", getProductByID)
		api.POST("", createProduct)
		api.PUT("/:id", updateProduct)
		api.DELETE("/:id", deleteProduct)
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "product-service",
			"time":    time.Now(),
		})
	})

	port := "8082"
	if cfg.Port != "" && cfg.Port != "8080" && cfg.Port != "8081" {
		port = cfg.Port
	}

	log.Printf("Product Service đang chạy trên port %s", port)
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

// GetProducts godoc
// @Summary Get all products
// @Description Get list of all products
// @Tags products
// @Produce json
// @Success 200 {object} models.Response
// @Router /products [get]
func getProducts(c *gin.Context) {
	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Lấy danh sách products thành công",
		Data:    products,
	})
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get a product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /products/{id} [get]
func getProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "ID không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	for _, product := range products {
		if product.ID == id {
			c.JSON(http.StatusOK, models.Response{
				Status:  "success",
				Message: "Lấy thông tin product thành công",
				Data:    product,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, models.ErrorResponse{
		Status:  "error",
		Message: "Không tìm thấy product",
		Error:   "Product not found",
	})
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the given data
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product data"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Router /products [post]
func createProduct(c *gin.Context) {
	var newProduct models.Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Dữ liệu không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	// Generate new ID
	newProduct.ID = len(products) + 1
	newProduct.CreatedAt = time.Now()
	newProduct.UpdatedAt = time.Now()

	products = append(products, newProduct)

	c.JSON(http.StatusCreated, models.Response{
		Status:  "success",
		Message: "Tạo product thành công",
		Data:    newProduct,
	})
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update product data by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Updated product data"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /products/{id} [put]
func updateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "ID không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Dữ liệu không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	for i, product := range products {
		if product.ID == id {
			updatedProduct.ID = id
			updatedProduct.CreatedAt = product.CreatedAt
			updatedProduct.UpdatedAt = time.Now()
			products[i] = updatedProduct

			c.JSON(http.StatusOK, models.Response{
				Status:  "success",
				Message: "Cập nhật product thành công",
				Data:    updatedProduct,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, models.ErrorResponse{
		Status:  "error",
		Message: "Không tìm thấy product",
		Error:   "Product not found",
	})
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /products/{id} [delete]
func deleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "ID không hợp lệ",
			Error:   err.Error(),
		})
		return
	}

	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...)
			c.JSON(http.StatusOK, models.Response{
				Status:  "success",
				Message: "Xóa product thành công",
				Data:    nil,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, models.ErrorResponse{
		Status:  "error",
		Message: "Không tìm thấy product",
		Error:   "Product not found",
	})
}
