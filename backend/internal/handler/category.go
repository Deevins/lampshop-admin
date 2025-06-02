package handler

import (
	"github.com/Deevins/lampshop-admin-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CategoryHandler обрабатывает запросы, связанные с категориями и атрибутами.
type CategoryHandler struct {
	categoryService service.CategoryService
}

// NewCategoryHandler создаёт новый CategoryHandler.
func NewCategoryHandler(cs service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: cs}
}

// GetCategories – GET /categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	cats, err := h.categoryService.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch categories"})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// GetAttributeOptions – GET /categories/:id/attributes
func (h *CategoryHandler) GetAttributeOptions(c *gin.Context) {
	id := c.Param("id")
	opts, err := h.categoryService.GetAttributeOptions(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "attributes not found for category"})
		return
	}
	c.JSON(http.StatusOK, opts)
}
