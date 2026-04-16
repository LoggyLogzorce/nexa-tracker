package priority

import (
	"errors"
	"net/http"
	"nexa-task-tracker/internal/pkg/events"
	"strconv"

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
	Title string `json:"title" binding:"required,min=1,max=50"`
	Color string `json:"color" binding:"omitempty"`
}

type UpdatePriorityRequest struct {
	Title *string `json:"title" binding:"omitempty,min=1,max=50"`
	Color *string `json:"color" binding:"omitempty"`
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
	var req CreatePriorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Создать модель приоритета
	priority := &Priority{
		ProjectID: projectID,
		Title:     req.Title,
		Color:     req.Color,
	}

	// 4. Вызвать сервис
	if err := h.service.Create(c.Request.Context(), priority); err != nil {
		// Обработка специфичных ошибок
		if errors.Is(err, ErrPriorityTitleExists) {
			response.Error(c, http.StatusConflict, "priority with this title already exists")
			return
		}
		if errors.Is(err, ErrColorFormat) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		// Общая ошибка
		response.Error(c, http.StatusInternalServerError, "failed to create priority")
		return
	}

	// 5. Вернуть созданный приоритет
	response.Success(c, http.StatusCreated, priority)
}

func (h *Handler) GetByProjectID(c *gin.Context) {
	// 1. Парсинг project ID из URL
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	// 2. Вызвать сервис
	priorities, err := h.service.GetByProjectID(c.Request.Context(), projectID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get priorities")
		return
	}

	// 3. Вернуть приоритеты
	response.Success(c, http.StatusOK, priorities)
}

func (h *Handler) Update(c *gin.Context) {
	// 1. Парсинг project ID из URL
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	// 2. Парсинг priority ID из URL
	priorityIDStr := c.Param("priority_id")
	priorityID, err := strconv.ParseUint(priorityIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid priority id")
		return
	}

	// 3. Парсинг и валидация запроса
	var req UpdatePriorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Проверить, что есть хотя бы одно поле для обновления
	if req.Title == nil && req.Color == nil {
		response.Error(c, http.StatusBadRequest, "no fields to update")
		return
	}

	// 5. Вызвать сервис
	priority, err := h.service.Update(c.Request.Context(), uint(priorityID), projectID, req)
	if err != nil {
		// Обработка специфичных ошибок
		if errors.Is(err, ErrPriorityNotFound) {
			response.Error(c, http.StatusNotFound, "priority not found")
			return
		}
		if errors.Is(err, ErrPriorityTitleExists) {
			response.Error(c, http.StatusConflict, "priority with this title already exists")
			return
		}
		if errors.Is(err, ErrColorFormat) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		// Общая ошибка
		response.Error(c, http.StatusInternalServerError, "failed to update priority")
		return
	}

	// 6. Вернуть обновленный приоритет
	response.Success(c, http.StatusOK, priority)
}

func (h *Handler) Delete(c *gin.Context) {
	// 1. Парсинг project ID из URL
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	// 2. Парсинг priority ID из URL
	priorityIDStr := c.Param("priority_id")
	priorityID, err := strconv.ParseUint(priorityIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid priority id")
		return
	}

	// 3. Вызвать сервис
	if err := h.service.Delete(c.Request.Context(), uint(priorityID), projectID); err != nil {
		// Обработка специфичных ошибок
		if errors.Is(err, ErrPriorityNotFound) {
			response.Error(c, http.StatusNotFound, "priority not found")
			return
		}
		// Общая ошибка
		response.Error(c, http.StatusInternalServerError, "failed to delete priority")
		return
	}

	// 4. Вернуть успешный ответ
	response.Success(c, http.StatusOK, gin.H{"message": "priority deleted successfully"})
}
