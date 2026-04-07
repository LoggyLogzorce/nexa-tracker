package auth

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"nexa-task-tracker/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// setCookie устанавливает refresh token cookie с правильными security настройками
func setCookie(c *gin.Context, name, value, domain string, sameSite http.SameSite, maxAge int) {
	// Определить Secure flag в зависимости от окружения
	secure := true
	if os.Getenv("ENV") == "development" {
		secure = false
		domain = ""
	}

	c.SetSameSite(sameSite)
	c.SetCookie(
		name,
		value,
		maxAge,
		"/",
		domain,
		secure,
		true,
	)
}

// deleteCookie удаляет cookie
func deleteCookie(c *gin.Context, name, domain string, sameSite http.SameSite) {
	secure := true
	if os.Getenv("ENV") == "development" {
		secure = false
		domain = ""
	}

	c.SetSameSite(sameSite)
	c.SetCookie(
		name,
		"",
		-1, // maxAge = -1 удаляет cookie
		"/",
		domain,
		secure,
		true,
	)
}

type Handler struct {
	service       Service
	domain        string
	sameSite      http.SameSite
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewHandler(service Service, domain string, sameSite http.SameSite, accessExp, refreshExp time.Duration) *Handler {
	return &Handler{service: service,
		domain:        domain,
		sameSite:      sameSite,
		accessExpiry:  accessExp,
		refreshExpiry: refreshExp,
	}
}

type RegisterRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=8"`
	Name     *string `json:"name" binding:"omitempty,min=1,max=50"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Code2FA  string `json:"code_2fa,omitempty"`
}

type Verify2FARequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Call service
	err := h.service.Register(req.Email, req.Password, req.Name)
	if err != nil {
		// Check for specific errors
		if errors.Is(err, ErrEmailAlreadyExists) {
			// Don't reveal that email exists (security)
			response.Error(c, http.StatusConflict, "Registration failed")
			return
		}

		// Internal server error
		response.Error(c, http.StatusInternalServerError, "Failed to register user")
		return
	}

	// Success
	response.Success(c, http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Extract metadata
	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	// Call service
	accessToken, refreshToken, err := h.service.Login(
		req.Email,
		req.Password,
		userAgent,
		ipAddress,
	)

	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		log.Println(err)
		response.Error(c, http.StatusInternalServerError, "Login failed")
		return
	}

	// Set refresh token in HttpOnly cookie
	setCookie(c, "refresh_token", refreshToken, h.domain, h.sameSite, int(h.refreshExpiry.Seconds())) // 7 days

	// Success - don't return refresh_token in JSON
	response.Success(c, http.StatusOK, gin.H{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   int(h.accessExpiry.Seconds()), // 15 minutes in seconds
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	// Read refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Refresh token not found")
		return
	}

	// Call service
	accessToken, newRefreshToken, err := h.service.RefreshToken(refreshToken)

	if err != nil {
		// При любой ошибке - удалить cookie
		deleteCookie(c, "refresh_token", h.domain, h.sameSite)

		if errors.Is(err, ErrInvalidToken) {
			response.Error(c, http.StatusUnauthorized, "Invalid refresh token")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Token refresh failed")
		return
	}

	// Set new refresh token in cookie
	setCookie(c, "refresh_token", newRefreshToken, h.domain, h.sameSite, int(h.refreshExpiry.Seconds()))

	// Success - don't return refresh_token in JSON
	response.Success(c, http.StatusOK, gin.H{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   int(h.accessExpiry.Seconds()),
	})
}

func (h *Handler) Logout(c *gin.Context) {
	// Read refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		// Cookie нет - уже разлогинен, но это OK
		deleteCookie(c, "refresh_token", h.domain, h.sameSite) // на всякий случай
		response.Success(c, http.StatusOK, gin.H{
			"message": "Logged out successfully",
		})
		return
	}

	// Отозвать refresh token в БД
	if err := h.service.Logout(refreshToken); err != nil {
		// Даже если ошибка в БД - удалить cookie
		deleteCookie(c, "refresh_token", h.domain, h.sameSite)
		response.Error(c, http.StatusInternalServerError, "Logout failed")
		return
	}

	// Удалить cookie
	deleteCookie(c, "refresh_token", h.domain, h.sameSite)

	response.Success(c, http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func (h *Handler) Setup2FA(c *gin.Context) {
	// TODO: Generate TOTP secret and QR code
	c.JSON(http.StatusOK, gin.H{"message": "setup 2FA endpoint"})
}

func (h *Handler) Verify2FA(c *gin.Context) {
	// TODO: Verify TOTP code
	c.JSON(http.StatusOK, gin.H{"message": "verify 2FA endpoint"})
}

func (h *Handler) Enable2FA(c *gin.Context) {
	// TODO: Enable 2FA for user
	c.JSON(http.StatusOK, gin.H{"message": "enable 2FA endpoint"})
}

func (h *Handler) Disable2FA(c *gin.Context) {
	// TODO: Disable 2FA for user
	c.JSON(http.StatusOK, gin.H{"message": "disable 2FA endpoint"})
}
