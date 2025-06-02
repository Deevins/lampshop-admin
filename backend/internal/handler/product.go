package handler

import (
	"github.com/Deevins/lampshop-admin-backend/internal/domain"
	"github.com/Deevins/lampshop-admin-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProductHandler обрабатывает CRUD-операции для товаров.
type ProductHandler struct {
	productService service.ProductService
}

// NewProductHandler создаёт новый ProductHandler.
func NewProductHandler(ps service.ProductService) *ProductHandler {
	return &ProductHandler{productService: ps}
}

// GetProducts – GET /products
func (h *ProductHandler) GetProducts(c *gin.Context) {
	prods, err := h.productService.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch products"})
		return
	}
	c.JSON(http.StatusOK, prods)
}

// GetProductByID – GET /products/:id
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}
	prod, err := h.productService.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, prod)
}

// CreateProduct – POST /products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var input domain.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	created, err := h.productService.CreateProduct(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create product"})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// UpdateProduct – PUT /products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}
	var input domain.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	updated, err := h.productService.UpdateProduct(id, input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteProduct – DELETE /products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}
	err = h.productService.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}
