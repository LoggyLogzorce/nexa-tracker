package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,min=1,max=50"`
	Email string `json:"email" binding:"omitempty,email"`
}

func (h *Handler) GetMe(c *gin.Context) {
	// TODO: Get current user from context
	c.JSON(http.StatusOK, gin.H{"message": "get me endpoint"})
}

func (h *Handler) UpdateMe(c *gin.Context) {
	// TODO: Update current user
	c.JSON(http.StatusOK, gin.H{"message": "update me endpoint"})
}

func (h *Handler) DeleteMe(c *gin.Context) {
	// TODO: Delete current user
	c.JSON(http.StatusOK, gin.H{"message": "delete me endpoint"})
}
