package status

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"nexa-task-tracker/internal/pkg/events"
	"nexa-task-tracker/internal/pkg/response"
	"strings"
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
	Name       string `json:"name" binding:"required,min=1,max=50"`
	Color      string `json:"color" binding:"omitempty"`
	OrderIndex int    `json:"order_index" binding:"omitempty,min=0"`
}

type UpdateStatusRequest struct {
	Name       string `json:"name"`
	Color      string `json:"color"`
	OrderIndex int    `json:"order_index"`
}

func (h *Handler) Create(c *gin.Context) {
	// 1. Парсинг project ID из URL
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	// 2. Парсинг и валидация запроса
	var req CreateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Создать модель статуса
	status := &Status{
		ProjectID:  projectID,
		Name:       req.Name,
		Color:      req.Color,
		OrderIndex: req.OrderIndex,
	}

	// 4. Вызвать сервис
	if err := h.service.Create(c.Request.Context(), status); err != nil {
		// Обработка специфичных ошибок
		if errors.Is(err, ErrStatusNameExists) {
			response.Error(c, http.StatusConflict, "status with this name already exists")
			return
		}
		if errors.Is(err, ErrDuplicateOrderIndex) {
			response.Error(c, http.StatusBadRequest, "status with this order_index already exists")
			return
		}
		// Проверка на ошибку валидации hex-цвета
		if strings.Contains(err.Error(), "invalid color format") {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		// Общая ошибка
		response.Error(c, http.StatusInternalServerError, "failed to create status")
		return
	}

	// 5. Вернуть созданный статус
	response.Success(c, http.StatusCreated, status)
}

func (h *Handler) GetByProjectID(c *gin.Context) {
	// 1. Парсинг project ID из URL
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	// 3. Вызвать сервис с проверкой прав
	statuses, err := h.service.GetByProjectID(c.Request.Context(), projectID)
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
