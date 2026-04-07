package participant

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

type AddParticipantRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
	Role   string `json:"role" binding:"required,oneof=owner member read_only"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=owner member read_only"`
}

func (h *Handler) AddParticipant(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"message": "add participant endpoint"})
}

func (h *Handler) GetByProjectID(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"message": "get participants by project endpoint"})
}

func (h *Handler) UpdateRole(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"message": "update participant role endpoint"})
}

func (h *Handler) RemoveParticipant(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"message": "remove participant endpoint"})
}
