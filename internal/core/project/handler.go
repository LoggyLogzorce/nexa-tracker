package project

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
	// TODO: List projects
	c.JSON(http.StatusOK, gin.H{"message": "list projects endpoint"})
}

func (h *Handler) Update(c *gin.Context) {
	// TODO: Update project
	c.JSON(http.StatusOK, gin.H{"message": "update project endpoint"})
}

func (h *Handler) Delete(c *gin.Context) {
	// TODO: Delete project
	c.JSON(http.StatusOK, gin.H{"message": "delete project endpoint"})
}
