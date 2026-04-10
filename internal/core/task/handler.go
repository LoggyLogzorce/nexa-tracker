package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required,min=1,max=100"`
	Description *string   `json:"description"`
	ProjectID   uuid.UUID `json:"project_id" binding:"required"`
	StatusID    *uint     `json:"status_id"`
	PriorityID  *uint     `json:"priority_id"`
	AssigneeID  *string   `json:"assignee_id" binding:"omitempty,uuid"`
	Deadline    *string   `json:"deadline" binding:"omitempty"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title" binding:"omitempty,min=1,max=100"`
	Description *string `json:"description"`
	StatusID    *uint   `json:"status_id"`
	PriorityID  *uint   `json:"priority_id"`
	AssigneeID  *string `json:"assignee_id" binding:"omitempty,uuid"`
	Deadline    *string `json:"deadline"`
}

func (h *Handler) Create(c *gin.Context) {
	// TODO: Implement task creation
	c.JSON(http.StatusOK, gin.H{"message": "create task endpoint"})
}

func (h *Handler) GetByID(c *gin.Context) {
	// TODO: Get task by ID
	c.JSON(http.StatusOK, gin.H{"message": "get task endpoint"})
}

func (h *Handler) List(c *gin.Context) {
	// TODO: List tasks with filters
	c.JSON(http.StatusOK, gin.H{"message": "list tasks endpoint"})
}

func (h *Handler) Update(c *gin.Context) {
	// TODO: Update task
	c.JSON(http.StatusOK, gin.H{"message": "update task endpoint"})
}

func (h *Handler) Delete(c *gin.Context) {
	// TODO: Delete task
	c.JSON(http.StatusOK, gin.H{"message": "delete task endpoint"})
}
