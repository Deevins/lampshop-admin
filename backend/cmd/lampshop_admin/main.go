// main.go

package main

import (
	"bytes"
	"context"
	"github.com/Deevins/lampshop-admin-backend/internal/handler"
	"github.com/Deevins/lampshop-admin-backend/internal/infra"
	"github.com/Deevins/lampshop-admin-backend/internal/repository"
	"github.com/Deevins/lampshop-admin-backend/internal/service"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1) Подключаемся к PostgreSQL
	ctx := context.Background()
	dbPool, err := infra.NewPostgresPool(ctx)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// 2) Инициализируем репозитории и сервисы (для локальных данных: auth и categories)
	authRepo := repository.NewAuthRepository(dbPool)
	attrRepo := repository.NewAttributeRepository(dbPool)
	authService := service.NewAuthService(authRepo)
	categoryService := service.NewCategoryService(attrRepo)
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// 3) Инициализируем внешние HTTP-клиенты для order- и product-сервисов
	orderServiceURL := "http://localhost:8082"
	productServiceURL := "http://localhost:8081"

	orderClient := infra.NewOrderServiceClient(orderServiceURL)
	productClient := infra.NewProductServiceClient(productServiceURL)

	// 4) Настраиваем Gin + CORS
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 5) Открытые маршруты
	router.POST("/login", authHandler.Login)
	router.GET("/categories", categoryHandler.GetCategories)
	router.GET("/categories/:id/attributes", categoryHandler.GetAttributeOptions)

	// 6) Защищенные маршруты с JWT-миддлварой
	protected := router.Group("/")
	//protected.Use(AuthMiddleware())

	// ===== Order-service прокси-ручки =====
	protected.GET("/orders", func(c *gin.Context) {
		//token := extractBearerToken(c.Request.Header.Get("Authorization"))
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		resp, err := orderClient.GetAllOrders(token)
		infra.ProxyResponse(c, resp, err)
	})

	// ===== Order-service прокси-ручки =====
	protected.PUT("/orders/:id/status", func(c *gin.Context) {
		//token := extractBearerToken(c.Request.Header.Get("Authorization"))
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		orderID := c.Param("id")
		bodyBuf := new(bytes.Buffer)
		_, _ = io.Copy(bodyBuf, c.Request.Body)

		resp, err := orderClient.UpdateOrderStatus(token, orderID, bodyBuf)
		infra.ProxyResponse(c, resp, err)
	})

	protected.GET("/orders/active", func(c *gin.Context) {
		//token := extractBearerToken(c.Request.Header.Get("Authorization"))
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		resp, err := orderClient.GetActiveOrders(token)
		infra.ProxyResponse(c, resp, err)
	})
	protected.GET("/orders/:id/status", func(c *gin.Context) {
		//token := extractBearerToken(c.Request.Header.Get("Authorization"))
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		orderID := c.Param("id")
		resp, err := orderClient.GetOrderStatus(token, orderID)
		infra.ProxyResponse(c, resp, err)
	})
	protected.POST("/checkout", func(c *gin.Context) {
		//token := extractBearerToken(c.Request.Header.Get("Authorization"))
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		bodyBuf := new(bytes.Buffer)
		_, _ = io.Copy(bodyBuf, c.Request.Body)
		resp, err := orderClient.CreateOrder(token, bodyBuf)
		infra.ProxyResponse(c, resp, err)
	})

	// ===== Product-service прокси-ручки =====
	protected.GET("/products", func(c *gin.Context) {
		//token := extractBearerToken(c.Request.Header.Get("Authorization"))
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		resp, err := productClient.ListProducts(token)
		infra.ProxyResponse(c, resp, err)
	})
	protected.POST("/products", func(c *gin.Context) {
		//token := extractBearerToken(c.Request.Header.Get("Authorization"))
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		bodyBuf := new(bytes.Buffer)
		_, _ = io.Copy(bodyBuf, c.Request.Body)
		resp, err := productClient.CreateProduct(token, bodyBuf)
		infra.ProxyResponse(c, resp, err)
	})
	protected.GET("/products/:id", func(c *gin.Context) {
		//token := extractBearerToken(c.Request.Header.Get("Authorization"))
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		productID := c.Param("id")
		resp, err := productClient.GetProductByID(token, productID)
		infra.ProxyResponse(c, resp, err)
	})
	protected.PUT("/products/:id", func(c *gin.Context) {
		//token := extractBearerToken(c.Request.Header.Get("Authorization"))
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		productID := c.Param("id")
		bodyBuf := new(bytes.Buffer)
		_, _ = io.Copy(bodyBuf, c.Request.Body)
		resp, err := productClient.UpdateProduct(token, productID, bodyBuf)
		infra.ProxyResponse(c, resp, err)
	})
	protected.DELETE("/products/:id", func(c *gin.Context) {
		//token := extractBearerToken(c.Request.Header.Get("Authorization"))
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		productID := c.Param("id")
		resp, err := productClient.DeleteProduct(token, productID)
		infra.ProxyResponse(c, resp, err)
	})

	protected.GET("/products/upload-url", func(c *gin.Context) {
		//token := extractBearerToken(c.Request.Header.Get("Authorization"))
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		resp, err := productClient.GetUploadURL(token)
		infra.ProxyResponse(c, resp, err)
	})
	protected.POST("/products/notify-upload", func(c *gin.Context) {
		//token := extractBearerToken(c.Request.Header.Get("Authorization"))
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		bodyBuf := new(bytes.Buffer)
		_, _ = io.Copy(bodyBuf, c.Request.Body)
		resp, err := productClient.NotifyUpload(token, bodyBuf)
		infra.ProxyResponse(c, resp, err)
	})

	router.Run(":" + "8083")
}

// AuthMiddleware проверяет заголовок Authorization Bearer <token>
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
			c.Abort()
			return
		}
		tokenString := authHeader[7:]
		username, err := infra.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}
		c.Set("username", username)
		c.Next()
	}
}

// extractBearerToken убирает "Bearer " префикс и возвращает сам токен.
func extractBearerToken(headerValue string) string {
	if len(headerValue) > 7 && headerValue[:7] == "Bearer " {
		return headerValue[7:]
	}
	return ""
}
