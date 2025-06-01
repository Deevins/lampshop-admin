package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("my_super_secret_key")

const tokenTTL = time.Hour * 2

// Hardcoded-администратор (логин/пароль).
// В продакшене вы бы сверяли с базой, хэшировали пароль и т.д.
var adminUsername = "admin"
var adminPassword = "password123"

// ==============================
// Модели (модели, как и раньше)
// ==============================

// Category — модель категории
type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AttributeOption — модель опций атрибутов
type AttributeOption struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

// Product — товар с JSONB-атрибутами
type Product struct {
	ID          int                    `json:"id"`
	SKU         string                 `json:"sku"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	CategoryID  string                 `json:"categoryId"`
	IsActive    bool                   `json:"isActive"`
	ImageURL    string                 `json:"imageUrl"`
	Price       float64                `json:"price"`
	StockQty    int                    `json:"stockQty"`
	Attributes  map[string]interface{} `json:"attributes"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// OrderItem — элемент заказа
type OrderItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

// OrderStatus — возможные статусы заказа
type OrderStatus string

const (
	StatusPending    OrderStatus = "Pending"
	StatusProcessing OrderStatus = "Processing"
	StatusShipped    OrderStatus = "Shipped"
	StatusDelivered  OrderStatus = "Delivered"
)

// Order — модель заказа
type Order struct {
	ID           int         `json:"id"`
	CustomerName string      `json:"customerName"`
	Items        []OrderItem `json:"items"`
	TotalPrice   float64     `json:"totalPrice"`
	Status       OrderStatus `json:"status"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}

// UpdateStatusRequest — тело PUT /orders/:id/status
type UpdateStatusRequest struct {
	Status OrderStatus `json:"status"`
}

// LoginRequest — тело POST /login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse — ответ POST /login
type LoginResponse struct {
	Token string `json:"token"`
}

// ==============================
// In-Memory «База» (заглушка)
// ==============================

var (
	// Категории
	categories = []Category{
		{ID: "bulb", Name: "Лампочки"},
		{ID: "cable", Name: "Кабели"},
		{ID: "equipment", Name: "Оборудование"},
	}

	// Опции атрибутов
	attributeOptions = map[string][]AttributeOption{
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

	// Слайс товаров
	products      = []Product{}
	productsMutex = sync.RWMutex{}

	// Слайс заказов
	orders      = []Order{}
	ordersMutex = sync.RWMutex{}
)

func init() {
	// Инициализация одного товара и одного заказа «по умолчанию»
	p := Product{
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
	}
	products = append(products, p)

	o := Order{
		ID:           1,
		CustomerName: "Иван Иванов",
		Items: []OrderItem{
			{ProductID: 1, Quantity: 2},
		},
		TotalPrice: 1000,
		Status:     StatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	orders = append(orders, o)
}

// ==============================
// JWT-утилиты
// ==============================

// Создаёт JWT-токен с полем «sub» = username
func generateJWT(username string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Парсит и валидирует JWT, возвращает Subject (username) или ошибку
func parseJWT(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}
	return "", jwt.ErrTokenInvalidClaims
}

// ==============================
// Middleware для проверки JWT
// ==============================

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ожидаем заголовок Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется токен (Bearer)"})
			c.Abort()
			return
		}
		tokenString := authHeader[7:]
		username, err := parseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный или просроченный токен"})
			c.Abort()
			return
		}

		// В контексте Gin можно сохранить username, если нужно
		c.Set("username", username)
		c.Next()
	}
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // React/Vite по умолчанию на 3000
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ===== Незащищённые маршруты =====

	// 1) Эндпоинт для логина
	router.POST("/login", loginHandler)

	// 2) Эндпоинты для получения категорий/атрибутов — тоже могут стоять под защитой или открытыми
	//    Пока сделаем их открытыми (они нужны, чтобы заполнить dropdown на фронте).
	router.GET("/categories", getCategoriesHandler)
	router.GET("/categories/:id/attributes", getAttributeOptionsHandler)

	// ===== Защищённые маршруты (JWT) =====
	protected := router.Group("/")
	protected.Use(AuthMiddleware())

	// Products
	protected.GET("/products", getProductsHandler)
	protected.GET("/products/:id", getProductByIDHandler)
	protected.POST("/products", createProductHandler)
	protected.PUT("/products/:id", updateProductHandler)
	protected.DELETE("/products/:id", deleteProductHandler)

	// Orders
	protected.GET("/orders", getOrdersHandler)
	protected.PUT("/orders/:id/status", updateOrderStatusHandler)

	// Запускаем сервер
	router.Run(":8080")
}

// ==============================
// Handlers
// ==============================

// ===== Логин =====
func loginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}
	if req.Username != adminUsername || req.Password != adminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}
	token, err := generateJWT(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сгенерировать токен"})
		return
	}
	c.JSON(http.StatusOK, LoginResponse{Token: token})
}

// ===== Categories / Attributes =====

func getCategoriesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, categories)
}

func getAttributeOptionsHandler(c *gin.Context) {
	categoryID := c.Param("id")
	if opts, exists := attributeOptions[categoryID]; exists {
		c.JSON(http.StatusOK, opts)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Категория не найдена"})
	}
}

// ===== Products =====

func getProductsHandler(c *gin.Context) {
	productsMutex.RLock()
	defer productsMutex.RUnlock()
	list := make([]Product, len(products))
	copy(list, products)
	c.JSON(http.StatusOK, list)
}

func getProductByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}
	productsMutex.RLock()
	defer productsMutex.RUnlock()
	for _, p := range products {
		if p.ID == id {
			c.JSON(http.StatusOK, p)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
}

func createProductHandler(c *gin.Context) {
	var newProduct Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные товара"})
		return
	}
	productsMutex.Lock()
	defer productsMutex.Unlock()
	maxID := 0
	for _, p := range products {
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	newID := maxID + 1
	now := time.Now()
	created := Product{
		ID:          newID,
		SKU:         newProduct.SKU,
		Name:        newProduct.Name,
		Description: newProduct.Description,
		CategoryID:  newProduct.CategoryID,
		IsActive:    newProduct.IsActive,
		ImageURL:    newProduct.ImageURL,
		Price:       newProduct.Price,
		StockQty:    newProduct.StockQty,
		Attributes:  newProduct.Attributes,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	products = append(products, created)
	c.JSON(http.StatusCreated, created)
}

func updateProductHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}
	var upd Product
	if err := c.ShouldBindJSON(&upd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные товара"})
		return
	}
	productsMutex.Lock()
	defer productsMutex.Unlock()
	for idx, p := range products {
		if p.ID == id {
			products[idx].SKU = upd.SKU
			products[idx].Name = upd.Name
			products[idx].Description = upd.Description
			products[idx].CategoryID = upd.CategoryID
			products[idx].IsActive = upd.IsActive
			products[idx].ImageURL = upd.ImageURL
			products[idx].Price = upd.Price
			products[idx].StockQty = upd.StockQty
			products[idx].Attributes = upd.Attributes
			products[idx].UpdatedAt = time.Now()
			c.JSON(http.StatusOK, products[idx])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
}

func deleteProductHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}
	productsMutex.Lock()
	defer productsMutex.Unlock()
	for idx, p := range products {
		if p.ID == id {
			products = append(products[:idx], products[idx+1:]...)
			c.JSON(http.StatusOK, gin.H{"result": "успешно удалён"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
}

// ===== Orders =====

func getOrdersHandler(c *gin.Context) {
	ordersMutex.RLock()
	defer ordersMutex.RUnlock()
	list := make([]Order, len(orders))
	copy(list, orders)
	c.JSON(http.StatusOK, list)
}

func updateOrderStatusHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID заказа"})
		return
	}
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}
	ordersMutex.Lock()
	defer ordersMutex.Unlock()
	for idx, o := range orders {
		if o.ID == id {
			orders[idx].Status = req.Status
			orders[idx].UpdatedAt = time.Now()
			c.JSON(http.StatusOK, orders[idx])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Заказ не найден"})
}
