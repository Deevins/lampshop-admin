package handler

import (
	"github.com/Deevins/lampshop-admin-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler обрабатывает логин.
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler создаёт новый AuthHandler.
func NewAuthHandler(as service.AuthService) *AuthHandler {
	return &AuthHandler{authService: as}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type loginResponse struct {
	Token string `json:"token"`
}

// Login – POST /login
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, loginResponse{Token: token})
}
