package status

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"nexa-task-tracker/internal/ctxkeys"
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
	// 1. Парсинг project ID из URL
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	// 2. Получить userID из контекста
	userID, exists := c.Get(ctxkeys.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	// 3. Вызвать сервис с проверкой прав
	statuses, err := h.service.GetByProjectID(c.Request.Context(), projectID, userID.(uuid.UUID))
	if err != nil {
		// Обработка специфичных ошибок
		if errors.Is(err, ErrProjectNotFound) {
			response.Error(c, http.StatusNotFound, "project not found")
			return
		}
		if errors.Is(err, ErrProjectAccessDenied) {
			response.Error(c, http.StatusForbidden, "access denied")
			return
		}
		// Общая ошибка
		response.Error(c, http.StatusInternalServerError, "failed to get statuses")
		return
	}

	// 4. Вернуть статусы
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
