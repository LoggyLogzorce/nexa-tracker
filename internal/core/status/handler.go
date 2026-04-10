package status

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"nexa-task-tracker/internal/pkg/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
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
