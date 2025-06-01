package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)


type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AttributeOption — описание доступного атрибута в категории
type AttributeOption struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

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

// OrderItem — элемент в заказе, ссылается на Product.ID и количество
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

type UpdateStatusRequest struct {
	Status OrderStatus `json:"status"`
}

// ---------------------------------------------------
//     In-Memory “База данных” (заглушка)
// ---------------------------------------------------

var (
	// Категории (имитируем, что эти данные брались бы из БД)
	categories = []Category{
		{ID: "bulb", Name: "Лампочки"},
		{ID: "cable", Name: "Кабели"},
		{ID: "equipment", Name: "Оборудование"},
	}

	// Какие атрибуты есть у каждой категории (имитируем, что в БД есть таблица attribute_options)
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
	productsMutex = sync.RWMutex{} // для безопасного доступа из нескольких горутин

	// Слайс заказов
	orders      = []Order{}
	ordersMutex = sync.RWMutex{}
)

func init() {
	// Инициализируем один товар и один заказ “по умолчанию” для примера

	// 1. Товар
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

	// 2. Заказ
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

// ---------------------------------------------------
//     Основная функция — настройка Gin и маршрутов
// ---------------------------------------------------

func main() {
	router := gin.Default()
	router.Use(corsMiddleware())

	// ========== Категории и атрибуты ==========
	router.GET("/categories", getCategoriesHandler)
	router.GET("/categories/:id/attributes", getAttributeOptionsHandler)

	// ========== CRUD для товаров ==========
	router.GET("/products", getProductsHandler)
	router.GET("/products/:id", getProductByIDHandler)
	router.POST("/products", createProductHandler)
	router.PUT("/products/:id", updateProductHandler)
	router.DELETE("/products/:id", deleteProductHandler)

	// ========== Просмотр и изменение статуса заказов ==========
	router.GET("/orders", getOrdersHandler)
	router.PUT("/orders/:id/status", updateOrderStatusHandler)

	// Запускаем на порту 8080
	router.Run(":8080")
}

// ---------------------------------------------------
//     Handlers для Категорий и Атрибутов
// ---------------------------------------------------

// getCategoriesHandler возвращает все категории
func getCategoriesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, categories)
}

// getAttributeOptionsHandler возвращает список AttributeOption для данной categoryId
func getAttributeOptionsHandler(c *gin.Context) {
	categoryID := c.Param("id")
	if opts, exists := attributeOptions[categoryID]; exists {
		c.JSON(http.StatusOK, opts)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Категория не найдена"})
	}
}

// ---------------------------------------------------
//     Handlers для Products
// ---------------------------------------------------

// getProductsHandler — GET /products
func getProductsHandler(c *gin.Context) {
	productsMutex.RLock()
	defer productsMutex.RUnlock()

	// Возвращаем копию, чтобы внешний код не мутировал наш in-memory-слайс
	list := make([]Product, len(products))
	copy(list, products)
	c.JSON(http.StatusOK, list)
}

// getProductByIDHandler — GET /products/:id
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

// createProductHandler — POST /products
func createProductHandler(c *gin.Context) {
	var newProduct Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные товара"})
		return
	}

	productsMutex.Lock()
	defer productsMutex.Unlock()

	// Сгенерируем новый ID как max(ID)+1
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

// updateProductHandler — PUT /products/:id
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
			// Меняем только те поля, которые пришли из JSON
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

// deleteProductHandler — DELETE /products/:id
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
			// Удаляем элемент из слайса
			products = append(products[:idx], products[idx+1:]...)
			c.JSON(http.StatusOK, gin.H{"result": "успешно удалён"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
}

// ---------------------------------------------------
//     Handlers для Orders
// ---------------------------------------------------

// getOrdersHandler — GET /orders
func getOrdersHandler(c *gin.Context) {
	ordersMutex.RLock()
	defer ordersMutex.RUnlock()

	// Возвращаем копию
	list := make([]Order, len(orders))
	copy(list, orders)
	c.JSON(http.StatusOK, list)
}

// updateOrderStatusHandler — PUT /orders/:id/status
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

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
