package status

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"nexa-task-tracker/internal/pkg/events"
	"nexa-task-tracker/internal/pkg/response"
	"strconv"
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
	Name       *string `json:"name" binding:"omitempty,min=1,max=50"`
	Color      *string `json:"color" binding:"omitempty"`
	OrderIndex *int    `json:"order_index" binding:"omitempty,min=0"`
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
	// 1. Парсинг project ID из URL
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	// 2. Парсинг status ID из URL
	statusIDStr := c.Param("status_id")
	statusID, err := strconv.ParseUint(statusIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid status id")
		return
	}

	// 3. Парсинг и валидация запроса
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Проверить, что есть хотя бы одно поле для обновления
	if req.Color == nil && req.Name == nil && req.OrderIndex == nil {
		response.Error(c, http.StatusBadRequest, "no fields to update")
		return
	}

	// 5. Вызвать сервис
	status, err := h.service.Update(c.Request.Context(), uint(statusID), projectID, req)
	if err != nil {
		// Обработка специфичных ошибок
		if errors.Is(err, ErrStatusNotFound) {
			response.Error(c, http.StatusNotFound, "status not found")
			return
		}
		if errors.Is(err, ErrStatusNameExists) {
			response.Error(c, http.StatusConflict, "status with this name already exists")
			return
		}
		// Проверка на ошибку валидации hex-цвета
		if errors.Is(err, ErrColorFormat) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		// Общая ошибка
		response.Error(c, http.StatusInternalServerError, "failed to update status")
		return
	}

	// 6. Вернуть обновленный статус
	response.Success(c, http.StatusOK, status)
}

func (h *Handler) Delete(c *gin.Context) {
	// 1. Парсинг project ID из URL
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	// 2. Парсинг status ID из URL
	statusIDStr := c.Param("status_id")
	statusID, err := strconv.ParseUint(statusIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid status id")
		return
	}

	// 3. Вызвать сервис
	if err := h.service.Delete(c.Request.Context(), uint(statusID), projectID); err != nil {
		// Обработка специфичных ошибок
		if errors.Is(err, ErrStatusNotFound) {
			response.Error(c, http.StatusNotFound, "status not found")
			return
		}
		// Общая ошибка
		response.Error(c, http.StatusInternalServerError, "failed to delete status")
		return
	}

	// 4. Вернуть успешный ответ
	response.Success(c, http.StatusOK, gin.H{"message": "status deleted successfully"})
}
