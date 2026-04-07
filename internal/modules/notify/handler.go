package notify

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

func (h *Handler) List(c *gin.Context) {
	// TODO: Get notifications for current user
	c.JSON(http.StatusOK, gin.H{"message": "list notifications endpoint"})
}

func (h *Handler) MarkAsRead(c *gin.Context) {
	// TODO: Mark notification as read
	c.JSON(http.StatusOK, gin.H{"message": "mark as read endpoint"})
}

func (h *Handler) MarkAllAsRead(c *gin.Context) {
	// TODO: Mark all notifications as read
	c.JSON(http.StatusOK, gin.H{"message": "mark all as read endpoint"})
}
