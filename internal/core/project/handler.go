package project

import (
	"errors"
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
	// 1. Получить userID из контекста (аутентификация)
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	// 2. Парсинг и валидация запроса
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Создать модель проекта
	project := &Project{
		Title:       req.Title,
		Description: req.Description,
	}

	// 4. Вызвать сервис
	if err := h.service.Create(c.Request.Context(), project, userID.(uuid.UUID)); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create project")
		return
	}

	// 5. Вернуть созданный проект
	response.Success(c, http.StatusCreated, project)
}

func (h *Handler) GetByID(c *gin.Context) {
	// 1. Парсинг project ID из URL
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	// 2. Получить userID из контекста (установлен middleware.Auth)
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	// 3. Вызвать сервис с проверкой прав
	project, err := h.service.GetByID(c.Request.Context(), projectID, userID.(uuid.UUID))
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
		response.Error(c, http.StatusInternalServerError, "failed to get project")
		return
	}

	// 4. Вернуть проект
	response.Success(c, http.StatusOK, project)
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
	// 1. Парсинг project ID из URL
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	// 2. Получить userID из контекста
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	// 3. Парсинг и валидация запроса
	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Создать модель проекта для обновления
	project := &Project{
		ID:          projectID,
		Title:       req.Title,
		Description: req.Description,
	}

	// 5. Вызвать сервис с проверкой прав
	if err := h.service.Update(c.Request.Context(), project, userID.(uuid.UUID)); err != nil {
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
		response.Error(c, http.StatusInternalServerError, "failed to update project")
		return
	}

	// 6. Вернуть обновленный проект
	response.Success(c, http.StatusOK, project)
}

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	if err := h.service.Delete(c.Request.Context(), projectID, userID.(uuid.UUID)); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to delete project")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "project deleted successfully"})
}
