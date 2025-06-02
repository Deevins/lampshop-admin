package main

import (
	"github.com/Deevins/lampshop-admin-backend/internal/domain"
	"github.com/Deevins/lampshop-admin-backend/internal/handler"
	"github.com/Deevins/lampshop-admin-backend/internal/infra"
	"github.com/Deevins/lampshop-admin-backend/internal/repository"
	"github.com/Deevins/lampshop-admin-backend/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// ============================
	// 1) Инициализация “базы” в memory
	// ============================

	// Категории
	initialCategories := []domain.Category{
		{ID: "bulb", Name: "Лампочки"},
		{ID: "cable", Name: "Кабели"},
		{ID: "equipment", Name: "Оборудование"},
	}
	// Опции атрибутов (JSONB-структура)
	initialAttributeOptions := map[string][]domain.AttributeOption{
		"bulb": {
			{Key: "power", Label: "Мощность (Вт)", Type: "number"},
			{Key: "color", Label: "Цвет", Type: "text"},
			{Key: "temperature", Label: "Температура (K)", Type: "number"},
			{Key: "socketType", Label: "Тип цоколя", Type: "text"},
		},
		"cable": {
			{Key: "length", Label: "Длина (м)", Type: "number"},
			{Key: "material", Label: "Материал", Type: "text"},
			{Key: "color", Label: "Цвет", Type: "text"},
		},
		"equipment": {
			{Key: "manufacturer", Label: "Производитель", Type: "text"},
			{Key: "model", Label: "Модель", Type: "text"},
			{Key: "warranty", Label: "Гарантия (мес.)", Type: "number"},
		},
	}

	// Один пример товара
	initialProducts := []domain.Product{
		{
			ID:          1,
			SKU:         "BULB-007",
			Name:        "EcoBright 7W",
			Description: "Энергоэффективная лампочка для дома.",
			CategoryID:  "bulb",
			IsActive:    true,
			ImageURL:    "https://santhimetaleshop.in/cdn/shop/files/Untitleddesign_26a5d7f4-82b7-4e7a-ac43-068a31086beb.png?v=1694498000&width=1445",
			Price:       500,
			StockQty:    20,
			Attributes: map[string]interface{}{
				"power":       7,
				"color":       "Тёплый белый",
				"temperature": 2700,
				"socketType":  "E27",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Один пример заказа
	initialOrders := []domain.Order{
		{
			ID:           1,
			CustomerName: "Иван Иванов",
			Items: []domain.OrderItem{
				{ProductID: 1, Quantity: 2},
			},
			TotalPrice: 1000,
			Status:     domain.StatusPending,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	// ============================
	// 2) Репозитории
	// ============================
	catRepo := repository.NewInMemoryCategoryRepo(initialCategories)
	attrRepo := repository.NewInMemoryAttributeRepo(initialAttributeOptions)
	prodRepo := repository.NewInMemoryProductRepo(initialProducts)
	orderRepo := repository.NewInMemoryOrderRepo(initialOrders)

	// ============================
	// 3) Сервисы
	// ============================
	authSrv := service.NewAuthService()
	categorySrv := service.NewCategoryService(catRepo, attrRepo)
	productSrv := service.NewProductService(prodRepo)
	orderSrv := service.NewOrderService(orderRepo)

	// ============================
	// 4) Внешние клиенты (Infrastructure)
	// ============================
	externalClient := infra.NewInMemoryExternalClient()

	// ============================
	// 5) Хендлеры
	// ============================
	authHandler := handler.NewAuthHandler(authSrv)
	categoryHandler := handler.NewCategoryHandler(categorySrv)
	productHandler := handler.NewProductHandler(productSrv)
	orderHandler := handler.NewOrderHandler(orderSrv)

	// ============================
	// 6) Gin, CORS, Middleware
	// ============================
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Открытые маршруты
	router.POST("/login", authHandler.Login)
	router.GET("/categories", categoryHandler.GetCategories)
	router.GET("/categories/:id/attributes", categoryHandler.GetAttributeOptions)

	// Группа защищённых маршрутов: JWT-middleware
	protected := router.Group("/")
	protected.Use(AuthMiddleware(externalClient))

	// --- CRUD для продуктов ---
	protected.GET("/products", productHandler.GetProducts)
	protected.GET("/products/:id", productHandler.GetProductByID)
	protected.POST("/products", productHandler.CreateProduct)
	protected.PUT("/products/:id", productHandler.UpdateProduct)
	protected.DELETE("/products/:id", productHandler.DeleteProduct)

	// --- CRUD для заказов (только чтение + обновление статуса) ---
	protected.GET("/orders", orderHandler.GetOrders)
	protected.PUT("/orders/:id/status", func(c *gin.Context) {
		var req struct {
			Status domain.OrderStatus `json:"status"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
			return
		}
		updated, err := orderSrv.UpdateOrderStatus(id, req.Status)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		// Вызов внешнего сервиса — “уведомление” о смене статуса
		_ = externalClient.NotifyOrderStatusChange(updated.ID, string(updated.Status))
		c.JSON(http.StatusOK, updated)
	})

	// Запускаем HTTP-сервер на порту 8080
	router.Run(":8083")
}

// AuthMiddleware проверяет в заголовке Authorization Bearer <token> и парсит JWT.
// Если токен невалидный или отсутствует — возвращает 401.
func AuthMiddleware(externalClient infra.ExternalClient) gin.HandlerFunc {
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
		// При желании можно сохранить username в контексте Gin:
		c.Set("username", username)

		// Здесь можно вызывать externalClient, если требуется:
		// externalClient.SomeMethod(...)
		c.Next()
	}
}
