package project

import (
	"github.com/google/uuid"
	"net/http"
	"nexa-task-tracker/internal/middleware"
	"nexa-task-tracker/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type CreateProjectRequest struct {
	Title       string  `json:"title" binding:"required,min=1,max=50"`
	Description *string `json:"description" binding:"omitempty,max=255"`
}

type UpdateProjectRequest struct {
	Title       string  `json:"title" binding:"omitempty,min=1,max=50"`
	Description *string `json:"description" binding:"omitempty,max=255"`
}

func (h *Handler) Create(c *gin.Context) {
	// TODO: Implement project creation
	c.JSON(http.StatusOK, gin.H{"message": "create project endpoint"})
}

func (h *Handler) GetByID(c *gin.Context) {
	// TODO: Get project by ID
	c.JSON(http.StatusOK, gin.H{"message": "get project endpoint"})
}

func (h *Handler) List(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	projects, err := h.service.List(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get projects")
		return
	}

	response.Success(c, http.StatusOK, projects)
}

func (h *Handler) Update(c *gin.Context) {
	// TODO: Update project
	c.JSON(http.StatusOK, gin.H{"message": "update project endpoint"})
}

func (h *Handler) Delete(c *gin.Context) {
	// TODO: Delete project
	c.JSON(http.StatusOK, gin.H{"message": "delete project endpoint"})
}
