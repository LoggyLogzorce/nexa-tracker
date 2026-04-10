package priority

import (
	"net/http"
	"nexa-task-tracker/internal/pkg/events"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"nexa-task-tracker/internal/pkg/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service, eventBus *events.EventBus) *Handler {
	eventBus.Subscribe(events.ProjectDeleted, service.HandleProjectDeleted)
	eventBus.Subscribe(events.ProjectCreated, service.HandleProjectCreated)
	return &Handler{service: service}
}

type CreatePriorityRequest struct {
	Title string `json:"title" binding:"required"`
	Color string `json:"color"`
}

type UpdatePriorityRequest struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

func (h *Handler) Create(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"message": "create priority endpoint"})
}

func (h *Handler) GetByProjectID(c *gin.Context) {
	idStr := c.Param("project_id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	priorities, err := h.service.GetByProjectID(projectID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get priorities")
		return
	}

	response.Success(c, http.StatusOK, priorities)
}

func (h *Handler) Update(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"message": "update priority endpoint"})
}

func (h *Handler) Delete(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"message": "delete priority endpoint"})
}
