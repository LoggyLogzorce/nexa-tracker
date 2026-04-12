package user

import (
	"errors"
	"net/http"
	"nexa-task-tracker/internal/ctxkeys"
	"nexa-task-tracker/internal/pkg/response"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,min=1,max=50"`
	Email string `json:"email" binding:"omitempty,email"`
}

type DeleteUserRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

func (h *Handler) GetMe(c *gin.Context) {
	// Get user_id from context (set by Auth middleware)
	userID, exists := c.Get(ctxkeys.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	// Get user from database
	user, err := h.service.GetByID(userID.(uuid.UUID))
	if err != nil {
		// Handle user not found (deleted after token issued)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "user not found")
			return
		}
		// Other database errors
		response.Error(c, http.StatusInternalServerError, "failed to get user")
		return
	}

	// Return user data via DTO
	response.Success(c, http.StatusOK, user.ToResponse())
}

func (h *Handler) UpdateMe(c *gin.Context) {
	// 1. Получить user_id из контекста
	userID, exists := c.Get(ctxkeys.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	// 2. Получить текущего пользователя из БД
	user, err := h.service.GetByID(userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "user not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to get user")
		return
	}

	// 3. Парсинг и валидация запроса
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	// 4. Проверить что хотя бы одно поле передано
	if req.Name == "" && req.Email == "" {
		response.Error(c, http.StatusBadRequest, "no fields to update")
		return
	}

	// 5. Обновить Name (если передан)
	if req.Name != "" {
		// Валидация имени на спецсимволы
		if !validateName(req.Name) {
			response.Error(c, http.StatusBadRequest, "name contains invalid characters")
			return
		}
		user.Name = req.Name
	}

	// 6. Обновить Email (если передан и отличается от текущего)
	if req.Email != "" {
		// Привести к lowercase для сравнения
		newEmail := strings.ToLower(req.Email)
		currentEmail := strings.ToLower(user.Email)

		if newEmail != currentEmail {
			// Проверить уникальность email
			exists, err := h.service.EmailExists(newEmail, user.ID)
			if err != nil {
				response.Error(c, http.StatusInternalServerError, "failed to check email availability")
				return
			}
			if exists {
				response.Error(c, http.StatusConflict, "email already in use")
				return
			}
			user.Email = newEmail
		}
	}

	// 7. Сохранить в БД (GORM автоматически обновит UpdatedAt)
	if err := h.service.Update(user); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update user")
		return
	}

	// 8. Вернуть обновленные данные
	response.Success(c, http.StatusOK, user.ToResponse())
}

// validateName checks if name contains only letters, spaces, and hyphens
func validateName(name string) bool {
	// Разрешаем буквы (любые языки), пробелы, дефисы
	matched, _ := regexp.MatchString(`^[\p{L}\s\-]+$`, name)
	return matched
}

func (h *Handler) DeleteMe(c *gin.Context) {
	// 1. Получить user_id из контекста
	userID, exists := c.Get(ctxkeys.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	// 2. Парсинг запроса с паролем
	var req DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	// 3. Удалить пользователя (с проверкой пароля внутри)
	if err := h.service.Delete(userID.(uuid.UUID), req.Password); err != nil {
		switch {
		case errors.Is(err, ErrInvalidPassword):
			response.Error(c, http.StatusUnauthorized, "invalid password")
		case errors.Is(err, ErrUserOwnsProjects):
			response.Error(c, http.StatusConflict, "cannot delete user with owned projects. Please transfer or delete your projects first")
		case errors.Is(err, ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "user not found")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to delete user")
		}
		return
	}

	// 4. Успешное удаление
	response.Success(c, http.StatusOK, gin.H{"message": "user deleted successfully"})
}
