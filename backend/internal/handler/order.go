package handler

import (
	"github.com/Deevins/lampshop-admin-backend/internal/domain"
	"github.com/Deevins/lampshop-admin-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OrderHandler обрабатывает запросы, связанные с заказами.
type OrderHandler struct {
	orderService service.OrderService
}

// NewOrderHandler создаёт новый OrderHandler.
func NewOrderHandler(os service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: os}
}

// GetOrders – GET /orders
func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

type updateStatusRequest struct {
	Status domain.OrderStatus `json:"status"`
}

// UpdateOrderStatus – PUT /orders/:id/status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}
	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	updated, err := h.orderService.UpdateOrderStatus(id, req.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, updated)
}
