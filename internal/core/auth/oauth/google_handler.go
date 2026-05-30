package oauth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"nexa-task-tracker/internal/pkg/cookie"
	"nexa-task-tracker/internal/pkg/response"
	"time"
)

type Handler struct {
	googleOAuthService GoogleOAuthService
	frontendURL        string
	domain             string
	sameSite           http.SameSite
	accessExpiry       time.Duration
	refreshExpiry      time.Duration
}

func NewOAuthHandler(googleOAuthService GoogleOAuthService, frontendUrl, domain string,
	sameSite http.SameSite, accessExp, refreshExp time.Duration) *Handler {
	return &Handler{
		googleOAuthService: googleOAuthService,
		frontendURL:        frontendUrl,
		domain:             domain,
		sameSite:           sameSite,
		accessExpiry:       accessExp,
		refreshExpiry:      refreshExp,
	}
}

func (h *Handler) GoogleLogin(c *gin.Context) {
	state := uuid.NewString()

	url := h.googleOAuthService.GetAuthURL(state)

	response.Success(c, http.StatusOK, gin.H{
		"url": url,
	})
}

func (h *Handler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth/error?reason=missing_code")
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	refreshToken, err := h.googleOAuthService.Exchange(c.Request.Context(), code, &userAgent, &ip)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth/error?reason=exchange_failed")
		return
	}

	cookie.Set(c, "refresh_token", refreshToken, h.domain, h.sameSite, int(h.refreshExpiry))

	c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth/google/callback")
}
