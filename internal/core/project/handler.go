package project

import (
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
	// TODO: Get project by ID
	c.JSON(http.StatusOK, gin.H{"message": "get project endpoint"})
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
	// TODO: Update project
	c.JSON(http.StatusOK, gin.H{"message": "update project endpoint"})
}

func (h *Handler) Delete(c *gin.Context) {
	// TODO: Delete project
	c.JSON(http.StatusOK, gin.H{"message": "delete project endpoint"})
}
