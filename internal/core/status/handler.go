package status

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"nexa-task-tracker/internal/pkg/events"
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

type CreateStatusRequest struct {
	Name       string `json:"name" binding:"required"`
	Color      string `json:"color"`
	OrderIndex int    `json:"order_index"`
}

type UpdateStatusRequest struct {
	Name       string `json:"name"`
	Color      string `json:"color"`
	OrderIndex int    `json:"order_index"`
}

func (h *Handler) Create(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"message": "create status endpoint"})
}

func (h *Handler) GetByProjectID(c *gin.Context) {
	idStr := c.Param("project_id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	statuses, err := h.service.GetByProjectID(projectID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get statuses")
		return
	}

	response.Success(c, http.StatusOK, statuses)
}

func (h *Handler) Update(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"message": "update status endpoint"})
}

func (h *Handler) Delete(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"message": "delete status endpoint"})
}
