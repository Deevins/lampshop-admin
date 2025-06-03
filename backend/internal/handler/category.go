package handler

import (
	"github.com/Deevins/lampshop-admin-backend/internal/service"
	"github.com/google/uuid"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(cs service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: cs}
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	cats, err := h.categoryService.GetAllCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch categories"})
		return
	}
	if len(cats) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no categories found"})
		return
	}
	c.JSON(http.StatusOK, cats)
}

func (h *CategoryHandler) GetAttributeOptions(c *gin.Context) {
	idParam := c.Param("id")
	categoryID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}
	opts, err := h.categoryService.GetAttributeOptions(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "attributes not found for category"})
		return
	}
	c.JSON(http.StatusOK, opts)
}
